package main

import (
	"log"

	"github.com/google/gopacket"
)

func packetHandler(i int, in <-chan gopacket.Packet) {
	for packet := range in {
		headers, err := parsePacket(packet)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[DEBUG] headers :%d => %s", i, headers.dns)
	}
}
