package plugins

import (
	"net"
	
	"github.com/miekg/dns"
	"github.com/millken/logger"
)

type RecordTXTPlugin struct {
	Addr     net.IP
	RRheader dns.RR_Header
	Conf     map[string]interface{}
}

func (this *RecordTXTPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	if edns != nil {
		this.Addr = edns
	} else {
		this.Addr = remote
	}

	this.RRheader = rr_header
}

func (this *RecordTXTPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	logger.Debug("conf : %+v", conf)
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

func (this *RecordTXTPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	var ok bool
	var vv map[string]interface{}
	var vvv string
	for _, v := range records {
		switch vt := v.(type) {
		case map[string]interface{}:
			vv = v.(map[string]interface{})
		default:
			logger.Warn("records struct not an map[string]interface{} : %v", vt)
		}
		if _, ok = vv["record"]; !ok {
			logger.Warn("record key not exit")
			continue
		}
		switch vt := vv["record"].(type) {
		case string:
			vvv = vv["record"].(string)
			answer = append(answer, &dns.TXT{this.RRheader, []string{vvv}})
		default:
			logger.Warn("records value not an list : %s", vt)
		}
	}
	return
}

func init() {
	RegisterPlugin("TXT", dns.TypeTXT, func() interface{} {
		return new(RecordTXTPlugin)
	})
}
