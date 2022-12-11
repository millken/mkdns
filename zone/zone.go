package zone

import (
	"net"
	"sort"
	"strings"
	"sync"

	"github.com/millken/mkdns/dns"
)

type Key struct {
	Suffix string
	Type   dns.Type
}

type KeyList []Key

func (k KeyList) Len() int {
	return len(k)
}

//["*.a.com","*.2.a.com","*.1.2.6.a.com","*.6.a.com"]

func (k KeyList) Less(i, j int) bool {
	return strings.Count(k[i].Suffix, ".") < strings.Count(k[j].Suffix, ".")
}

func (k KeyList) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

type Options struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
}

type Zone struct {
	name    string
	mapG    map[Key]*dns.Record
	mapW    map[Key]*dns.Record
	keylist KeyList
	soa     *dns.SOA // cached soa
	ns      []net.NS // cached ns
	rw      sync.RWMutex
	opt     Options
}

func New(name string) *Zone {
	z := new(Zone)
	z.name = name
	z.mapG = make(map[Key]*dns.Record)
	z.mapW = make(map[Key]*dns.Record)
	z.opt.EdnsAddr = nil
	z.opt.RemoteAddr = nil
	return z
}

func (z *Zone) Add(suffix string, record *dns.Record) {
	z.rw.Lock()
	defer z.rw.Unlock()
	z.add(Key{suffix, record.Type}, record)
}
func (z *Zone) add(zkey Key, record *dns.Record) {
	if strings.IndexByte(zkey.Suffix, '*') >= 0 {
		z.mapW[zkey] = record
		z.keylist = append(z.keylist, zkey)
		sort.Sort(z.keylist)
	} else {
		z.mapG[zkey] = record
	}
}

func (z *Zone) Lookup(key Key) (record *dns.Record, found bool) {
	z.rw.RLock()
	defer z.rw.RUnlock()
	return z.lookup(key)
}

func (z *Zone) lookup(key Key) (record *dns.Record, found bool) {
	record, found = z.mapG[key]
	if found {
		return
	}
	for wk, wr := range z.mapW {
		if wk.Type == key.Type {
			wkSuf := wk.Suffix[1:]
			wkSufLen := len(wkSuf)
			kSufLen := len(key.Suffix)
			if wkSufLen >= kSufLen {
				continue
			}

			if strings.Count(wk.Suffix, ".") > strings.Count(key.Suffix, ".") {
				break
			}
			if wk.Suffix[1:] == key.Suffix[kSufLen-wkSufLen:] {
				record = wr
				found = true
			}
		}
	}
	return
}
