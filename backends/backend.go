package backends

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/millken/mkdns/types"
	"github.com/streamrail/concurrent-map"
	"github.com/zvelo/ttlru"
)

type Backend interface {
	Load() error
	Watch()
}

var (
	backends  = map[string]func(*url.URL) (Backend, error){}
	zones     map[string]*types.Zone
	zonemap   = cmap.New()
	zonecache = ttlru.New(500, 600*time.Second)
)

func Open(rawUrl string) (backend Backend, err error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("backend parameter error : %s", err)
	}
	factory, ok := backends[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("%s backend not supported", u.Scheme)
	}

	backend, err = factory(u)
	if err != nil {
		return
	}
	return
}

func Register(name string, backend func(*url.URL) (Backend, error)) {
	if _, dup := backends[name]; dup {
		log.Fatal("duplicate backend", name)
	}
	backends[name] = backend
}

func GetZones() map[string]*types.Zone {
	return zones
}

func GetRecords(domain string) (record []*types.RecordPb, err error) {
	darr := dns.SplitDomainName(domain)
	for i := len(darr) - 1; i >= 0; i-- {
		qarr := darr[i:]
		qkey := strings.Join(qarr, ".")
		if tmp, ok := zonemap.Get(qkey); ok {
			record = tmp.([]*types.RecordPb)
			return
		}
	}
	err = fmt.Errorf("record not found")
	return
}
