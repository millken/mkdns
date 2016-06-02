package plugins

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordMXPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordMXPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordMXPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	pref := int32(5)
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {

		if r.Preference > 0 {
			pref = r.Preference
		}
		for _, v := range r.Record {
			if !strings.HasSuffix(v, ".") {
				v = v + "."
			}
			answer = append(answer, &dns.MX{
				Hdr:        this.RRheader,
				Mx:         v,
				Preference: uint16(pref),
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("MX", dns.TypeMX, func() interface{} {
		return new(RecordMXPlugin)
	})
}
