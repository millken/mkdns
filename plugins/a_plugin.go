package plugins

import (
	"log"
	"net"
	//"fmt"
	"strings"

	"github.com/miekg/dns"
)

type RecordAPlugin struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
	RRheader   dns.RR_Header
}

func (this *RecordAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	this.EdnsAddr = edns
	this.RemoteAddr = remote
	this.RRheader = rr_header
}
func (this *RecordAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	log.Printf("conf : %+v, %T", conf, conf["value"])
	if _, ok := conf["value"]; !ok {
		return
	}
	for _, v := range conf["value"].([]interface{}) {
		ip := net.ParseIP(strings.TrimSpace(v.(string)))
		if ip == nil {
			continue
		}
		answer = append(answer, &dns.A{this.RRheader, ip})
	}
	return answer, nil
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
