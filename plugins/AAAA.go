package plugins

import (
	"net"
	//"fmt"
	"strings"
	//"reflect"

	"github.com/miekg/dns"
	"github.com/millken/logger"
)

type RecordAAAAPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
	Conf     map[string]interface{}
}

func (this *RecordAAAAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordAAAAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
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

func (this *RecordAAAAPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	var ok bool
	var vv map[string]interface{}
	var vvv []interface{}
	for _, v := range records {
		switch vt := v.(type) {
		case map[string]interface{}:
			vv = v.(map[string]interface{})
		default:
			logger.Warn("records struct not an map[string]interface{} : %s", vt)
		}
		if _, ok = vv["record"]; !ok {
			logger.Warn("record key not exit")
			continue
		}
		switch vt := vv["record"].(type) {
		case []interface{}:
			vvv = vv["record"].([]interface{})
		default:
			logger.Warn("records value not an list : %s", vt)
		}
		for _, vvvv := range vvv {
			ip := net.ParseIP(strings.TrimSpace(vvvv.(string)))
			if ip == nil {
				logger.Error("%s is not ipv6", strings.TrimSpace(vvvv.(string)))
				continue
			}
			answer = append(answer, &dns.AAAA{this.RRheader, ip})
		}
	}
	return
}

func init() {
	RegisterPlugin("AAAA", dns.TypeAAAA, func() interface{} {
		return new(RecordAAAAPlugin)
	})
}
