package main

import (
	//"logger"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/logger"
	"github.com/rcrowley/go-metrics"
)

type Handler struct {
	responseTimer metrics.Timer
}

func (h *Handler) query(Net string, w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true
	if e := m.IsEdns0(); e != nil {
		m.SetEdns0(4096, e.Do())
	}
	q := req.Question[0]
	domain := strings.ToLower(q.Name)
	zone := FindZoneByDomain(domain)

	var zone_name string
	if zone == nil {
		zone_name = "NULL"
	} else {
		zone_name = zone.Name
	}
	log.Printf("[zone %s] incoming %s %s %d from %s://%s", zone_name, req.Question[0].Name,
		dns.TypeToString[q.Qtype], req.MsgHdr.Id, Net, w.RemoteAddr())

	if zone == nil {
		m.SetRcode(req, dns.RcodeNameError)
		w.WriteMsg(m)
		return
	}

	realIp, _, _ := net.SplitHostPort(w.RemoteAddr().String())
	zone.Options.EdnsAddr = nil
	zone.Options.RemoteAddr = net.ParseIP(realIp)

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
					log.Printf("Got edns %s %s %d %s", e.Address, e.Family, e.SourceNetmask, e.SourceScope)
					if e.Address != nil {
						//edns = e
						zone.Options.EdnsAddr = e.Address
					}
				}
			}
		}
	}
	/*
		// TODO: set scope to 0 if there are no alternate responses
		if edns != nil {
			if edns.Family != 0 {
				if netmask < 16 {
					netmask = 16
				}
				edns.SourceScope = uint8(netmask)
				m.Extra = append(m.Extra, opt_rr)
			}
		}
	*/
	defer func() {
		if err := w.WriteMsg(m); err != nil {
			logger.Error("failure to return reply %s", err.Error())
		}
		return
	}()
	if q.Qclass == dns.ClassCHAOS {
		if q.Qtype == dns.TypeTXT {
			switch domain {
			case "bind.version":
				fallthrough
			case "id.server.":
				// TODO(miek): machine name to return
				hdr := dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
				m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{"localhost"}}}
				return
			}

		}
	}
	m, err := zone.FindRecord(req)
	if err != nil {
		m.Ns = append(m.Ns, zone.SoaRR())
		logger.Error("zone error : %s", err)
	} else {
		m.Ns = zone.NsRR()
	}
	m.Authoritative = true

	//loggerger.Debug("%s", m)

	w.WriteMsg(m)
	return
}
func (h *Handler) TCP(w dns.ResponseWriter, req *dns.Msg) {
	h.query("tcp", w, req)
}
func (h *Handler) UDP(w dns.ResponseWriter, req *dns.Msg) {
	h.query("udp", w, req)
}
