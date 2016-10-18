package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/miekg/dns"
	"github.com/millken/mkdns/backends"
	"github.com/millken/mkdns/drivers"
	"github.com/millken/mkdns/stats"
	"github.com/millken/mkdns/types"
)

const BPFFilter = "dst port 53"

type server struct {
	config        *Config
	io            drivers.PacketDataSourceCloser
	isStopped     bool
	forceQuitChan chan os.Signal
	txChan        chan types.PacketLayer
	rxChan        chan gopacket.Packet
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
		Device:  s.config.Server.Iface,
		Snaplen: 2048,
		Filter:  BPFFilter,
	}

	factory, ok := drivers.Drivers[s.config.Server.Driver]
	if !ok {
		log.Fatal(fmt.Sprintf("%s Packet driver not supported on this system", s.config.Server.Driver))
	}

	s.io, err = factory(options)
	if err != nil {
		return fmt.Errorf("driver: %s, interface: %s boot error: %s", s.config.Server.Driver, s.config.Server.Iface, err)
	}
	worker_num := 7
	if s.config.Server.WorkerNum > 0 {
		worker_num = s.config.Server.WorkerNum
	}
	for i := 0; i < worker_num; i++ {
		go s.decodePackets(i)
	}
	go s.readPackets()
	go s.sendPackets()
	s.signalWorker()
	return
}

func (s *server) decodePackets(worker_id int) {
	for packet := range s.rxChan {
		if s.isStopped {
			return
		}
		p, err := types.ParsePacket(packet)
		if err != nil || p.Dns == nil {
			// response syn ->ack
			if p.Tcp == nil {
				continue
			}
			if p.Tcp.SYN && !p.Tcp.ACK && !p.Tcp.FIN {
				buf := gopacket.NewSerializeBuffer()
				opts := gopacket.SerializeOptions{
					FixLengths:       true,
					ComputeChecksums: true,
				}
				ethMac := p.Ethernet.DstMAC
				p.Ethernet.DstMAC = p.Ethernet.SrcMAC
				p.Ethernet.SrcMAC = ethMac

				ipv4SrcIp := p.Ipv4.SrcIP
				p.Ipv4.Id = 0
				p.Ipv4.SrcIP = p.Ipv4.DstIP
				p.Ipv4.DstIP = ipv4SrcIp
				tcpSrcPort := p.Tcp.SrcPort
				p.Tcp.SrcPort = p.Tcp.DstPort
				p.Tcp.DstPort = tcpSrcPort
				p.Tcp.SYN = true
				p.Tcp.ACK = true
				p.Tcp.RST = false
				tcpSeq := p.Tcp.Seq
				p.Tcp.Seq = p.Tcp.Ack
				p.Tcp.Ack = uint32(tcpassembly.Sequence(tcpSeq).Add(1))
				//p.Tcp.Window = 512
				p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
				gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp)
				err = s.io.WritePacketData(buf.Bytes())
				continue
			}

			if p.Tcp.ACK && p.Tcp.PSH {
				applicationLayer := packet.ApplicationLayer()

				dnsMsg := new(dns.Msg)
				//tcp mode, remove length
				if err := dnsMsg.Unpack(applicationLayer.Payload()[2:]); err != nil {
					log.Printf("unpack dns error : %s, str: %s", err, hex.Dump(applicationLayer.Payload()))
					continue
				}
				p.Dns = dnsMsg
			}

			//response FIN && ACK
			if p.Tcp.FIN && p.Tcp.ACK {
				buf := gopacket.NewSerializeBuffer()
				opts := gopacket.SerializeOptions{
					FixLengths:       true,
					ComputeChecksums: true,
				}
				ethMac := p.Ethernet.DstMAC
				p.Ethernet.DstMAC = p.Ethernet.SrcMAC
				p.Ethernet.SrcMAC = ethMac

				ipv4SrcIp := p.Ipv4.SrcIP
				//p.Ipv4.Id = p.Ipv4.Id + 1
				p.Ipv4.SrcIP = p.Ipv4.DstIP
				p.Ipv4.DstIP = ipv4SrcIp
				tcpSrcPort := p.Tcp.SrcPort
				p.Tcp.SrcPort = p.Tcp.DstPort
				p.Tcp.DstPort = tcpSrcPort
				p.Tcp.FIN = false
				p.Tcp.ACK = true
				p.Tcp.RST = false
				tcpSeq := p.Tcp.Seq
				p.Tcp.Seq = p.Tcp.Ack
				p.Tcp.Ack = uint32(tcpassembly.Sequence(tcpSeq).Add(1))
				//p.Tcp.Window = 512
				p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
				gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp)
				err = s.io.WritePacketData(buf.Bytes())
				continue
			}
			log.Printf("[ERROR] parsePacket : %+v", err)
			if p.Dns == nil {
				continue
			}
		}
		req := p.Dns
		q := req.Question[0]
		domain := strings.ToLower(q.Name)
		zz, err := backends.GetRecords(domain)
		if err != nil {
			log.Printf("[WARN] %s %s", domain, err)
		}

		m := new(dns.Msg)
		if req != nil {
			m.SetReply(req)
			m.SetRcode(req, dns.RcodeNameError)
		}
		if e := m.IsEdns0(); e != nil {
			m.SetEdns0(4096, e.Do())
		}

		var zname string
		if zz == nil {
			zname = "NULL"
		} else {
			zname = zz.Name
		}
		log.Printf("[FINE] [zone %s] incoming %s %s %d from %s", zname, req.Question[0].Name,
			dns.TypeToString[q.Qtype], req.MsgHdr.Id, p.Ipv4.SrcIP)

		if zz == nil {
			m.SetRcode(req, dns.RcodeNameError)
		} else {
			counterKey := fmt.Sprintf("%s:%s", zz.Name, dns.TypeToString[q.Qtype])
			stats.NewCounter(counterKey).Add(1)
			zz.Options.EdnsAddr = nil
			zz.Options.RemoteAddr = p.Ipv4.SrcIP

			//var edns *dns.EDNS0_SUBNET
			//var opt_rr *dns.OPT

			for _, extra := range req.Extra {

				switch extra.(type) {
				case *dns.OPT:
					for _, o := range extra.(*dns.OPT).Option {
						//opt_rr = extra.(*dns.OPT)
						switch e := o.(type) {
						case *dns.EDNS0_NSID:
						case *dns.EDNS0_SUBNET:
							//log.Printf("[DEBUG] Got edns", e.Address, e.Family, e.SourceNetmask, e.SourceScope)
							if e.Address != nil {
								//edns = e
								zz.Options.EdnsAddr = e.Address
							}
						}
					}
				}
			}

			if q.Qclass == dns.ClassCHAOS {
				if q.Qtype == dns.TypeTXT {
					switch domain {
					case "bind.version":
						fallthrough
					case "id.server.":
						hdr := dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
						m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{"localhost"}}}
						return
					}
				}
			}
			m, err = zz.FindRecord(req)
			if err != nil {
				m.Ns = append(m.Ns, zz.SoaRR())
				log.Printf("[ERROR] zone error : %s", err)
			} else {
				m.Ns = zz.NsRR()
			}
		}
		m.Authoritative = true
		p.Dns = m
		s.txChan <- p
	}
}

func (s *server) sendPackets() {
	defer close(s.txChan)

	for {
		select {
		case <-s.forceQuitChan:
			return
		case p := (<-s.txChan):

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
				err = s.io.WritePacketData(buf.Bytes())
			}
			if p.Tcp != nil {
				buf := gopacket.NewSerializeBuffer()
				//p.Ipv4.Id = p.Ipv4.Id + 1
				tcpSrcPort := p.Tcp.SrcPort
				p.Tcp.SrcPort = p.Tcp.DstPort
				p.Tcp.DstPort = tcpSrcPort
				p.Tcp.Options = nil

				//send ack
				p.Tcp.SYN = false
				p.Tcp.ACK = true
				p.Tcp.RST = false
				p.Tcp.PSH = false
				//seq := p.Tcp.Seq
				tcpSeq := p.Tcp.Seq
				p.Tcp.Seq = p.Tcp.Ack
				p.Tcp.Ack = tcpSeq + uint32(len(p.Tcp.LayerPayload()))
				//p.Tcp.Window = 512
				p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
				gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp)
				err = s.io.WritePacketData(buf.Bytes())

				buf = gopacket.NewSerializeBuffer()
				//p.Ipv4.Id = p.Ipv4.Id + 1
				p.Tcp.PSH = true
				p.Tcp.ACK = true
				//p.Tcp.Seq = 1

				//p.Tcp.Ack = uint32(tcpassembly.Sequence(seq).Add(1))
				p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
				bs := make([]byte, 2)
				binary.BigEndian.PutUint16(bs, uint16(len(out)))
				gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp, gopacket.Payload(append(bs, out...)))
				err = s.io.WritePacketData(buf.Bytes())
			}

			if err != nil {
				log.Fatal(err)
			}
			if p.Tcp != nil {
				p.Ipv4.Id = p.Ipv4.Id + 1
				p.Tcp.SYN = false
				p.Tcp.ACK = true
				p.Tcp.FIN = true
				p.Tcp.PSH = false
				p.Tcp.Seq = uint32(len(out) + 3)
				//p.Tcp.Window = 0
				p.Tcp.SetNetworkLayerForChecksum(p.Ipv4)
				gopacket.SerializeLayers(buf, opts, p.Ethernet, p.Ipv4, p.Tcp)
				err = s.io.WritePacketData(buf.Bytes())
			}
		}
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
		s.rxChan <- packet
		if s.isStopped {
			break
		}
	}
}

func (s *server) Shutdown() {
	/*
		if s.io != nil {
			log.Print("[INFO] closing packet capture socket")
			s.io.Close()
		}
	*/
	log.Print("[INFO] stopping loop")
	s.isStopped = true

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
			s.forceQuitChan <- sig
			log.Println("Shutdown initiated.")
			return
		}
	}
}
