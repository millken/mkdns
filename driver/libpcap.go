package drivers

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func init() {
	SnifferRegister("libpcap", NewPcapSniffer)
}

type PcapHandle struct {
	handle *pcap.Handle
}

func NewPcapSniffer(options *DriverOptions) (PacketDataSourceCloser, error) {
	pcapWireHandle, err := pcap.OpenLive(options.Device, options.Snaplen, true, options.WireDuration)
	pcapHandle := PcapHandle{
		handle: pcapWireHandle,
	}
	err = pcapHandle.handle.SetBPFFilter(options.Filter)
	return &pcapHandle, err
}
func NewPcapWireSniffer(netDevice string, snaplen int32, wireDuration time.Duration, filter string) (*PcapHandle, error) {
	pcapWireHandle, err := pcap.OpenLive(netDevice, snaplen, true, wireDuration)
	pcapHandle := PcapHandle{
		handle: pcapWireHandle,
	}
	err = pcapHandle.handle.SetBPFFilter(filter)
	return &pcapHandle, err
}

func (p *PcapHandle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	return p.handle.ReadPacketData()
}

func (p *PcapHandle) Close() error {
	p.handle.Close()
	return nil
}
