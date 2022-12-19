package layers

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEthernet_GetAll(t *testing.T) {
	require := require.New(t)
	p := []byte{
		0x94, 0x94, 0x26, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x08, 0x00,
	}

	eth := *(*Ethernet)(&p)
	mac := eth.GetSrcAddress()
	fmt.Println(mac.String())
	mac = eth.GetDstAddress()
	fmt.Println(mac.String())
	typ := eth.GetEthernetType()
	fmt.Println(typ)
	require.Equal(net.HardwareAddr{0x94, 0x94, 0x26, 0x01, 0x02, 0x03}, eth.GetDstAddress())
	require.Equal(net.HardwareAddr{0x04, 0x05, 0x06, 0x07, 0x08, 0x09}, eth.GetSrcAddress())
	require.Equal(EthernetTypeIPv4, eth.GetEthernetType())
}

func TestEthernet_SetAll(t *testing.T) {
	require := require.New(t)
	p := make([]byte, 14)

	eth := *(*Ethernet)(&p)
	eth.SetSrcAddress(net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	eth.SetDstAddress(net.HardwareAddr{0x06, 0x05, 0x04, 0x03, 0x02, 0x01})
	eth.SetEthernetType(uint16(EthernetTypeIPv4))

	fmt.Printf("%x\n", eth)
	require.Equal(net.HardwareAddr{0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, eth.GetDstAddress())
	require.Equal(net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, eth.GetSrcAddress())
	require.Equal(EthernetTypeIPv4, eth.GetEthernetType())
}

func BenchmarkEthernetSet(b *testing.B) {
	p := make([]byte, 14)
	eth := *(*Ethernet)(&p)
	for i := 0; i < b.N; i++ {
		eth.SetSrcAddress(net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
		eth.SetDstAddress(net.HardwareAddr{0x06, 0x05, 0x04, 0x03, 0x02, 0x01})
		eth.SetEthernetType(uint16(EthernetTypeIPv4))
	}
}

func BenchmarkEthernetGet(b *testing.B) {
	p := []byte{
		0x94, 0x94, 0x26, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x08, 0x00,
	}

	eth := *(*Ethernet)(&p)
	for i := 0; i < b.N; i++ {
		eth.GetSrcAddress()
		eth.GetDstAddress()
		eth.GetEthernetType()
	}
}
