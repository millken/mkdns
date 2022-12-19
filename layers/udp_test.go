package layers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUDP_GetAll(t *testing.T) {
	require := require.New(t)
	p := []byte{
		0xef, 0x81, 0x00, 0x35,
		0x00, 0x24, 0x1f, 0x32,
	}

	udp := *(*UDP)(&p)

	require.Equal(uint16(0xef81), udp.GetSrcPort())
	require.Equal(uint16(0x0035), udp.GetDstPort())
	require.Equal(uint16(0x0024), udp.GetLength())
	require.Equal(uint16(0x1f32), udp.GetChecksum())
}

func TestUDP_SetAll(t *testing.T) {
	require := require.New(t)
	p := make([]byte, 8)

	udp := *(*UDP)(&p)
	udp.SetSrcPort(61313)
	udp.SetDstPort(53)
	udp.SetLength(36)
	udp.SetChecksum(7986)
	fmt.Printf("%x\n", p)
	require.Equal(uint16(0xef81), udp.GetSrcPort())
	require.Equal(uint16(0x0035), udp.GetDstPort())
	require.Equal(uint16(0x0024), udp.GetLength())
	require.Equal(uint16(0x1f32), udp.GetChecksum())
}

func BenchmarkUDPSet(b *testing.B) {
	p := make([]byte, 8)
	udp := *(*UDP)(&p)
	for i := 0; i < b.N; i++ {
		udp.SetSrcPort(61313)
		udp.SetDstPort(53)
		udp.SetLength(36)
		udp.SetChecksum(7986)
	}
}

func BenchmarkUDPGet(b *testing.B) {
	p := []byte{
		0xef, 0x81, 0x00, 0x35,
		0x00, 0x24, 0x1f, 0x32,
	}
	udp := *(*UDP)(&p)
	for i := 0; i < b.N; i++ {
		udp.GetSrcPort()
		udp.GetDstPort()
		udp.GetLength()
		udp.GetChecksum()
	}
}
