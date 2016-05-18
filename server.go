package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/millken/mkdns/drivers"
)

const BPFFilter = "udp and dst port 53"

type server struct {
	config           *Config
	sniffer          *Sniffer
	childStoppedChan chan bool
	forceQuitChan    chan os.Signal
	//handler  *Handler
	//rTimeout time.Duration
	//wTimeout time.Duration
}

func NewServer(config *Config) *server {
	snifferOption := &drivers.DriverOptions{
		DAQ:     "libpcap",
		Device:  "enp3s0",
		Snaplen: 1024,
		Filter:  BPFFilter,
	}
	sniffer := NewSniffer(snifferOption)
	return &server{
		config:           config,
		sniffer:          sniffer,
		forceQuitChan:    make(chan os.Signal, 1),
		childStoppedChan: make(chan bool, 0),
	}
}

func (s *server) Run() {
	s.sniffer.Start()
	signal.Notify(s.forceQuitChan, os.Interrupt)

	select {
	case <-s.forceQuitChan:
		log.Print("graceful shutdown: user force quit\n")
		log.Print("stopping sniffer")
		s.sniffer.Stop()
		log.Print("supervisor waiting for child to stop\n")
	case <-s.childStoppedChan:
		log.Print("graceful shutdown: packet-source stopped")
	}
}
