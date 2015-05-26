package plugins

import (
	"net"
	//"fmt"
	"strings"

	"github.com/miekg/dns"
)

type RecordAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
	Conf     map[string]interface{}
}

func (this *RecordAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	//log.Printf("conf : %+v", conf)
	var records []interface{}
	var ok bool
	this.Conf = conf
	if _, ok = conf["type"]; !ok {
		if _, ok = this.Conf["records"]; ok {
			records = this.Conf["records"].([]interface{})
		}
	} else {
		records = this.Conf["records"].([]interface{})
		record_type := conf["type"].(uint64)
		br := newBaseRecords(this.Addr, record_type, records)
		records = br.GetRecords()
	}
	return this.NormalRecord(records)
}

func (this *RecordAPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	var ok bool
	for _, v := range records {
		if _, ok = v.(map[string]interface{})["record"]; !ok {
			continue
		}
		for _, vv := range v.(map[string]interface{})["record"].([]interface{}) {
			ip := net.ParseIP(strings.TrimSpace(vv.(string)))
			if ip == nil {
				continue
			}
			answer = append(answer, &dns.A{this.RRheader, ip})
		}
	}
	return
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
