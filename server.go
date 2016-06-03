package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gopacket"
	"github.com/millken/mkdns/drivers"
	"github.com/millken/mkdns/stats"
	"github.com/millken/mkdns/types"
)

const BPFFilter = "udp and dst port 53"

type server struct {
	config        *Config
	io            drivers.PacketDataSourceCloser
	stats         stats.Directional
	isStopped     bool
	forceQuitChan chan os.Signal
	txChan        chan types.PacketLayer
	rxChan        chan gopacket.Packet
	//handler          *Handler
}

func NewServer(config *Config) *server {
	return &server{
		config:        config,
		forceQuitChan: make(chan os.Signal, 1),
		txChan:        make(chan types.PacketLayer),
		rxChan:        make(chan gopacket.Packet),
		isStopped:     false,
	}
}

func (s *server) Start() (err error) {
	options := &drivers.DriverOptions{
		Device:  "enp3s0",
		Snaplen: 2048,
		Filter:  BPFFilter,
	}

	factory, ok := drivers.Drivers["libpcap"]
	if !ok {
		log.Fatal(fmt.Sprintf("%s Packet driver not supported on this system", s.config.Server.Driver))
	}

	s.io, err = factory(options)
	if err != nil {
		return
	}
	for i := 0; i < 1; i++ {
		go packetHandler(i, s.rxChan, s.txChan)
	}
	go s.readPackets()
	go s.sendPackets()
	s.signalWorker()
	return
}

func (s *server) sendPackets() {
	defer close(s.txChan)

	for {
		p := (<-s.txChan)
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		ethMac := p.Ethernet.DstMAC
		p.Ethernet.DstMAC = p.Ethernet.SrcMAC
		p.Ethernet.SrcMAC = ethMac

		ipv4SrcIp := p.Ipv4.SrcIP
		p.Ipv4.SrcIP = p.Ipv4.DstIP
		p.Ipv4.DstIP = ipv4SrcIp

		out, err := p.Dns.Pack()
		if err != nil {
			log.Printf("dnsMsg Pack error :%s", err)
		}
		if p.Udp != nil {
			udpSrcPort := p.Udp.SrcPort
			p.Udp.SrcPort = p.Udp.DstPort
			p.Udp.DstPort = udpSrcPort
			p.Udp.SetNetworkLayerForChecksum(p.Ipv4)
			gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Udp, gopacket.Payload(out))
		}
		if p.Tcp != nil {
			tcpSrcPort := p.Tcp.SrcPort
			p.Tcp.SrcPort = p.Tcp.DstPort
			p.Tcp.DstPort = tcpSrcPort
			p.Tcp.PSH = true
			p.Tcp.ACK = true
			p.Tcp.FIN = true
			tcpSeq := p.Tcp.Seq
			p.Tcp.Seq = p.Tcp.Ack
			p.Tcp.Ack = tcpSeq + uint32(len(p.Tcp.LayerPayload()))
			p.Tcp.Window = 0
			p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
			gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp, gopacket.Payload(out))
		}
		err = s.io.WritePacketData(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		s.stats.Tx.Packets++
		//d.stats.Tx.Bytes += uint64(p.Metadata().CaptureInfo.CaptureLength)
	}
}

func (s *server) readPackets() {
	packetSource := s.io.PacketSource()
	defer close(s.rxChan)
	for {

		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("[ERROR] readPackets err: %s", err)
			continue
		}
		s.stats.Rx.Packets++

		s.rxChan <- packet
		if s.isStopped {
			break
		}
	}
}

func (s *server) Shutdown() {
	s.io.Close()
	s.isStopped = true
}

func (s *server) Stats() stats.Directional {
	return s.stats
}

func (s *server) signalWorker() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGINT)

	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP:
			log.Println("Reload initiated.")
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Shutdown()
			log.Println("Shutdown initiated.")
			return
		}
	}
}
