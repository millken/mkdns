package layers

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPv4_GetAll(t *testing.T) {
	require := require.New(t)
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
	require.Equal(uint8(4), ip4.GetVersion())
	require.Equal(uint8(5), ip4.GetIHL())
	require.Equal(uint8(0), ip4.GetTOS())
	require.Equal(uint16(60), ip4.GetLength())
	require.Equal(uint16(0x978b), ip4.GetID())
	require.Equal(uint16(0), ip4.GetFragOff())
	require.Equal(uint8(127), ip4.GetTTL())
	require.Equal(uint8(1), ip4.GetProtocol())
	require.Equal(uint16(0x784a), ip4.GetChecksum())
	require.Equal(net.IP{0x64, 0x61, 0x51, 0x6b}, ip4.GetSrcAddr())
	require.Equal(net.IP{0x64, 0x63, 0x11, 0xbc}, ip4.GetDstAddr())
	require.False(ip4.IsFlagReserved())
	require.False(ip4.IsFlagDontFrag())
	require.False(ip4.IsFlagMoreFrag())

	fmt.Println(ip4.GetVersion())
	fmt.Println(ip4.GetIHL())
	fmt.Println(ip4.GetTOS())
	fmt.Println(Swap16(ip4.GetLength()))
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
	require := require.New(t)
	p := make([]byte, 20)

	ipv4 := *(*IPv4)(&p)
	ipv4.SetVersion(4)
	ipv4.SetIHL(5)
	ipv4.SetTOS(64)
	ipv4.SetLength(20)
	ipv4.SetID(11)
	ipv4.SetFlagReserved(false)
	ipv4.SetFlagDontFrag(true)
	ipv4.SetFlagMoreFrag(false)
	ipv4.SetFragOff()
	ipv4.SetTTL(6)
	ipv4.SetProtocol(1)
	ipv4.SetSrcAddr(net.IP{1, 1, 1, 1})
	ipv4.SetDstAddr(net.IP{2, 2, 2, 2})

	ipv4.SetChecksum()

	fmt.Printf("%x\n", p)
	require.Equal(uint8(4), ipv4.GetVersion())
	require.Equal(uint8(5), ipv4.GetIHL())
	require.Equal(uint8(64), ipv4.GetTOS())
	require.Equal(uint16(20), ipv4.GetLength())
	require.Equal(uint16(11), ipv4.GetID())
	require.Equal(uint16(0), ipv4.GetFragOff())
	require.Equal(uint8(6), ipv4.GetTTL())
	require.Equal(uint8(1), ipv4.GetProtocol())
	require.Equal(uint16(0x6e99), ipv4.GetChecksum())
	require.Equal(net.IP{1, 1, 1, 1}, ipv4.GetSrcAddr())
	require.Equal(net.IP{2, 2, 2, 2}, ipv4.GetDstAddr())
	require.False(ipv4.IsFlagReserved())
	require.True(ipv4.IsFlagDontFrag())
	require.False(ipv4.IsFlagMoreFrag())
}

func TestIPv4_Checksum(t *testing.T) {
	require := require.New(t)
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
	ip4.SetChecksum()
	require.Equal(uint16(0x784a), ip4.GetChecksum())
}

func BenchmarkIPv4Set(b *testing.B) {
	p := make([]byte, 20)
	ipv4 := *(*IPv4)(&p)
	for i := 0; i < b.N; i++ {
		ipv4.SetVersion(4)
		ipv4.SetIHL(5)
		ipv4.SetTOS(64)
		ipv4.SetLength(20)
		ipv4.SetID(11)
		ipv4.SetFlagReserved(false)
		ipv4.SetFlagDontFrag(true)
		ipv4.SetFlagMoreFrag(false)
		ipv4.SetFragOff()
		ipv4.SetTTL(6)
		ipv4.SetProtocol(1)
		ipv4.SetSrcAddr(net.IP{1, 1, 1, 1})
		ipv4.SetDstAddr(net.IP{2, 2, 2, 2})
		ipv4.SetChecksum()
	}
}

func BenchmarkIPv4Get(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		ip4.GetVersion()
		ip4.GetIHL()
		ip4.GetTOS()
		ip4.GetLength()
		ip4.GetID()
		ip4.GetFragOff()
		ip4.GetTTL()
		ip4.GetProtocol()
		ip4.GetChecksum()
		ip4.GetSrcAddr()
		ip4.GetDstAddr()
		ip4.IsFlagReserved()
		ip4.IsFlagDontFrag()
		ip4.IsFlagMoreFrag()
	}
}
