package driver

import (
	"time"

	"github.com/google/gopacket"
)

type DriverOptions struct {
	DAQ          string
	Filename     string
	Device       string
	Snaplen      int32
	WireDuration time.Duration
	Filter       string
}

// PacketDataSource is an interface for some source of packet data.
type PacketDataSourceCloser interface {
	// ReadPacketData returns the next packet available from this data source.
	// It returns:
	//  data:  The bytes of an individual packet.
	//  ci:  Metadata about the capture
	//  err:  An error encountered while reading packet data.  If err != nil,
	//    then data/ci will be ignored.
	ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error)
	// Close closes the ethernet sniffer and returns nil if no error was found.
	Close() error
}
