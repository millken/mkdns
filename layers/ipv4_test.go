package layers

import (
	"fmt"
	"net"
	"testing"
)

func TestIPv4_GetAll(t *testing.T) {
	p := []byte{
		0x45,
		0x00,
		0x00, 0x3c,
		0x97, 0x8b,
		0x00, 0x00,
		0x7f,
		0x01,
		0x78, 0x4a,
		0x64, 0x61, 0x51, 0x6b,
		0x64, 0x63, 0x11, 0xbc,
	}

	ip4 := *(*IPv4)(&p)

	fmt.Println(ip4.GetVersion())
	fmt.Println(ip4.GetIHL())
	fmt.Println(ip4.GetTOS())
	fmt.Println(Swap16(ip4.GetTotalLen()))
	fmt.Println(Swap16(ip4.GetID()))
	fmt.Println(Swap16(ip4.GetFragOff()))
	fmt.Println(ip4.GetTTL())
	fmt.Println(ip4.GetProtocol())
	fmt.Println(Swap16(ip4.GetChecksum()))
	fmt.Println(ip4.GetSrcAddr())
	fmt.Println(ip4.GetDstAddr())
	fmt.Println(ip4.IsFlagReserved())
	fmt.Println(ip4.IsFlagDontFrag())
	fmt.Println(ip4.IsFlagMoreFrag())
}

func TestIPv4_SetAll(t *testing.T) {
	p := make([]byte, 20)

	ipv4 := *(*IPv4)(&p)
	ipv4.SetVersion(4)
	ipv4.SetIHL(20)
	ipv4.SetTOS(64)
	ipv4.SetTotalLen(Swap16(20))
	ipv4.SetID(Swap16(1))
	ipv4.SetFlagReserved(false)
	ipv4.SetFlagDontFrag(true)
	ipv4.SetFlagMoreFrag(false)
	ipv4.SetFragOff(Swap16(0))
	ipv4.SetTTL(6)
	ipv4.SetProtocol(1)
	ipv4.SetSrcAddr(net.IP{1, 1, 1, 1})
	ipv4.SetDstAddr(net.IP{2, 2, 2, 2})

	ipv4.SetChecksum(0)
	ipv4.SetChecksum(Swap16(TCPIPChecksum(p[:20], 0)))

	fmt.Printf("%x\n", p)
}
