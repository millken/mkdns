package plugins

import (
	"net"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordTXTPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordTXTPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordTXTPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			answer = append(answer, &dns.TXT{
				Hdr: this.RRheader,
				Txt: []string{v},
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("TXT", dns.TypeTXT, func() interface{} {
		return new(RecordTXTPlugin)
	})
}
