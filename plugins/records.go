package plugins

import (
	"log"
	"net"

	"github.com/millken/mkdns/types"
)

/*
 *
 type & 1  normal
 type & 2  weight
 type & 4  geo
*/

const (
	NORMAL   = 1
	WEIGHT   = 2
	GEO      = 4
)

var upChooseRecord = -1
var currentWeight int = 0

type BaseRecords struct {
	Addr        net.IP
	RecordValue []*types.Record_Value
	State       int32
}

func newBaseRecords(addr net.IP, state int32, rv []*types.Record_Value) *BaseRecords {
	return &BaseRecords{
		Addr:        addr,
		RecordValue: rv,
		State:       state,
	}
}

func (this *BaseRecords) GetRecords() (rrv []*types.Record_Value) {
	rrv = this.RecordValue
	if this.State&GEO == GEO {
		rrv = this.GeoRecord(rrv)
	}

	if this.State&WEIGHT == WEIGHT {
		rrv = this.WeightRecord(rrv)
	}
	return
}

func (this *BaseRecords) getMaxWeight(rv []*types.Record_Value) int {
	maxweight := 0
	for _, r := range rv {
		weight := int(r.Weight)
		if weight > maxweight {
			maxweight = weight
		}
	}
	return maxweight
}

func (this *BaseRecords) GeoRecord(rv []*types.Record_Value) (rrv []*types.Record_Value) {
	var default_rrv, country_rrv, continent_rrv []*types.Record_Value
	hitContinent := false
	hitCountry := false
	country, continent, netmask := geoIP.GetCountry(this.Addr)
	log.Printf("[FINE] geoip= %s, country= %s, continent=%s, netmask=%d", this.Addr, country, continent, netmask)

	for _, v := range rv {
		if v.Country == "" && v.Continent == "" {
			default_rrv = append(default_rrv, v)
		}
		if v.Country != "" {
			if v.Country == country {
				hitCountry = true
				country_rrv = append(country_rrv, v)
			}
		} else if v.Continent != "" && v.Continent == continent {
			hitContinent = true
			continent_rrv = append(continent_rrv, v)
		}
		//log.Printf("%+v, %+v", _country, _continent)
	}
	if hitCountry {
		rrv = country_rrv
	} else if hitContinent {
		rrv = continent_rrv
	} else {
		rrv = default_rrv
	}
	return
}

//http://www.wangshangyou.com/go/126.html
func (this *BaseRecords) WeightRecord(rv []*types.Record_Value) (rrv []*types.Record_Value) {
	rlen := len(rv)
	if rlen == 0 {
		return
	}

	maxweight := this.getMaxWeight(rv)
	for {
		upChooseRecord = (upChooseRecord + 1) % rlen
		if upChooseRecord == 0 {
			currentWeight = currentWeight - 2
			if currentWeight <= 0 {
				currentWeight = maxweight
			}
		}
		if weight := int(rv[upChooseRecord].Weight); weight >= currentWeight {
			//log.Printf("%+v", records[upChooseRecord])
			rrv = append(rrv, rv[upChooseRecord])
			break
		}
	}

	return
}

func getBaseRecord(state int32, addr net.IP, rv []*types.Record_Value) (rrv []*types.Record_Value) {

	br := newBaseRecords(addr, state, rv)
	rrv = br.GetRecords()
	return
}
