package layers

import "unsafe"

// Calculate the TCP/IP checksum defined in rfc1071.  The passed-in csum is any
// initial checksum data that's already been computed.
// GRE, ICMPv4, TCP/IP can use it.
//go:inline
func TCPIPChecksum(data []byte, baseCSum uint32) uint16 {
	// 避免重复获取长度
	length := len(data)
	// 计算偶数部分
	for i := 0; i < length>>1; i++ {
		baseCSum += uint32(data[i*2])<<8 + uint32(data[i*2+1])
	}
	// 如果是奇数就把最后一位加上
	if length&0x01 == 0x01 {
		baseCSum += uint32(data[length]) << 8
	}
	for baseCSum > 0xffff {
		baseCSum = (baseCSum >> 16) + (baseCSum & 0xffff)
	}
	return ^uint16(baseCSum)
}

func checksum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}
	for {
		// Break when sum is less or equals to 0xFFFF
		if csum <= 65535 {
			break
		}
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}

func IsBigEndian() bool {
	var i uint16 = 0x0001
	return (*[2]byte)(unsafe.Pointer(&i))[0] == 0x00
}

func Swap16(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func Swap32(i uint32) uint32 {
	b0 := (i & 0x000000ff) << 24
	b1 := (i & 0x0000ff00) << 8
	b2 := (i & 0x00ff0000) >> 8
	b3 := (i & 0xff000000) >> 24

	return b0 | b1 | b2 | b3
}
