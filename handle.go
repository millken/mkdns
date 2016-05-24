package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/miekg/dns"
)

func packetHandler(i int, in <-chan gopacket.Packet, out chan PacketLayer) {
	for packet := range in {
		p, err := parsePacket(packet)
		if err != nil || p.dns == nil {
			log.Printf("[ERROR] parsePacket %s", err)
			continue
		}
		log.Printf("[FINE] worker : %d ,data:\n%s", i, p.dns)
		req := p.dns
		m := new(dns.Msg)
		if req != nil {
			m.SetReply(req)
			m.SetRcode(req, dns.RcodeNameError)
		}
		m.Authoritative = true
		if e := m.IsEdns0(); e != nil {
			m.SetEdns0(4096, e.Do())
		}
		p.dns = m
		out <- p
	}
}
