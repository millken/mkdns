package layers

import (
	"fmt"
	"testing"
)

func TestTCP_GetAll(t *testing.T) {
	p := []byte{
		0x01, 0xbb, 0xec, 0x91,
		0x39, 0x98, 0x3d, 0xa5,
		0x9f, 0xb8, 0x7b, 0x02,
		0x50, 0x10, 0x08, 0x00,
		0x36, 0x0b, 0x00, 0x00,
	}

	tcp := *(*TCP)(&p)

	fmt.Println(Swap16(tcp.GetSrcPort()))
	fmt.Println(Swap16(tcp.GetDstPort()))
	fmt.Println(Swap32(tcp.GetSeq()))
	fmt.Println(Swap32(tcp.GetAckSeq()))
	fmt.Println(tcp.GetDataOffset())
	fmt.Println(tcp.GetReserved())
	fmt.Println(tcp.IsFlagCWR())
	fmt.Println(tcp.IsFlagECE())
	fmt.Println(tcp.IsFlagUrg())
	fmt.Println(tcp.IsFlagAck())
	fmt.Println(tcp.IsFlagPsh())
	fmt.Println(tcp.IsFlagRst())
	fmt.Println(tcp.IsFlagSyn())
	fmt.Println(tcp.IsFlagFin())
	fmt.Println(Swap16(tcp.GetWindow()))
	fmt.Println(Swap16(tcp.GetChecksum()))
	fmt.Println(Swap16(tcp.GetUrgPointer()))
}

func TestTCP_SetAll(t *testing.T) {
	p := make([]byte, 20)

	tcp := *(*TCP)(&p)
	tcp.SetSrcPort(Swap16(443))
	tcp.SetDstPort(Swap16(60561))
	tcp.SetSeq(Swap32(966278565))
	tcp.SetAckSeq(Swap32(2679667458))
	tcp.SetDataOffset(20)
	tcp.SetFlagCWR(false)
	tcp.SetFlagECE(false)
	tcp.SetFlagUrg(false)
	tcp.SetFlagAck(true)
	tcp.SetFlagPsh(false)
	tcp.SetFlagRst(false)
	tcp.SetFlagSyn(false)
	tcp.SetFlagFin(false)
	tcp.SetWindow(Swap16(2048))
	tcp.SetChecksum(Swap16(0x360b))
	tcp.SetUrgPointer(0)

	fmt.Printf("%x\n", p)
}
