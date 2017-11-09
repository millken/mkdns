package plugins

import (
	"log"
	"net"
	"strconv"
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
	var priority int
	var mx string
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {

		for _, v := range r.Record {
			vv := strings.SplitN(v, " ", 2)
			if len(vv) != 2 {
				log.Printf("[ERROR] MX record format incorrect: %s", v)
				continue
			}
			priority, err = strconv.Atoi(vv[0])
			if err != nil {
				priority = 5
			}
			mx = strings.TrimSpace(vv[1])
			if !strings.HasSuffix(mx, ".") {
				mx = mx + "."
			}
			answer = append(answer, &dns.MX{
				Hdr:        this.RRheader,
				Mx:         mx,
				Preference: uint16(priority),
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
