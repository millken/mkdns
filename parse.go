package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

type PacketLayer struct {
	ethernet *layers.Ethernet
	ipv4     *layers.IPv4
	tcp      *layers.TCP
	udp      *layers.UDP
	dns      *dns.Msg
}

func parsePacket(packet gopacket.Packet) (PacketLayer, error) {
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
