package plugins

import (
	"log"
	"net"
	//"fmt"
	"strings"

	"github.com/miekg/dns"
)

type RecordAPlugin bool

func (this *RecordAPlugin) Filter(rr_header dns.RR_Header, conf map[string]interface{}) (answer []dns.RR, err error) {
	//log.Printf("conf : %+v, %T", conf, conf["value"])
	if _, ok := conf["value"]; !ok {
		return
	}
	for _, v := range conf["value"].([]interface{}) {
		ip := net.ParseIP(strings.TrimSpace(v.(string)))
		if ip == nil {
			continue
		}
		answer = append(answer, &dns.A{rr_header, ip})
	}
	return answer, nil
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
