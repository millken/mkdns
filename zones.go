package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

type Vzones map[string]string

var (
	zones     map[string]*Zone
	zonesLock = new(sync.RWMutex)
)

func LoadZones(zfile string) {
	logger.Info("loading zone: %s", zfile)
	p, err := os.Open(zfile)
	if err != nil {
		logger.Exit("error opening zones file: %s", err)
	}
	contents, err := ioutil.ReadAll(p)
	if err != nil {
		logger.Exit("error reading zones file: %s", err)
	}
	blists := string(contents)
	blist := strings.Split(blists, "\n")
	reg := regexp.MustCompile(`\s+|\t+`)

	temp := make(map[string]*Zone)
	for _, b := range blist {
		b = reg.ReplaceAllString(b, " ")
		if len(b) < 1 {
			continue
		}

		alist := strings.Split(b, " ")
		zone := NewZone()
		zone.Name = alist[0]
		err = zone.LoadFile(alist[1])

		if err != nil {
			logger.Error(err)
			continue
		}
		temp[alist[0]] = zone

		logger.Debug("%s=%s", alist[0], alist[1])
	}
	zonesLock.Lock()
	zones = temp
	zonesLock.Unlock()
}

func GetZones() map[string]*Zone {
	zonesLock.RLock()
	defer zonesLock.RUnlock()
	return zones
}

func FindZoneByDomain(domain string) *Zone {
	zones := GetZones()
	darr := dns.SplitDomainName(domain)
	for i := len(darr) - 1; i >= 0; i-- {
		qarr := darr[i:]
		qkey := strings.Join(qarr, ".")
		if zone, ok := zones[qkey]; ok {
			return zone
		}
	}
	return nil
}
