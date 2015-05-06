package plugins

import (
	"log"
	"net"
	//"fmt"
	"strings"

	"github.com/miekg/dns"
)

/*
 *
 type & 1  view
 type & 2  weight
 type & 4  geo
*/

const (
	VIEW   = 1
	WEIGHT = 2
	GEO    = 4
)

var upChooseRecord = -1
var currentWeight int = 0

type RecordAPlugin struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
	RRheader   dns.RR_Header
	Conf       map[string]interface{}
}

func (this *RecordAPlugin) New(edns, remote net.IP, rr_header dns.RR_Header) {
	this.EdnsAddr = edns
	this.RemoteAddr = remote
	this.RRheader = rr_header
}

func (this *RecordAPlugin) Filter(conf map[string]interface{}) (answer []dns.RR, err error) {
	log.Printf("conf : %+v", conf)
	this.Conf = conf
	if _, ok := conf["type"]; !ok {
		return this.NormalRecord(this.Conf["records"].([]interface{}))
	}
	record_type := conf["type"].(uint64)
	if record_type&VIEW == VIEW {
		return this.ViewRecord()
	}
	if record_type&WEIGHT == WEIGHT {
		return this.WeightRecord()
	}
	if record_type&GEO == GEO {
		return this.GeoRecord()
	}
	return
}

func (this *RecordAPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	for _, v := range records {
		ip := net.ParseIP(strings.TrimSpace(v.(string)))
		if ip == nil {
			continue
		}
		answer = append(answer, &dns.A{this.RRheader, ip})
	}
	return
}

//http://www.wangshangyou.com/go/126.html
func (this *RecordAPlugin) WeightRecord() (answer []dns.RR, err error) {
	maxweight := this.getMaxWeight()
	//log.Printf("maxweight : %d", maxweight)
	records := this.Conf["records"].([]interface{})
	for {
		upChooseRecord = (upChooseRecord + 1) % len(records)
		if upChooseRecord == 0 {
			currentWeight = currentWeight - 2
			if currentWeight <= 0 {
				currentWeight = maxweight
			}
		}
		if weight := int(records[upChooseRecord].(map[string]interface{})["weight"].(uint64)); weight >= currentWeight {
			//log.Printf("%+v", records[upChooseRecord])
			return this.NormalRecord(records[upChooseRecord].(map[string]interface{})["record"].([]interface{}))
		}
	}

	return
}

func (this *RecordAPlugin) ViewRecord() (answer []dns.RR, err error) {
	return
}

func (this *RecordAPlugin) GeoRecord() (answer []dns.RR, err error) {
	var _country, _continent string
	var answer_records []interface{}
	country, continent, netmask := geoIP.GetCountry(this.EdnsAddr)
	log.Printf("geoip= %s, country= %s, continent=%s, netmask=%d", this.EdnsAddr, country, continent, netmask)
	records := this.Conf["records"].([]interface{})
	for _, v := range records {
		vv := v.(map[string]interface{})
		if _, ok := vv["country"]; ok {
			_country = vv["country"].(string)
		} else {
			_country = ""
		}
		if _, ok := vv["continent"]; ok {
			_continent = vv["continent"].(string)
		} else {
			_continent = ""
		}
		if _country != "" && _country == country {
			return this.NormalRecord(vv["record"].([]interface{}))
		}
		if _continent != "" && _continent == continent {
			answer_records = vv["record"].([]interface{})
		}

		//log.Printf("%+v, %+v", _country, _continent)
	}
	return this.NormalRecord(answer_records)
}

func (this *RecordAPlugin) getMaxWeight() int {
	maxweight := 0
	for _, v := range this.Conf["records"].([]interface{}) {
		weight := int(v.(map[string]interface{})["weight"].(uint64))
		if weight > maxweight {
			maxweight = weight
		}
	}
	return maxweight
}

func init() {
	RegisterPlugin("A", dns.TypeA, func() interface{} {
		return new(RecordAPlugin)
	})
}
