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
)

const BPFFilter = "udp and dst port 53"

type server struct {
	config        *Config
	io            drivers.PacketDataSourceCloser
	stats         stats.Directional
	isStopped     bool
	forceQuitChan chan os.Signal
	txChan        chan PacketLayer
	rxChan        chan gopacket.Packet
	//handler          *Handler
}

func NewServer(config *Config) *server {
	return &server{
		config:        config,
		forceQuitChan: make(chan os.Signal, 1),
		txChan:        make(chan PacketLayer),
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
		ethMac := p.ethernet.DstMAC
		p.ethernet.DstMAC = p.ethernet.SrcMAC
		p.ethernet.SrcMAC = ethMac

		ipv4SrcIp := p.ipv4.SrcIP
		p.ipv4.SrcIP = p.ipv4.DstIP
		p.ipv4.DstIP = ipv4SrcIp

		out, err := p.dns.Pack()
		if err != nil {
			log.Printf("dnsMsg Pack error :%s", err)
		}
		if p.udp != nil {
			udpSrcPort := p.udp.SrcPort
			p.udp.SrcPort = p.udp.DstPort
			p.udp.DstPort = udpSrcPort
			p.udp.SetNetworkLayerForChecksum(p.ipv4)
			gopacket.SerializeLayers(buf, opts, p.ethernet, p.ipv4, p.udp, gopacket.Payload(out))
		}
		if p.tcp != nil {
			tcpSrcPort := p.tcp.SrcPort
			p.tcp.SrcPort = p.tcp.DstPort
			p.tcp.DstPort = tcpSrcPort
			p.tcp.PSH = true
			p.tcp.ACK = true
			p.tcp.FIN = true
			tcpSeq := p.tcp.Seq
			p.tcp.Seq = p.tcp.Ack
			p.tcp.Ack = tcpSeq + uint32(len(p.tcp.LayerPayload()))
			p.tcp.Window = 0
			p.tcp.SetNetworkLayerForChecksum(p.ipv4)
			gopacket.SerializeLayers(buf, opts, p.ethernet, p.ipv4, p.tcp, gopacket.Payload(out))
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
