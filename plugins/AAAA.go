package plugins

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type RecordAAAAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
	Conf     map[string]interface{}
}

func (this *RecordAAAAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordAAAAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	records := getBaseRecord(this.Addr, conf)
	return this.NormalRecord(records)
}

func (this *RecordAAAAPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
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
				log.Printf("[ERROR] %s is not ipv6", strings.TrimSpace(v.(string)))
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
