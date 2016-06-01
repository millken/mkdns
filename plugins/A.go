package plugins

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordAPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	records := getBaseRecord(state, this.Addr, rv)
	return this.NormalRecord(records)
}

func (this *RecordAPlugin) NormalRecord(rv []*types.Record_Value) (answer []dns.RR, err error) {
	var r []interface{}
	var e error
	for _, record := range records {
		r, e = getProofRecord(record)
		if e != nil {
			err = e
			continue
		}
		for _, v := range r {
			ip := net.ParseIP(strings.TrimSpace(v.(string)))
			if ip == nil {
				log.Printf("[ERROR] %s is not ipv4", strings.TrimSpace(v.(string)))
				continue
			}
			answer = append(answer, &dns.A{
				Hdr: this.RRheader,
				A:   ip})
		}
	}
	return
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
