package plugins

import (
	"net"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordSPFPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordSPFPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordSPFPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			answer = append(answer, &dns.SPF{
				Hdr: this.RRheader,
				Txt: []string{v},
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("SPF", dns.TypeSPF, func() interface{} {
		return new(RecordSPFPlugin)
	})
}
