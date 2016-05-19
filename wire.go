package main

import (
	"encoding/hex"
	"io"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/miekg/dns"
	"github.com/millken/mkdns/drivers"
)

func craftAnswer(ethernetLayer *layers.Ethernet, ipLayer *layers.IPv4, dnsLayer *layers.DNS, udpLayer *layers.UDP) []byte {

	var err error
	//if not a question return
	if dnsLayer.QR {
		return nil
	}

	//must build every layer to send DNS packets
	ethMac := ethernetLayer.DstMAC
	ethernetLayer.DstMAC = ethernetLayer.SrcMAC
	ethernetLayer.SrcMAC = ethMac

	ipSrc := ipLayer.SrcIP
	ipLayer.SrcIP = ipLayer.DstIP
	ipLayer.DstIP = ipSrc

	srcPort := udpLayer.SrcPort
	udpLayer.SrcPort = udpLayer.DstPort
	udpLayer.DstPort = srcPort
	err = udpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}

	var answer layers.DNSResourceRecord
	answer.Type = layers.DNSTypeA
	answer.Class = layers.DNSClassIN
	answer.TTL = 200
	answer.IP = net.ParseIP("123.123.21.31")

	dnsLayer.QR = true

	for _, q := range dnsLayer.Questions {
		if q.Type != layers.DNSTypeA || q.Class != layers.DNSClassIN {
			continue
		}

		answer.Name = q.Name

		dnsLayer.Answers = append(dnsLayer.Answers, answer)
		dnsLayer.ANCount = dnsLayer.ANCount + 1
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	err = gopacket.SerializeLayers(buf, opts, ethernetLayer, ipLayer, udpLayer, nil)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}

	log.Printf("[DEBUG] dns :%v", hex.Dump(dnsLayer.Payload()))
	return buf.Bytes()
}

type Wire struct {
	options          *drivers.DriverOptions
	packetDataSource drivers.PacketDataSourceCloser
	isStopped        bool
}

func NewWire(options *drivers.DriverOptions) *Wire {
	i := Wire{
		options: options,
	}
	return &i
}

func (i *Wire) Start() {
	handle, err := pcap.OpenLive("enp3s0", 1024, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	handle.SetBPFFilter("dst port 53")
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		headers, err := parsePacket(packet)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[DEBUG] headers : %s", headers.dns)
		req := headers.dns
		m := new(dns.Msg)
		m.SetReply(req)
		m.Authoritative = true
		if e := m.IsEdns0(); e != nil {
			m.SetEdns0(4096, e.Do())
		}
		m.SetRcode(req, dns.RcodeNameError)
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		ethernetPacket := &layers.Ethernet{
			SrcMAC: headers.ethernet.DstMAC,
			DstMAC: headers.ethernet.SrcMAC,
		}
		ipv4Packet := &layers.IPv4{
			SrcIP: headers.ipv4.DstIP,
			DstIP: headers.ipv4.SrcIP,
		}
		udpPacket := &layers.UDP{
			SrcPort: headers.udp.DstPort,
			DstPort: headers.udp.SrcPort,
		}
		udpPacket.SetNetworkLayerForChecksum(ipv4Packet)
		out, err := m.Pack()
		if err != nil {
			log.Printf("dnsMsg Pack error :%s", err)
		}
		gopacket.SerializeLayers(buf, opts,
			ethernetPacket,
			ipv4Packet,
			udpPacket,
			gopacket.Payload(out))
		log.Printf("dns send packet: ---- rawPacket ----\n%v\n", hex.Dump(buf.Bytes()))
		handle.WritePacketData(buf.Bytes())
	}
}

func (i *Wire) Stop() {
	log.Print("[INFO] sniffer: sending stopCapureChan signal")
	i.isStopped = true
}

func (i *Wire) Close() {
	if i.packetDataSource != nil {
		log.Print("[INFO] closing packet capture socket")
		i.packetDataSource.Close()
	}
	log.Print("[INFO] stopping the sniffer decode loop")
	i.isStopped = true
	log.Print("[INFO] done.")
}
