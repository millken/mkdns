package layers

import (
	"encoding/binary"
)

// struct udphdr {
//	__be16	source;
//	__be16	dest;
//	__be16	len;
//	__sum16	check;
// };

type UDP []byte

const LengthUDP = 8

func (u *UDP) GetSrcPort() uint16 {
	return binary.BigEndian.Uint16((*u)[0:2])
}

func (u *UDP) SetSrcPort(p uint16) {
	binary.BigEndian.PutUint16((*u)[0:2], p)
}

func (u *UDP) GetDstPort() uint16 {
	return binary.BigEndian.Uint16((*u)[2:4])
}

func (u *UDP) SetDstPort(p uint16) {
	binary.BigEndian.PutUint16((*u)[2:4], p)
}

func (u *UDP) GetLength() uint16 {
	return binary.BigEndian.Uint16((*u)[4:6])
}

func (u *UDP) SetLength(l uint16) {
	binary.BigEndian.PutUint16((*u)[4:6], l)
}

func (u *UDP) GetChecksum() uint16 {
	return binary.BigEndian.Uint16((*u)[6:8])
}

func (u *UDP) SetChecksum(l uint16) {
	binary.BigEndian.PutUint16((*u)[6:8], l)
}
