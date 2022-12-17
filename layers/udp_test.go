package layers

import (
	"fmt"
	"testing"
)

func TestUDP_GetAll(t *testing.T) {
	p := []byte{
		0xef, 0x81, 0x00, 0x35,
		0x00, 0x24, 0x1f, 0x32,
	}

	udp := *(*UDP)(&p)

	fmt.Println(Swap16(udp.GetSrcPort()))
	fmt.Println(Swap16(udp.GetDstPort()))
	fmt.Println(Swap16(udp.GetLen()))
	fmt.Println(Swap16(udp.GetChecksum()))
}

func TestUDP_SetAll(t *testing.T) {
	p := make([]byte, 8)

	udp := *(*UDP)(&p)
	udp.SetSrcPort(Swap16(61313))
	udp.SetDstPort(Swap16(53))
	udp.SetLen(Swap16(36))
	udp.SetChecksum(Swap16(7986))

	fmt.Printf("%x\n", p)
}
