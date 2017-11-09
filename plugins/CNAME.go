package plugins

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordCNAMEPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordCNAMEPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordCNAMEPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			answer = append(answer, &dns.CNAME{
				Hdr:    this.RRheader,
				Target: strings.TrimSpace(v),
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("CNAME", dns.TypeCNAME, func() interface{} {
		return new(RecordCNAMEPlugin)
	})
}
