package plugins

import (
	"net"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordNSPlugin struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
	RRheader   dns.RR_Header
}

func (this *RecordNSPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	this.EdnsAddr = edns
	this.RemoteAddr = remote
	this.RRheader = rr_header
}

func (this *RecordNSPlugin) Filter(rv []*types.Record_Value) (answer []dns.RR, err error) {
	return this.NormalRecord(rv)
}

func (this *RecordNSPlugin) NormalRecord(rv []*types.Record_Value) (answer []dns.RR, err error) {
	for _, r := range rv {
		for _, v := range r.Record {
			answer = append(answer, &dns.NS{
				Hdr: this.RRheader,
				Ns:  dns.Fqdn(v),
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("NS", dns.TypeNS, func() interface{} {
		return new(RecordNSPlugin)
	})
}
