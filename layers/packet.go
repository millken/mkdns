package layers

import (
	"errors"
	"sync"
)

type Packet struct {
	Ethernet Ethernet
	IPv4     IPv4
	UDP      UDP
	TCP      TCP
	Payload  []byte
}
type DecodeOptions struct {
	NoCopy bool
}

var Default = DecodeOptions{}

func ParsePacket(packet *Packet, frame []byte, options DecodeOptions) error {
	// Decode Ethernet
	packet.Ethernet = (*(*Ethernet)(&frame))
	ethType := packet.Ethernet.GetEthernetType()
	if ethType == EthernetTypeIPv4 {
		// Decode IPv4
		ipv4Raw := frame[LengthEthernet:]
		if len(ipv4Raw) < 20 {
			return errors.New("IPv4 packet too short")
		}
		packet.IPv4 = (*(*IPv4)(&ipv4Raw))

		// Decode UDP
		if packet.IPv4.GetProtocol() == IPProtocolUDP {
			ipv4HdrLen := packet.IPv4.GetIHL()
			udpRaw := ipv4Raw[ipv4HdrLen*4:]
			if len(udpRaw) < 8 {
				return errors.New("UDP packet too short")
			}
			packet.UDP = (*(*UDP)(&udpRaw))
			packet.Payload = udpRaw[LengthUDP:]
			return nil
		}
		// Decode TCP

	}

	return nil
}

var packetPool = sync.Pool{
	New: func() interface{} {
		return new(Packet)
	},
}

// AcquirePacket returns a new packet  from the pool.
func AcquirePacket() *Packet {
	return packetPool.Get().(*Packet)
}

// ReleasePacket returnes the packet to the pool.
func ReleasePacket(p *Packet) {
	packetPool.Put(p)
}
