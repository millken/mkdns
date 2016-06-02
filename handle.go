package main

import (
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/miekg/dns"
	"github.com/millken/mkdns/backends"
)

func packetHandler(i int, in <-chan gopacket.Packet, out chan PacketLayer) {
	for packet := range in {
		p, err := parsePacket(packet)
		if err != nil || p.dns == nil {
			log.Printf("[ERROR] parsePacket %s", err)
			continue
		}
		req := p.dns
		q := req.Question[0]
		domain := strings.ToLower(q.Name)
		zz, err := backends.GetRecords(domain)
		if err != nil {
			log.Printf("[WARN] %s %s", domain, err)
		}
		log.Printf("zz : %+v", zz)

		m := new(dns.Msg)
		if req != nil {
			m.SetReply(req)
			m.SetRcode(req, dns.RcodeNameError)
		}
		if e := m.IsEdns0(); e != nil {
			m.SetEdns0(4096, e.Do())
		}

		var zname string
		if zz == nil {
			zname = "NULL"
		} else {
			zname = zz.Name
		}
		log.Printf("[FINE] [zone %s] incoming %s %s %d from %s", zname, req.Question[0].Name,
			dns.TypeToString[q.Qtype], req.MsgHdr.Id, p.ipv4.SrcIP)

		if zz == nil {
			m.SetRcode(req, dns.RcodeNameError)
		} else {

			zz.Options.EdnsAddr = nil
			zz.Options.RemoteAddr = p.ipv4.SrcIP

			//var edns *dns.EDNS0_SUBNET
			//var opt_rr *dns.OPT

			for _, extra := range req.Extra {

				switch extra.(type) {
				case *dns.OPT:
					for _, o := range extra.(*dns.OPT).Option {
						//opt_rr = extra.(*dns.OPT)
						switch e := o.(type) {
						case *dns.EDNS0_NSID:
						case *dns.EDNS0_SUBNET:
							//log.Printf("[DEBUG] Got edns", e.Address, e.Family, e.SourceNetmask, e.SourceScope)
							if e.Address != nil {
								//edns = e
								zz.Options.EdnsAddr = e.Address
							}
						}
					}
				}
			}

			if q.Qclass == dns.ClassCHAOS {
				if q.Qtype == dns.TypeTXT {
					switch domain {
					case "bind.version":
						fallthrough
					case "id.server.":
						hdr := dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
						m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{"localhost"}}}
						return
					}
				}
			}
			m, err = zz.FindRecord(req)
			if err != nil {
				m.Ns = append(m.Ns, zz.SoaRR())
				log.Printf("[ERROR] zone error : %s", err)
			} else {
				m.Ns = zz.NsRR()
			}
		}
		m.Authoritative = true
		p.dns = m
		out <- p
	}
}
