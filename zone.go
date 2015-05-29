package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/Zverushko/punycode"
	"github.com/miekg/dns"
	"github.com/millken/logger"
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

type ORecord struct {
	Label string
	Ttl   int
	Type  uint16
	Value map[string]interface{}
}
type Record struct {
	Ttl   int
	Value map[string]interface{}
}

type Zck struct {
	L string
	T uint16
}

type Records map[Zck]*Record

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
	zone.Records = make(map[Zck]*Record)
	zone.Regexp = make(map[Zck]*Record)
	zone.Options.EdnsAddr = nil
	zone.Options.RemoteAddr = nil
	return zone
}

func (z *Zone) LoadFile(file string) (err error) {
	//var a, b []byte
	records := make(Records)
	regexp_records := make(Records)
	line := 0
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	inputFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening zone file: %s", err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanBytes)
	a := ""
	ahead := ""
	b := 1.0
	aheadspace := false

	for scanner.Scan() {
		//logger.Debug("%v %s", scanner.Bytes(), scanner.Text())
		c := scanner.Text()
		//logger.Debug("c=%s, unicode.IsSpace([]rune(c)[0])=%v", c, unicode.IsSpace([]rune(c)[0]))
		if c == "\n" && int(b) == 1 {
			line = line + 1
			text := strings.TrimSpace(a)

			logger.Finest("(%s:%d)=%s", file, line, a)
			a = ""
			if len(text) == 0 || text[0:1] == ";" {
				continue
			}
			r, rerr := z.parseLine(line, text)
			if rerr != nil {
				logger.Warn("parse zone %s:%d error: %q", file, line, rerr)
				continue
			}
			zck := Zck{
				L: r.Label,
				T: r.Type,
			}
			if strings.Contains(r.Label, "*") {
				regexp_records[zck] = &Record{
					Ttl:   r.Ttl,
					Value: r.Value,
				}
			} else {
				records[zck] = &Record{
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

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}
	z.Records = records
	z.Regexp = regexp_records
	logger.Fine("zone(%s)\n %v", file, z)
	return nil
}

/*
func (z *Zone) String() string {
	ret := fmt.Sprintf("SOA { %s %s\n %d %d %d %d %d\n }", z.Soa.Mname, z.Soa.Nname, z.Soa.Serial, z.Soa.Refresh, z.Soa.Retry, z.Soa.Expire, z.Soa.Minttl)
	for zck, r := range z.Records {
		ret = fmt.Sprintf("%s\n (label)=%s, (type)=%s, (ttl)=%d, (value)=%s", ret, zck.L, zck.T, r.Ttl, r.Value)
	}
	return ret
}
*/
func (z *Zone) parseLine(line int, text string) (record *ORecord, err error) {
	var dtype uint16
	logger.Fine("line[%d] = [%s]", line, text)
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
	record = &ORecord{
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
	logger.Debug("zone : %s, SOA=%s", z.Name, z.Soa)
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
		logger.Error("plugin: %d not register", dns.TypeNS)
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
	var zck Zck
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

	logger.Debug("z.Name=%s, q.Name=%s, slab=%s, q.Qtype=%d, z.Options=%+v", z.Name, q.Name, slab, q.Qtype, z.Options)
	zck = Zck{L: slab, T: q.Qtype}
	if record, ok = z.Records[zck]; !ok {
		for z, r := range z.Regexp {
			regexp_record := strings.Replace(z.L, "*", "\\w+", -1)
			if z.T == q.Qtype && regexpcache.MustCompile(regexp_record).MatchString(slab) {
				if len(z.L) > len(tlab) {
					logger.Debug("z.L = %s, tlab= %s", z.L, tlab)
					tlab = z.L
					record = r
				}
				ok = true
				logger.Debug("hit regexp : [%s] %s", slab, regexp_record)
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
		logger.Warn("record not found[%s] : %s", q.Name, dns.TypeToString[q.Qtype])
		err = fmt.Errorf("record not foud :%s=>%s", dns.TypeToString[q.Qtype], q.Name)
	}
	return
}
