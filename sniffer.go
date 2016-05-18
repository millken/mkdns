package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/millken/mkdns/drivers"
)

type TimedRawPacket struct {
	Timestamp time.Time
	RawPacket []byte
}

type Sniffer struct {
	options          *drivers.DriverOptions
	packetDataSource drivers.PacketDataSourceCloser
	isStopped        bool
	decodePacketChan chan TimedRawPacket
	stopDecodeChan   chan bool
}

type PacketManifest struct {
	Timestamp time.Time
	Flow      *TcpIpFlow
	RawPacket []byte
	IPv4      *layers.IPv4
	IPv6      *layers.IPv6
	TCP       layers.TCP
	Payload   gopacket.Payload
}

func NewSniffer(options *drivers.DriverOptions) *Sniffer {
	i := Sniffer{
		options:          options,
		decodePacketChan: make(chan TimedRawPacket),
		stopDecodeChan:   make(chan bool),
	}
	return &i
}

func (i *Sniffer) GetStartedChan() chan bool {
	return make(chan bool)
}

// Start... starts the TCP attack inquisition!
func (i *Sniffer) Start() {
	//i.setupHandle()

	go i.capturePackets()
	go i.decodePackets()
}

func (i *Sniffer) Stop() {
	log.Print("[INFO] sniffer: sending stopCapureChan signal")
	i.isStopped = true
	i.stopDecodeChan <- true
}

func (i *Sniffer) Close() {
	if i.packetDataSource != nil {
		log.Print("[INFO] closing packet capture socket")
		i.packetDataSource.Close()
	}
	log.Print("[INFO] stopping the sniffer decode loop")
	i.isStopped = true
	log.Print("[INFO] done.")
}

func (i *Sniffer) setupHandle() {
	var err error
	var what string

	factory, ok := drivers.Drivers[i.options.DAQ]
	if !ok {
		log.Fatal(fmt.Sprintf("%s Sniffer not supported on this system", i.options.DAQ))
	}
	i.packetDataSource, err = factory(i.options)

	if err != nil {
		panic(fmt.Sprintf("Failed to acquire DataAcQuisition source: %s", err))
	}

	what = fmt.Sprintf("interface %s", i.options.Device)

	log.Printf("[INFO] Starting %s packet capture on %s", i.options.DAQ, what)
}

func (i *Sniffer) capturePackets() {
	handle, err := pcap.OpenLive("enp3s0", 1024, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	handle.SetBPFFilter("udp")
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		headers, err := Parse(packet)
		if err != nil {
			log.Println(err)
		}
		b, err := json.MarshalIndent(headers, "", "  ")
		if err != nil {
			log.Println(err)
			// Skip packet if JSON marshalling errors
			continue
		}
		log.Printf("[DEBUG] %s", b)
		/*
			rawPacket, captureInfo, err := i.packetDataSource.ReadPacketData()
			if err == io.EOF {
				log.Print("ReadPacketData got EOF\n")
				i.Close()
				i.Stop()
				return
			}
			log.Printf("[DEBUG] ---- rawPacket ----\n%v\n", hex.Dump(rawPacket))
			if err != nil {
				log.Printf("packet capure read error: %s", err)
				continue
			}
			timedPacket := TimedRawPacket{
				Timestamp: captureInfo.Timestamp,
			}
			timedPacket.RawPacket = make([]byte, len(rawPacket))
			copy(timedPacket.RawPacket, rawPacket)
			i.decodePacketChan <- timedPacket
			if i.isStopped {
				break
			}
		*/
	}
}

func (i *Sniffer) decodePackets() {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var payload gopacket.Payload

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &payload)
	decoded := make([]gopacket.LayerType, 0, 4)

	for {
		select {
		case <-i.stopDecodeChan:
			return
		case timedRawPacket := <-i.decodePacketChan:
			newPayload := new(gopacket.Payload)
			payload = *newPayload
			err := parser.DecodeLayers(timedRawPacket.RawPacket, &decoded)
			if err != nil {
				continue
			}

			packetManifest := PacketManifest{
				Timestamp: timedRawPacket.Timestamp,
				Payload:   payload,
				IPv6:      nil,
				IPv4:      nil,
			}
			foundNetLayer := false

			for _, typ := range decoded {
				switch typ {
				case layers.LayerTypeIPv4:
					packetManifest.IPv4 = &ip4
					foundNetLayer = true
				case layers.LayerTypeIPv6:
					packetManifest.IPv6 = &ip6
					foundNetLayer = true
				case layers.LayerTypeTCP:
					if foundNetLayer {
						flow := TcpIpFlow{}
						if packetManifest.IPv6 == nil {
							// IPv4 case
							flow = NewTcpIpFlowFromFlows(ip4.NetworkFlow(), tcp.TransportFlow())
						} else if packetManifest.IPv4 == nil {
							// IPv6 case
							flow = NewTcpIpFlowFromFlows(ip6.NetworkFlow(), tcp.TransportFlow())
						} else {
							panic("wtf")
						}

						packetManifest.Flow = &flow
						packetManifest.TCP = tcp
						//i.dispatcher.ReceivePacket(&packetManifest)
					} else {
						log.Println("could not find IPv4 or IPv6 layer, inoring")
					}
				} // switch
			} // for

		} // select
	} // for
}
