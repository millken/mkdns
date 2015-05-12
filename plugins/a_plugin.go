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
	//log.Printf("conf : %+v", conf)
	var records []interface{}
	var ok bool
	this.Conf = conf
	if _, ok = conf["type"]; !ok {
		if _, ok = this.Conf["records"]; ok {
			records = this.Conf["records"].([]interface{})
		}
	}else{
	records = this.Conf["records"].([]interface{})
	record_type := conf["type"].(uint64)
	if record_type&GEO == GEO {
		records = this.GeoRecord(records)
	}
	if record_type&VIEW == VIEW {
		return this.ViewRecord()
	}
	if record_type&WEIGHT == WEIGHT {
		records = this.WeightRecord(records)
	}
}
	return this.NormalRecord(records)
}

func (this *RecordAPlugin) NormalRecord(records []interface{}) (answer []dns.RR, err error) {
	var ok bool
	for _, v := range records {
		if _, ok = v.(map[string]interface{})["record"]; !ok {
			continue
		}
		for  _, vv := range v.(map[string]interface{})["record"].([]interface{}) {
		ip := net.ParseIP(strings.TrimSpace(vv.(string)))
		if ip == nil {
			continue
		}
		answer = append(answer, &dns.A{this.RRheader, ip})
		}
	}
	return
}

//http://www.wangshangyou.com/go/126.html
func (this *RecordAPlugin) WeightRecord(records []interface{}) (answer []interface{}) {
	var ok bool
	var w uint64
	rlen := len(records)
	if rlen == 0 {
		return
	}
	maxweight := this.getMaxWeight()
	//log.Printf("maxweight : %d", maxweight)
	for {
		upChooseRecord = (upChooseRecord + 1) % rlen
		if upChooseRecord == 0 {
			currentWeight = currentWeight - 2
			if currentWeight <= 0 {
				currentWeight = maxweight
			}
		}
		if _, ok = records[upChooseRecord].(map[string]interface{})["weight"]; !ok {
			w = 0
		}else{
			w = records[upChooseRecord].(map[string]interface{})["weight"].(uint64)
		}		
		if weight := int(w); weight >= currentWeight {
			//log.Printf("%+v", records[upChooseRecord])
			answer = append(answer, records[upChooseRecord])
			break
		}
	}

	return
}

func (this *RecordAPlugin) ViewRecord() (answer []dns.RR, err error) {
	return
}

func (this *RecordAPlugin) GeoRecord(records []interface{}) (answer []interface{}) {
	var _country, _continent string
	var remote_addr net.IP
	var def_answer []interface{}
	hitGeo := false
	if this.EdnsAddr != nil {
		remote_addr = this.EdnsAddr
	}else{
		remote_addr = this.RemoteAddr
	}
	country, continent, netmask := geoIP.GetCountry(remote_addr)
	log.Printf("geoip= %s, country= %s, continent=%s, netmask=%d", remote_addr, country, continent, netmask)

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
		if _country == "" && _continent == "" {
			def_answer = append(def_answer, v)
		}
		if (_country != "" && _country == country) || (_continent != "" && _continent == continent) {
			hitGeo = true
			answer = append(answer, v)
		}
		//log.Printf("%+v, %+v", _country, _continent)
	}
	if !hitGeo {
		answer = def_answer
	}
	return 
}

func (this *RecordAPlugin) getMaxWeight() int {
	var ok bool
	var w uint64
	maxweight := 0
	for _, v := range this.Conf["records"].([]interface{}) {
		if _, ok = v.(map[string]interface{})["weight"]; !ok {
			w = 0
		}else{
			w = v.(map[string]interface{})["weight"].(uint64)
		}
		weight := int(w)
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
