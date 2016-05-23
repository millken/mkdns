package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/millken/mkdns/drivers"
	"github.com/millken/mkdns/stats"
)

const BPFFilter = "udp and dst port 53"

type server struct {
	config           *Config
	io               drivers.PacketDataSourceCloser
	stats            stats.Directional
	childStoppedChan chan bool
	forceQuitChan    chan os.Signal
	txChan           chan PacketLayer
	rxChan           chan gopacket.Packet
	//handler          *Handler
	//rTimeout time.Duration
	//wTimeout time.Duration
}

func NewServer(config *Config) *server {
	return &server{
		config:           config,
		forceQuitChan:    make(chan os.Signal, 1),
		childStoppedChan: make(chan bool, 0),
		rxChan:           make(chan gopacket.Packet),
	}
}

func (s *server) Start() (err error) {
	options := &drivers.DriverOptions{
		Device:  "enp3s0",
		Snaplen: 2048,
		Filter:  BPFFilter,
	}

	factory, ok := drivers.Drivers["libpcap"]
	if !ok {
		log.Fatal(fmt.Sprintf("%s Packet driver not supported on this system", s.config.Server.Driver))
	}

	s.io, err = factory(options)
	if err != nil {
		return
	}
	for i := 0; i < 8; i++ {
		go packetHandler(i, s.rxChan)
	}
	go s.readPackets()
	go s.sendPackets()
	return
}

func (s *server) sendPackets() {
	for {
		//p := (<-s.txChan)
		//s.io.WritePacketData(p.Data())
		s.stats.Tx.Packets++
		//d.stats.Tx.Bytes += uint64(p.Metadata().CaptureInfo.CaptureLength)
	}
}

func (s *server) readPackets() {
	packetSource := s.io.PacketSource()
	for true {

		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		s.stats.Rx.Packets++

		s.rxChan <- packet
		//handlePacket(packet)
	}
}
