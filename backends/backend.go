package backends

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/millken/mkdns/cache"
	"github.com/millken/mkdns/types"
	"github.com/millken/mkdns/zone"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"github.com/zvelo/ttlru"
)

type Backend interface {
	Load()
}

var (
	backends  = map[string]func(*url.URL) (Backend, error){}
	zonemap   = cache.NewZoneMap()
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

func GetRecords(domain string) (zz *zone.Zone, err error) {
	var dbp types.Records
	tldomain, err := publicsuffix.Domain(domain[0 : len(domain)-1])
	if err != nil {
		return
	}
	v, e := zonecache.Get(tldomain)
	if e {
		zz = v.(*zone.Zone)
		return
	}
	if tmp, ok := zonemap.Get(tldomain); ok {
		dbp, err = types.DecodeByProtobuff(tmp)
		if err != nil {
			return nil, fmt.Errorf("%s decode protobuf err %s", tldomain, err)
		}
		zz = zone.New()
		if dbp.Domain == "" {
			dbp.Domain = tldomain
		}
		if err = zz.ParseRecords(dbp); err == nil {
			zonecache.Set(tldomain, zz)
		}
		return
	}
	err = fmt.Errorf("%s record not found", tldomain)
	return
}
