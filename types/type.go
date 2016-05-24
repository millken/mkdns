package types

import (
	"log"

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
		log.Printf("[DEBUG] SrcMac: %s -> DstMac: %s", layer.ethernet.SrcMAC, layer.ethernet.DstMAC)
	}
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		layer.ipv4 = ipLayer.(*layers.IPv4)
		log.Printf("[DEBUG] IPv4: SrcIP:  %s -> DstIP: %s", layer.ipv4.SrcIP, layer.ipv4.DstIP)
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		layer.udp = udpLayer.(*layers.UDP)
		log.Printf("[DEBUG] UDP: SrcPort:  %d -> DstPort: %d", layer.udp.SrcPort, layer.udp.DstPort)
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		layer.tcp = tcpLayer.(*layers.TCP)
		log.Printf("[DEBUG] TCP: SrcPort:  %d -> DstPort: %d", layer.tcp.SrcPort, layer.tcp.DstPort)
		log.Printf("[DEBUG] TCP: SYN:  %v | ACK: %v | FIN: %v", layer.tcp.SYN, layer.tcp.ACK, layer.tcp.FIN)
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
