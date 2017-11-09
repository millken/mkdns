package plugins

import (
	"net"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordSOAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordSOAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordSOAPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		answer = append(answer, &dns.SOA{
			Hdr:     this.RRheader,
			Ns:      r.Soa.Mname,
			Mbox:    r.Soa.Nname,
			Serial:  r.Soa.Serial,
			Refresh: r.Soa.Refresh,
			Retry:   r.Soa.Retry,
			Expire:  r.Soa.Expire,
			Minttl:  r.Soa.Minttl,
		})
	}
	return
}

func init() {
	RegisterPlugin("SOA", dns.TypeSOA, func() interface{} {
		return new(RecordSOAPlugin)
	})
}
