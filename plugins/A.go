package plugins

import (
	"log"
	"net"
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
	var vv map[string]interface{}
	var vvv []interface{}
	for _, v := range records {
		switch vt := v.(type) {
		case map[string]interface{}:
			vv = v.(map[string]interface{})
		default:
			log.Printf("[WARING] records struct not an map[string]interface{} : %v", vt)
		}
		if _, ok = vv["record"]; !ok {
			log.Printf("[WARING] record key not exit")
			continue
		}
		switch vt := vv["record"].(type) {
		case []interface{}:
			vvv = vv["record"].([]interface{})
		default:
			log.Printf("[WARING] records value not an list : %s", vt)
		}
		for _, vvvv := range vvv {
			ip := net.ParseIP(strings.TrimSpace(vvvv.(string)))
			if ip == nil {
				log.Printf("[ERROR] %s is not ipv4", strings.TrimSpace(vvvv.(string)))
				continue
			}
			answer = append(answer, &dns.A{
				Hdr: this.RRheader,
				A:   ip})
		}
	}
	return
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
