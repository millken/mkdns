package layers

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacket(t *testing.T) {
	require := require.New(t)
	type IPv4Option struct {
		OptionType   uint8
		OptionLength uint8
		OptionData   []byte
	}
	type Ethernet struct {
		SrcMAC, DstMAC net.HardwareAddr
		EthernetType   uint16
	}
	type IPv4 struct {
		Version    uint8
		IHL        uint8
		TOS        uint8
		Length     uint16
		Id         uint16
		Flags      uint8
		FragOffset uint16
		TTL        uint8
		Protocol   uint8
		Checksum   uint16
		SrcIP      net.IP
		DstIP      net.IP
		Options    []IPv4Option
		Padding    []byte
	}
	type UDP struct {
		SrcPort, DstPort uint16
		Length           uint16
		Checksum         uint16
	}
	tests := []struct {
		Data     string
		Ethernet Ethernet
		IPv4     IPv4
		UDP      UDP
	}{
		{
			Data: `b025aa3ed5b33c9c0f17c8d408004500004d407200004011b7c1c0a8009bc0a80081cb90003500397d4c7b2401200001000000000001047465737403636f6d0000010001000029100000000000000c000a00088e547bff740acb63`,
			Ethernet: Ethernet{
				DstMAC:       net.HardwareAddr{0xb0, 0x25, 0xaa, 0x3e, 0xd5, 0xb3},
				SrcMAC:       net.HardwareAddr{0x3c, 0x9c, 0x0f, 0x17, 0xc8, 0xd4},
				EthernetType: 0x0800,
			},
			IPv4: IPv4{
				Version:    0x45 >> 4,
				IHL:        0x45 & 0x0f,
				TOS:        0x00,
				Length:     0x004d,
				Id:         0x4072,
				Flags:      0x00,
				FragOffset: 0x0000,
				TTL:        0x40,
				Protocol:   0x11,
				Checksum:   0xb7c1,
				SrcIP:      net.IP{0xc0, 0xa8, 0x00, 0x9b},
				DstIP:      net.IP{0xc0, 0xa8, 0x00, 0x81},
				Options:    []IPv4Option{},
				Padding:    []byte{},
			},
			UDP: UDP{
				SrcPort:  0xcb90,
				DstPort:  0x0035,
				Length:   0x0039,
				Checksum: 0x7d4c,
			},
		},
	}
	for _, tt := range tests {
		raw, err := hex.DecodeString(tt.Data)
		if err != nil {
			t.Fatal(err)
		}
		packet := &Packet{}
		err = ParsePacket(packet, raw, Default)
		require.NoError(err)
		require.Equal(tt.Ethernet.SrcMAC, packet.Ethernet.GetSrcAddress())
		require.Equal(tt.Ethernet.DstMAC, packet.Ethernet.GetDstAddress())
		require.Equal(tt.Ethernet.EthernetType, packet.Ethernet.GetEthernetType())
		require.Equal(tt.IPv4.Version, packet.IPv4.GetVersion())
		require.Equal(tt.IPv4.IHL, packet.IPv4.GetIHL())
		require.Equal(tt.IPv4.TOS, packet.IPv4.GetTOS())
		require.Equal(tt.IPv4.Length, packet.IPv4.GetLength())
		require.Equal(tt.IPv4.Id, packet.IPv4.GetID())
		require.Equal(tt.IPv4.Flags, packet.IPv4.GetFlags())
		require.Equal(tt.IPv4.FragOffset, packet.IPv4.GetFragOff())
		require.Equal(tt.IPv4.TTL, packet.IPv4.GetTTL())
		require.Equal(tt.IPv4.Protocol, packet.IPv4.GetProtocol())
		require.Equal(tt.IPv4.Checksum, packet.IPv4.GetChecksum())
		require.Equal(tt.IPv4.SrcIP, packet.IPv4.GetSrcAddr())
		require.Equal(tt.IPv4.DstIP, packet.IPv4.GetDstAddr())
		// require.Equal(tt.IPv4.Options, packet.IPv4.GetOptions())
		// require.Equal(tt.IPv4.Padding, packet.IPv4.GetPadding())
		require.Equal(tt.UDP.SrcPort, packet.UDP.GetSrcPort())
		require.Equal(tt.UDP.DstPort, packet.UDP.GetDstPort())
		require.Equal(tt.UDP.Length, packet.UDP.GetLength())
		require.Equal(tt.UDP.Checksum, packet.UDP.GetChecksum())
	}
}

func BenchmarkPacket(b *testing.B) {
	raw, err := hex.DecodeString(`b025aa3ed5b33c9c0f17c8d408004500004d407200004011b7c1c0a8009bc0a80081cb90003500397d4c7b2401200001000000000001047465737403636f6d0000010001000029100000000000000c000a00088e547bff740acb63`)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		packet := AcquirePacket()
		err := ParsePacket(packet, raw, Default)
		ReleasePacket(packet)
		if err != nil {
			b.Fatal(err)
		}
	}
}
