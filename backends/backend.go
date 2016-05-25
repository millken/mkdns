package backends

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/millken/mkdns/types"
)

type Backend interface {
	Load() error
	Watch()
}

var (
	backends      = map[string]func(*url.URL) (Backend, error){}
	zones         map[string]*types.Zone
	zonesLock     = new(sync.RWMutex)
	lastReadZones time.Time
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
	zonesLock.RLock()
	defer zonesLock.RUnlock()
	return zones
}
