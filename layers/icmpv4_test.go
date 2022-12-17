package layers

import (
	"fmt"
	"testing"
)

func TestICMPv4_GetAll(t *testing.T) {
	p := []byte{
		0x08, 0x00, 0x8f, 0x3e, 0x04, 0x04, 0x00, 0x01,
	}

	icmp4 := *(*ICMPv4)(&p)
	fmt.Println(icmp4.GetType())
	fmt.Println(icmp4.GetCode())
	fmt.Println(Swap16(icmp4.GetChecksum()))
	fmt.Println(Swap16(icmp4.GetID()))
	fmt.Println(Swap16(icmp4.GetSequence()))
}

func TestICMPv4_SetAll(t *testing.T) {
	p := make([]byte, 8)
	icmp4 := *(*ICMPv4)(&p)
	icmp4.SetType(ICMPv4TypeEchoRequest)
	icmp4.SetCode(ICMPv4CodeNet)
	icmp4.SetChecksum(0)
	icmp4.SetID(Swap16(1))
	icmp4.SetSequence(Swap16(1))

	icmp4.SetChecksum(Swap16(TCPIPChecksum(p, 0)))

	fmt.Printf("%x\n", p)
}
