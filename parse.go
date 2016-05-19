package main

import (
	"github.com/millken/mkdns/protocols"
	"github.com/miekg/dns"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)
type PacketLayer struct {
	ethernet *layers.Ethernet
	ipv4   *layers.IPv4
    tcp *layers.TCP
	udp *layers.UDP
	dns *dns.Msg
}
func parsePacket(packet gopacket.Packet)(PacketLayer, error ) {
	layer := PacketLayer{}
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		layer.ethernet = ethernetLayer.(*layers.Ethernet)
	}
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		layer.ipv4 = ipLayer.(*layers.IPv4)
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		layer.udp = udpLayer.(*layers.UDP)
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		layer.tcp = tcpLayer.(*layers.TCP)
	}

	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dnsLayerMsg := dnsLayer.(*layers.DNS)

		contents := dnsLayerMsg.BaseLayer.LayerContents()

		dnsMsg := new(dns.Msg)
		if err := dnsMsg.Unpack(contents); err != nil {
			return layer, err
		}

		layer.dns = dnsMsg
	}
	
	return layer, nil
}

// Parse parses a packet header.
func Parse(packet gopacket.Packet) (map[string]interface{}, error) {
	packetHeaders := make(map[string]interface{})

	metaData := packet.Metadata()

	// Include packet timestamp
	packetHeaders["timestamp"] = (&metaData.CaptureInfo.Timestamp).String()

	// If this packet has an Ethernet frame, include it's header
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		packetHeaders["ethernet"] = ethernetLayer.(*layers.Ethernet)
	}

	// If this is an ICMP packet, include it's header
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer != nil {
		packetHeaders["icmpv4"] = protocols.ICMPv4Parser(icmpLayer)
	}

	// If this is an IPv4 packet, include it's header
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		packetHeaders["ipv4"] = protocols.IPv4Parser(ipLayer)
	}

	// If this is a UDP datagram, include it's header
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		packetHeaders["udp"] = protocols.UDPParser(udpLayer)
	}

	// If this is a TCP segment, include it's header
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		packetHeaders["tcp"] = protocols.TCPParser(tcpLayer)
	}

	// If this packet has a DNS payload, include it's data
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, err := protocols.DNSParser(dnsLayer)
		if err != nil {
			return nil, err
		}
		packetHeaders["dns"] = dns
	}

	return packetHeaders, nil
}
