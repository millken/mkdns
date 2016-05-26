package plugins

import (
	"fmt"
	"log"
	"net"
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

type BaseRecords struct {
	Addr    net.IP
	Records []interface{}
	RType   uint64
}

func newBaseRecords(addr net.IP, rtype uint64, records []interface{}) *BaseRecords {
	return &BaseRecords{
		Addr:    addr,
		Records: records,
		RType:   rtype,
	}
}

func (this *BaseRecords) GetRecords() (answer []interface{}) {
	var records []interface{}
	records = this.Records
	if this.RType&GEO == GEO {
		records = this.GeoRecord(records)
	}

	if this.RType&WEIGHT == WEIGHT {
		records = this.WeightRecord(records)
	}
	return records
}

func (this *BaseRecords) getMaxWeight(records []interface{}) int {
	var ok bool
	var w uint64
	maxweight := 0
	for _, v := range records {
		if _, ok = v.(map[string]interface{})["weight"]; !ok {
			w = 0
		} else {
			w = v.(map[string]interface{})["weight"].(uint64)
		}
		weight := int(w)
		if weight > maxweight {
			maxweight = weight
		}
	}
	return maxweight
}

func (this *BaseRecords) GeoRecord(records []interface{}) (answer []interface{}) {
	var _country, _continent string
	var default_answer, country_answer, continent_answer []interface{}
	hitContinent := false
	hitCountry := false
	country, continent, netmask := geoIP.GetCountry(this.Addr)
	log.Printf("[FINE] geoip= %s, country= %s, continent=%s, netmask=%d", this.Addr, country, continent, netmask)

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
			default_answer = append(default_answer, v)
		}
		if _country != "" {
			if _country == country {
				hitCountry = true
				country_answer = append(country_answer, v)
			}
		} else if _continent != "" && _continent == continent {
			hitContinent = true
			continent_answer = append(continent_answer, v)
		}
		//log.Printf("%+v, %+v", _country, _continent)
	}
	if hitCountry {
		answer = country_answer
	} else if hitContinent {
		answer = continent_answer
	} else {
		answer = default_answer
	}
	return
}

//http://www.wangshangyou.com/go/126.html
func (this *BaseRecords) WeightRecord(records []interface{}) (answer []interface{}) {
	var ok bool
	var w uint64
	rlen := len(records)
	if rlen == 0 {
		return
	}

	maxweight := this.getMaxWeight(records)
	log.Printf("[FINE] maxweight : %d ", maxweight)
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
		} else {
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

func getBaseRecord(addr net.IP, cf map[string]interface{}) (records []interface{}) {
	var ok bool
	if _, ok = cf["type"]; !ok {
		if _, ok = cf["records"]; ok {
			records = cf["records"].([]interface{})
		}
	} else {
		records = cf["records"].([]interface{})
		record_type := cf["type"].(uint64)
		br := newBaseRecords(addr, record_type, records)
		records = br.GetRecords()
	}
	return
}

func getProofRecord(record interface{}) (result []interface{}, err error) {
	var rv map[string]interface{}

	switch rt := record.(type) {
	case map[string]interface{}:
		rv = record.(map[string]interface{})
	default:
		return nil, fmt.Errorf("records struct not an map[string]interface{} : %v", rt)
	}
	if _, ok := rv["record"]; !ok {
		return nil, fmt.Errorf("record not exit : %v", rv)
	}
	switch vt := rv["record"].(type) {
	case []interface{}:
		result = rv["record"].([]interface{})
	default:
		return nil, fmt.Errorf("records value not an list : %s", vt)
	}
	return
}
