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
		if req != nil {
			m.SetReply(req)
			m.SetRcode(req, dns.RcodeNameError)
		}
		m.Authoritative = true
		if e := m.IsEdns0(); e != nil {
			m.SetEdns0(4096, e.Do())
		}
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		ethMac := headers.ethernet.DstMAC
		headers.ethernet.DstMAC = headers.ethernet.SrcMAC
		headers.ethernet.SrcMAC = ethMac

		ipv4SrcIp := headers.ipv4.SrcIP
		headers.ipv4.SrcIP = headers.ipv4.DstIP
		headers.ipv4.DstIP = ipv4SrcIp

		udpSrcPort := headers.udp.SrcPort
		headers.udp.SrcPort = headers.udp.DstPort
		headers.udp.DstPort = udpSrcPort
		headers.udp.SetNetworkLayerForChecksum(headers.ipv4)
		out, err := m.Pack()
		if err != nil {
			log.Printf("dnsMsg Pack error :%s", err)
		}
		gopacket.SerializeLayers(buf, opts,
			headers.ethernet,
			headers.ipv4,
			headers.udp,
			gopacket.Payload(out))
		log.Printf("dns send packet:\n"+
			"---- ethernet ----\n "+
			"---- ipv4o ----\n "+
			"---- ipv4 ----\n "+
			"---- udp ----\n\n "+
			"---- rawPacket ----\n%v\n",
			hex.Dump(buf.Bytes()))
		err = handle.WritePacketData(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
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
