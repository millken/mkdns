package types

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

type PacketLayer struct {
	Ethernet *layers.Ethernet
	Ipv4     *layers.IPv4
	Tcp      *layers.TCP
	Udp      *layers.UDP
	Dns      *dns.Msg
}

func DecodeByProtobuff(data []byte) (r Records, err error) {
	r = Records{}
	err = proto.Unmarshal(data, &r)
	return
}

func ParsePacket(packet gopacket.Packet) (PacketLayer, error) {
	layer := PacketLayer{}
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		layer.Ethernet = ethernetLayer.(*layers.Ethernet)
		log.Printf("[DEBUG] SrcMac: %s -> DstMac: %s", layer.Ethernet.SrcMAC, layer.Ethernet.DstMAC)
	}
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		layer.Ipv4 = ipLayer.(*layers.IPv4)
		log.Printf("[DEBUG] IPv4: SrcIP:  %s -> DstIP: %s", layer.Ipv4.SrcIP, layer.Ipv4.DstIP)
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		layer.Udp = udpLayer.(*layers.UDP)
		log.Printf("[DEBUG] UDP: SrcPort:  %d -> DstPort: %d", layer.Udp.SrcPort, layer.Udp.DstPort)
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		layer.Tcp = tcpLayer.(*layers.TCP)
		log.Printf("[DEBUG] TCP: SrcPort:  %d -> DstPort: %d", layer.Tcp.SrcPort, layer.Tcp.DstPort)
		log.Printf("[DEBUG] TCP: SYN:  %v | ACK: %v | FIN: %v", layer.Tcp.SYN, layer.Tcp.ACK, layer.Tcp.FIN)
	}

	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dnsLayerMsg := dnsLayer.(*layers.DNS)

		contents := dnsLayerMsg.BaseLayer.LayerContents()

		dnsMsg := new(dns.Msg)
		if err := dnsMsg.Unpack(contents); err != nil {
			return layer, err
		}

		layer.Dns = dnsMsg
	}

	return layer, nil
}
