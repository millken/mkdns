package types

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/Zverushko/punycode"
	"github.com/miekg/dns"
	"github.com/millken/mkdns/plugins"
	"github.com/ugorji/go/codec"
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

type OrigRecord struct {
	Label string
	Ttl   int
	Type  uint16
	Value map[string]interface{}
}
type Record struct {
	Ttl   int
	Value map[string]interface{}
}

type ZoneKey struct {
	Label string
	Type  uint16
}

type Records map[ZoneKey]*Record

type ZoneOptions struct {
	EdnsAddr   net.IP
	RemoteAddr net.IP
}

type Zone struct {
	Name    string
	Records Records
	Regexp  Records
	Soa     dns.RR
	Ns      []dns.RR
	Options ZoneOptions
}

func NewZone() *Zone {
	zone := new(Zone)
	//zone.Soa = []dns.RR
	zone.Records = make(map[ZoneKey]*Record)
	zone.Regexp = make(map[ZoneKey]*Record)
	zone.Options.EdnsAddr = nil
	zone.Options.RemoteAddr = nil
	return zone
}

func (z *Zone) ParseBody(data []byte) (err error) {
	records := make(Records)
	regexp_records := make(Records)
	line := 0

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanBytes)
	a := ""
	ahead := ""
	b := 1.0
	aheadspace := false

	for scanner.Scan() {
		//log.Printf("[DEBUG] %v %s", scanner.Bytes(), scanner.Text())
		c := scanner.Text()
		//log.Printf("[DEBUG] c=%s, unicode.IsSpace([]rune(c)[0])=%v", c, unicode.IsSpace([]rune(c)[0]))
		if c == "\n" && int(b) == 1 {
			line = line + 1
			text := strings.TrimSpace(a)

			a = ""
			if len(text) == 0 || text[0:1] == ";" {
				continue
			}
			r, rerr := z.parseLine(line, text)
			if rerr != nil {
				log.Printf("[WARN] %d : %q", line, rerr)
				continue
			}
			zk := ZoneKey{
				Label: r.Label,
				Type:  r.Type,
			}
			if strings.Contains(r.Label, "*") {
				regexp_records[zk] = &Record{
					Ttl:   r.Ttl,
					Value: r.Value,
				}
			} else {
				records[zk] = &Record{
					Ttl:   r.Ttl,
					Value: r.Value,
				}
			}

		}

		if c == "\"" && ahead != "\\" {
			b = math.Abs(b - 1)
		}

		if unicode.IsSpace([]rune(c)[0]) {
			if c == " " && (ahead == " " || aheadspace == true) {
				c = ""
			} else
			//hack "
			if int(b) == 0 {
				c = "\t"
			} else {
				c = " "
			}
			aheadspace = true
		} else {
			aheadspace = false
		}
		ahead = c
		a = a + c

	}

	if err = scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}
	z.Records = records
	z.Regexp = regexp_records
	return nil
}

func (z *Zone) parseLine(line int, text string) (record *OrigRecord, err error) {
	var dtype uint16
	log.Printf("[FINE] line[%d] = [%s]", line, text)
	textlist := strings.SplitN(text, " ", 4)
	if len(textlist) < 4 {
		return nil, errors.New("line formart error")
	}
	label, err := z.parseLabel(textlist[0])
	if err != nil {
		return nil, fmt.Errorf("parselable error: %q", err)
	}

	punycode_lable, err := punycode.ToASCII(label)
	if err != nil {
		return nil, fmt.Errorf("parse label to punycode error: %q", err)
	}

	ttl, err := z.parseTtl(textlist[1])
	if err != nil {
		return nil, fmt.Errorf("parsettl error: %q", err)
	}

	rtype := strings.ToUpper(textlist[2])
	if dtype, err = plugins.DnsType(rtype); err != nil {
		return nil, fmt.Errorf("parsetype error: %q", err)
	}

	value, err := z.parseValue(textlist[3])
	if err != nil {
		return nil, fmt.Errorf("%s : %s", textlist[3], err)
	}
	record = &OrigRecord{
		Label: punycode_lable,
		Ttl:   ttl,
		Type:  dtype,
		Value: value,
	}
	if label == "@" && dtype == dns.TypeSOA {
		z.setSoaRR(ttl, value)
	}

	if label == "@" && dtype == dns.TypeNS {
		z.setNsRR(ttl, value)
	}

	return record, nil
}

func (z *Zone) setSoaRR(ttl int, conf map[string]interface{}) {

	rr_header := dns.RR_Header{
		Name:   z.Name + ".",
		Rrtype: dns.TypeSOA,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}
	z.Soa = &dns.SOA{
		Hdr:     rr_header,
		Ns:      conf["mname"].(string),
		Mbox:    conf["nname"].(string),
		Serial:  uint32(conf["serial"].(uint64)),
		Refresh: uint32(conf["refresh"].(uint64)),
		Retry:   uint32(conf["retry"].(uint64)),
		Expire:  uint32(conf["expire"].(uint64)),
		Minttl:  uint32(conf["minttl"].(uint64)),
	}
	log.Printf("[FINE] zone : %s, SOA=%s", z.Name, z.Soa)
}

func (z *Zone) setNsRR(ttl int, value map[string]interface{}) {
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
	z.Ns, _ = plugin.Filter(value)
}

func (z *Zone) parseLabel(label string) (ret string, err error) {
	return label, nil
}

func (z *Zone) parseTtl(ttl string) (ret int, err error) {
	newttl, err := strconv.Atoi(ttl)
	if err != nil {
		return 0, err
	}
	return newttl, nil
}

func (z *Zone) parseValue(value string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	_codec := codec.JsonHandle{}
	_codec.MapType = reflect.TypeOf(map[string]interface{}(nil))
	dec := codec.NewDecoderBytes([]byte(value), &_codec)
	err := dec.Decode(&ret)

	return ret, err
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
	var tlab string
	var ok bool
	var zk ZoneKey
	record := new(Record)
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
	zk = ZoneKey{Label: slab, Type: q.Qtype}
	if record, ok = z.Records[zk]; !ok {
		for z, r := range z.Regexp {
			regexp_record := strings.Replace(z.Label, "*", "\\w+", -1)
			if z.Type == q.Qtype && regexpcache.MustCompile(regexp_record).MatchString(slab) {
				if len(z.Label) > len(tlab) {
					log.Printf("[DEBUG] z.L = %s, tlab= %s", z.Label, tlab)
					tlab = z.Label
					record = r
				}
				ok = true
				log.Printf("[DEBUG] hit regexp : [%s] %s", slab, regexp_record)
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
		m.Answer, err = plugin.Filter(record.Value)
		/*
				switch rrtype {
				case dns.TypeA:
					ip := net.ParseIP(strings.TrimSpace("1.1.1.1"))
					answer = &dns.A{rr_header, ip}
				}
				m.Answer = append(m.Answer, answer)
			}
		*/
	} else {
		log.Printf("[WARING] record not found[%s] : %s", q.Name, dns.TypeToString[q.Qtype])
		err = fmt.Errorf("record not foud :%s=>%s", dns.TypeToString[q.Qtype], q.Name)
	}
	return
}
