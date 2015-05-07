package plugins

import (
	"net"

	"github.com/miekg/dns"
)

type RecordSOAPlugin struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
	RRheader   dns.RR_Header
}

func (this *RecordSOAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	this.EdnsAddr = edns
	this.RemoteAddr = remote
	this.RRheader = rr_header
}

func (this *RecordSOAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	answer = append(answer, &dns.SOA{
		Hdr:     this.RRheader,
		Ns:      conf["mname"].(string),
		Mbox:    conf["nname"].(string),
		Serial:  uint32(conf["serial"].(uint64)),
		Refresh: uint32(conf["refresh"].(uint64)),
		Retry:   uint32(conf["retry"].(uint64)),
		Expire:  uint32(conf["expire"].(uint64)),
		Minttl:  uint32(conf["minttl"].(uint64)),
	})
	return
}

func init() {
	RegisterPlugin("SOA", dns.TypeSOA, func() interface{} {
		return new(RecordSOAPlugin)
	})
}
