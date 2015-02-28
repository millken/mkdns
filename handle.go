package main

import (
	"github.com/miekg/dns"
	"github.com/rcrowley/go-metrics"
	"strings"
)

type Handler struct {
	responseTimer metrics.Timer
}

func (h *Handler) UDP(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true
	if e := m.IsEdns0(); e != nil {
		m.SetEdns0(4096, e.Do())
	}
	q := req.Question[0]
	domain := strings.ToLower(q.Name)
	zone := FindZoneByDomain(domain)
	if zone == nil {
		m.SetRcode(req, dns.RcodeNameError)
		return
	}

	defer func() {
		logger.Debug("(domain)=%s, (q)=%v\"%s\"", domain, q, q.String())
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
	switch q.Qtype {
	case dns.TypeNS:
	case dns.TypeA, dns.TypeAAAA:
		m , _ = zone.FindRecord(req)
		return
	default:
		fallthrough
	case dns.TypeSRV, dns.TypeANY:
	}

	m.SetRcode(req, dns.RcodeNameError)
	w.WriteMsg(m)
	return
}
