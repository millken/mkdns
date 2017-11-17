package plugins

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
)

type RecordCAAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
}

func atou8(s string) uint8 {
	u64, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		log.Printf("[ERROR]CAA atou8 failed (%v) (err=%v", s, err)
		return 0
	}
	return uint8(u64)
}

func (this *RecordCAAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordCAAPlugin) Filter(state int32, rv []*types.Record_Value) (answer []dns.RR, err error) {
	rv = getBaseRecord(state, this.Addr, rv)
	for _, r := range rv {
		for _, v := range r.Record {
			t := strings.Fields(v)
			if len(t) != 3 {
				log.Printf("[ERROR] %s is not caa", strings.TrimSpace(v))
				continue
			}

			answer = append(answer, &dns.CAA{
				Hdr:   this.RRheader,
				Flag:  atou8(t[0]),
				Tag:   t[1],
				Value: t[2],
			})
		}
	}
	return
}

func init() {
	RegisterPlugin("CAA", dns.TypeCAA, func() interface{} {
		return new(RecordCAAPlugin)
	})
}
