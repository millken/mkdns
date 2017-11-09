package plugins

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordSRVPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func (this *RecordSRVPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordSRVPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	var priority, weight, port int
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			vv := strings.SplitN(v, " ", 4)
			if len(vv) != 4 {
				log.Printf("[ERROR] SRV record format incorrect: %s", v)
				continue
			}
			priority, err = strconv.Atoi(vv[0])
			if err != nil {
				continue
			}
			weight, err = strconv.Atoi(vv[1])
			if err != nil {
				continue
			}
			port, err = strconv.Atoi(vv[2])
			if err != nil {
				continue
			}
			answer = append(answer, &dns.SRV{
				Hdr:      this.RRheader,
				Priority: uint16(priority),
				Weight:   uint16(weight),
				Port:     uint16(port),
				Target:   strings.TrimSpace(vv[3]),
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("SRV", dns.TypeSRV, func() interface{} {
		return new(RecordSRVPlugin)
	})
}
