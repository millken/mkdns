package plugins

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordAAAAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordAAAAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordAAAAPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			ip := net.ParseIP(strings.TrimSpace(v))
			if ip == nil {
				log.Printf("[ERROR] %s is not ipv6", strings.TrimSpace(v))
				continue
			}
			answer = append(answer, &dns.AAAA{
				Hdr:  this.RRheader,
				AAAA: ip})
		}
	}
	return
}

func init() {
	RegisterPlugin("AAAA", dns.TypeAAAA, func() interface{} {
		return new(RecordAAAAPlugin)
	})
}
