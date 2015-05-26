package plugins

import (
	"log"
	"net"

	"github.com/miekg/dns"
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

func (this *RecordNSPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	return this.NormalRecord(conf["record"].([]interface{}))
	return
}

func (this *RecordNSPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	for _, v := range records {
		log.Printf("ns=%s", v.(string))
		answer = append(answer, &dns.NS{this.RRheader, dns.Fqdn(v.(string))})
	}
	return
}

func init() {
	RegisterPlugin("NS", dns.TypeNS, func() interface{} {
		return new(RecordNSPlugin)
	})
}
