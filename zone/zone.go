package zone

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Zverushko/punycode"
	"github.com/miekg/dns"
	"github.com/millken/mkdns/plugins"
	"github.com/millken/mkdns/types"
	"github.com/umisama/go-regexpcache"
)

type Soa struct {
	Mname   string
	Nname   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minttl  uint32
}

type ZoneRecord struct {
	State int32
	Ttl   int
	Value []*types.Record_Value
}

type ZoneKey struct {
	Name string
	Type uint16
}

type ZoneRecords map[ZoneKey]*ZoneRecord

type ZoneOptions struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
}

type Zone struct {
	Name      string
	Records   ZoneRecords
	Regexp    ZoneRecords
	RegexpKey []ZoneKey
	Soa       dns.RR
	Ns        []dns.RR
	Options   ZoneOptions
}

func New() *Zone {
	z := new(Zone)
	//zone.Soa = []dns.RR
	z.Records = make(map[ZoneKey]*ZoneRecord)
	z.Regexp = make(map[ZoneKey]*ZoneRecord)
	z.Options.EdnsAddr = nil
	z.Options.RemoteAddr = nil
	return z
}

func (z *Zone) setSoaRR(ttl int32, rvs *types.Record_Value_Soa) {

	rr_header := dns.RR_Header{
		Name:   z.Name + ".",
		Rrtype: dns.TypeSOA,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}
	z.Soa = &dns.SOA{
		Hdr:     rr_header,
		Ns:      rvs.Mname,
		Mbox:    rvs.Nname,
		Serial:  rvs.Serial,
		Refresh: rvs.Refresh,
		Retry:   rvs.Retry,
		Expire:  rvs.Expire,
		Minttl:  rvs.Minttl,
	}
	//log.Printf("[FINE] zone : %s, SOA=%s", z.Name, z.Soa)
}

func (z *Zone) ParseRecords(rs types.Records) (err error) {
	var dtype uint16
	z.Name = rs.Domain
	hasRegexp := false
	for _, r := range rs.Records {
		dtype, err = plugins.DnsType(r.Type)
		if err != nil {
			return
		}
		if r.Name == "@" {
			if dtype == dns.TypeSOA {
				z.setSoaRR(r.Ttl, r.Value[0].Soa)
			} else if dtype == dns.TypeNS {
				z.setNsRR(r.Ttl, r.State, r.Value)
			}
		}
		name, err := punycode.ToASCII(r.Name)
		if err != nil {
			return fmt.Errorf("get punycode error: %s", err)
		}
		zk := ZoneKey{
			Name: name,
			Type: dtype,
		}
		if strings.Contains(name, "*") {
			hasRegexp = true
			zk.Name = strings.Replace(strings.Replace(zk.Name, "*", `\w+`, -1), ".", `\.`, -1) + "$"

			z.Regexp[zk] = &ZoneRecord{
				State: r.State,
				Ttl:   int(r.Ttl),
				Value: r.Value,
			}
		} else {
			z.Records[zk] = &ZoneRecord{
				State: r.State,
				Ttl:   int(r.Ttl),
				Value: r.Value,
			}
		}

	}
	if hasRegexp {
		for k, _ := range z.Regexp {
			z.RegexpKey = append(z.RegexpKey, k)
		}
		for i := 1; i < len(z.RegexpKey); i++ {
			for j := 0; j < len(z.RegexpKey)-i; j++ {
				iwL := strings.Count(z.RegexpKey[j].Name, `\w+`)
				idL := strings.Count(z.RegexpKey[j].Name, `\.`)
				jwL := strings.Count(z.RegexpKey[j+1].Name, `\w+`)
				jdL := strings.Count(z.RegexpKey[j+1].Name, `\.`)
				iL := iwL + idL*2 + (idL+1-iwL)*3
				jL := jwL + jdL*2 + (jdL+1-jwL)*3
				if iL < jL { //sort
					z.RegexpKey[j], z.RegexpKey[j+1] = z.RegexpKey[j+1], z.RegexpKey[j]
				}

			}
		}

	}
	return nil
}

func (z *Zone) setNsRR(ttl int32, state int32, rv []*types.Record_Value) {
	rr_header := dns.RR_Header{
		Name:   z.Name + ".",
		Rrtype: dns.TypeNS,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}
	plugin := plugins.Get(dns.TypeNS).(plugins.Plugin)
	if plugin == nil {
		log.Printf("[ERROR] plugin: %d not register", dns.TypeNS)
		return
	}
	plugin.New(z.Options.EdnsAddr, z.Options.RemoteAddr, rr_header)
	z.Ns, _ = plugin.Filter(state, rv)
}

func (z *Zone) SoaRR() dns.RR {
	return z.Soa
}

func (z *Zone) NsRR() []dns.RR {
	return z.Ns
}

func (z *Zone) FindRecord(req *dns.Msg) (m *dns.Msg, err error) {
	//var answer dns.RR
	var slab string
	var ok bool
	var zk ZoneKey
	record := new(ZoneRecord)
	q := req.Question[0]
	m = new(dns.Msg)
	m.SetReply(req)

	//rrtype := q.Qtype
	if len(q.Name) == len(z.Name)+1 {
		slab = "@"
	} else {
		tl := len(q.Name) - len(z.Name) - 2
		slab = strings.ToLower(q.Name[0:tl])
	}

	log.Printf("[DEBUG] z.Name=%s, q.Name=%s, slab=%s, q.Qtype=%d, z.Options=%+v", z.Name, q.Name, slab, q.Qtype, z.Options)
	zk = ZoneKey{Name: slab, Type: q.Qtype}
	if record, ok = z.Records[zk]; !ok {
		for _, zz := range z.RegexpKey {
			r := z.Regexp[zz]
			if zz.Type == q.Qtype && regexpcache.MustCompile(zz.Name).MatchString(slab) {
				record = r
				ok = true
				log.Printf("[DEBUG] hit regexp : [%s] %s", slab, zz.Name)
				break
			}
		}
	}
	if ok {
		rr_header := dns.RR_Header{
			Name:   q.Name,
			Rrtype: q.Qtype,
			Class:  dns.ClassINET,
			Ttl:    uint32(record.Ttl),
		}
		plugin := plugins.Get(q.Qtype).(plugins.Plugin)
		if plugin == nil {
			err = fmt.Errorf("plugin: %d not register", q.Qtype)
			return
		}
		plugin.New(z.Options.EdnsAddr, z.Options.RemoteAddr, rr_header)
		m.Answer, err = plugin.Filter(record.State, record.Value)

	} else {
		err = fmt.Errorf("record not foud :%s=>%s", dns.TypeToString[q.Qtype], q.Name)
	}
	return
}
