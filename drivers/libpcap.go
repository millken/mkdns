package drivers

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func init() {
	DriverRegister("libpcap", NewPcapHandle)
}

type PcapHandle struct {
	handle *pcap.Handle
}

func NewPcapHandle(options *DriverOptions) (PacketDataSourceCloser, error) {
	log.Printf("[INFO] %s %s", options.Device, options.Filter)
	pcapWireHandle, err := pcap.OpenLive(options.Device, options.Snaplen, true, pcap.BlockForever)
	pcapHandle := PcapHandle{
		handle: pcapWireHandle,
	}
	err = pcapHandle.handle.SetBPFFilter(options.Filter)
	return &pcapHandle, err
}

func (p *PcapHandle) PacketSource() *gopacket.PacketSource {
	return gopacket.NewPacketSource(p.handle, p.handle.LinkType())
}

func (p *PcapHandle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	return p.handle.ReadPacketData()
}

func (p *PcapHandle) WritePacketData(data []byte) (err error) {
	return p.handle.WritePacketData(data)
}

func (p *PcapHandle) Close() error {
	p.handle.Close()
	return nil
}
