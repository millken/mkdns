package plugins

import (
	"log"
	"net"
	//"fmt"
	"strings"
	//"reflect"

	"github.com/miekg/dns"
)

type RecordMXPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
	Conf     map[string]interface{}
}

func (this *RecordMXPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordMXPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
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

func (this *RecordMXPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	var ok bool
	var vv, vvvvv map[string]interface{}
	var vvv []interface{}
	var pref uint16
	for _, v := range records {
		switch vt := v.(type) {
		case map[string]interface{}:
			vv = v.(map[string]interface{})
		default:
			log.Printf("records value error, type= %v", vt)
		}
		if _, ok = vv["record"]; !ok {
			log.Printf("record not ok")
			continue
		}
		switch vt := vv["record"].(type) {
		case []interface{}:
			vvv = vv["record"].([]interface{})
		default:
			log.Printf("records value error, type= %v", vt)
		}
		for _, vvvv := range vvv {
			switch vt := vvvv.(type) {
			case map[string]interface{}:
				vvvvv = vvvv.(map[string]interface{})
			default:
				log.Printf("record error, type=%v", vt)
			}
			if _, ok = vvvvv["value"]; !ok {
				log.Printf("value not ok")
				continue
			}
			value := vvvvv["value"].(string)
			if !strings.HasSuffix(value, ".") {
				value = value + "."
			}
			if _, ok = vvvvv["mx"]; !ok {
				pref = 5
			} else {
				pref = uint16(vvvvv["mx"].(uint64))
			}

			dns_RR := &dns.MX{
				Hdr:        this.RRheader,
				Mx:         value,
				Preference: pref,
			}
			answer = append(answer, dns_RR)
		}
	}
	return
}

func init() {
	RegisterPlugin("MX", dns.TypeMX, func() interface{} {
		return new(RecordMXPlugin)
	})
}
