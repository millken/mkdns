package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"net"
	"strconv"
	"strings"
	"unicode"

	"github.com/miekg/dns"
)

type TypeDef int

const (
	D_SOA TypeDef = iota
	D_A
	D_AAAA
	D_CNAME
	D_NS
	D_TXT
	P_ALIAS
	P_VIEW
	P_WEIGHT
	P_ERROR
)

var (
	typeDefStrings = [...]string{"SOA", "A", "AAAA", "CNAME", "NS", "TXT", "!ALIAS", "!VIEW", "!WEIGHT", "!ERROR"}
)

func (t TypeDef) String() string {
	if t < 0 || int(t) > len(typeDefStrings) {
		return "UNKNOWN"
	}
	return typeDefStrings[int(t)]
}

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
	Type  TypeDef
	Value string
}
type Record struct {
	Ttl   int
	Value string
}

type Zck struct {
	L string
	T TypeDef
}

type Records map[Zck][]*Record

type Zone struct {
	Name string
	Records Records
	Soa     *Soa
}

func NewZone() *Zone {
	zone := new(Zone)
	zone.Soa = new(Soa)
	zone.Records = make(map[Zck][]*Record)
	return zone
}

func (z *Zone) LoadFile(file string) (err error) {
	//var a, b []byte
	records := make(Records)
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
		logger.Debug("c=%s, unicode.IsSpace([]rune(c)[0])=%v", c, unicode.IsSpace([]rune(c)[0]))
		if c == "\n" && int(b) == 1 {
			line = line + 1
			text := strings.TrimSpace(a)

			logger.Trace("(%s:%d)=%s", file, line, a)
			a = ""
			if len(text) == 0 || text[0:1] == ";" {
				continue
			}
			r, rerr := z.parseLine(line, text)
			if rerr != nil {
				logger.Warn("parse zone %s:%d error: %q", file, line, rerr)
				continue
			}
			if r.Type == D_SOA {
				z.Soa, err = z.parseSoa(strings.TrimSpace(r.Value))
				if err != nil {
					return fmt.Errorf("parse soa %s:%d error: %q", file, line, err)
				}
				continue
			}
			zck := Zck{
				L: r.Label,
				T: r.Type,
			}
			records[zck] = append(records[zck], &Record{
				Ttl:   r.Ttl,
				Value: r.Value,
			})

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
		}else{
			aheadspace = false
		}
		ahead = c
		a = a + c

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}
	z.Records = records
	logger.Debug("zone(%s)\n %s", file, z)
	return nil
}

func (z *Zone) String() string {
	ret := fmt.Sprintf("SOA { %s %s\n %d %d %d %d %d\n }", z.Soa.Mname, z.Soa.Nname, z.Soa.Serial, z.Soa.Refresh, z.Soa.Retry, z.Soa.Expire, z.Soa.Minttl)
	for zck, r := range z.Records {
		for _, re := range r {
			ret = fmt.Sprintf("%s\n (label)=%s, (type)=%s, (ttl)=%d, (value)=%s", ret, zck.L, zck.T, re.Ttl, re.Value)
		}
	}
	return ret
}

func (z *Zone) parseSoa(s string) (soa *Soa, err error) {

	slist := strings.SplitN(s, " ", 7)
	if len(slist) < 7 {
		return nil, errors.New("line formart error")
	}
	//soa = new(Soa)
	mname := slist[0]
	nname := slist[1]
	serial, err := strconv.Atoi(slist[2])
	if err != nil {
		return nil, err
	}
	refresh, err := strconv.Atoi(slist[3])
	if err != nil {
		return nil, err
	}
	retry, err := strconv.Atoi(slist[3])
	if err != nil {
		return nil, err
	}
	expire, err := strconv.Atoi(slist[3])
	if err != nil {
		return nil, err
	}
	minttl, err := strconv.Atoi(slist[3])
	if err != nil {
		return nil, err
	}
	soa = &Soa{
		Mname:   mname,
		Nname:   nname,
		Serial:  uint32(serial),
		Refresh: uint32(refresh),
		Retry:   uint32(retry),
		Expire:  uint32(expire),
		Minttl:  uint32(minttl),
	}

	return
}

func (z *Zone) parseLine(line int, text string) (record *ORecord, err error) {

	logger.Debug("line[%d] = [%s]", line, text)
	textlist := strings.SplitN(text, " ", 4)
	if len(textlist) < 4 {
		return nil, errors.New("line formart error")
	}
	label, err := z.parseLabel(textlist[0])
	if err != nil {
		return nil, fmt.Errorf("parselable error: %q", err)
	}

	ttl, err := z.parseTtl(textlist[1])
	if err != nil {
		return nil, fmt.Errorf("parsettl error: %q", err)
	}

	typedef, err := z.parseType(textlist[2])
	if err != nil {
		return nil, fmt.Errorf("parsetype error: %q", err)
	}

	value, err := z.parseValue(textlist[3])
	if err != nil {
		return nil, fmt.Errorf("parsevalue error: %q", err)
	}
	record = &ORecord{
		Label: label,
		Ttl:   ttl,
		Type:  typedef,
		Value: value,
	}

	logger.Debug("%v", textlist)
	return record, nil
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

func (z *Zone) parseType(typedef string) (ret TypeDef, err error) {
	switch strings.ToUpper(typedef) {
	case "A":
		ret = D_A
	case "AAAA":
		ret = D_AAAA
	case "CNAME":
		ret = D_CNAME
	case "NS":
		ret = D_NS
	case "TXT":
		ret = D_TXT
	case "SOA":
		ret = D_SOA
	default:
		ret = P_ERROR
	}
	if ret == P_ERROR {
		return ret, fmt.Errorf("unknown type: %s", typedef)
	}
	return ret, nil
}

func (z *Zone) parseValue(value string) (ret string, err error) {
	return value, nil
}

func (z *Zone) FindRecord(req *dns.Msg) (m *dns.Msg, err error) {
	var typedef TypeDef
	var answer dns.RR
	var slab string
	q := req.Question[0]
	m = new(dns.Msg)
	m.SetReply(req)	
	typedef = P_ERROR
	rrtype := q.Qtype
	switch q.Qtype {
		case dns.TypeA:
			typedef = D_A
		case dns.TypeAAAA:
			typedef = D_AAAA
		case dns.TypeCNAME:
			typedef = D_CNAME
		case dns.TypeNS:
			typedef = D_NS
	}
	if len(q.Name) == len(z.Name) + 1 {
		slab = "@" 
	}else {
		tl := len(q.Name) - len(z.Name) - 2
		slab = strings.ToLower(q.Name[0:tl])
	}

	logger.Debug("z.Name=%s, q.Name=%s, slab=%s", z.Name, q.Name, slab)
	zck := Zck{L: slab, T: typedef}
	if r, ok := z.Records[zck]; ok {
		for _, re := range r {
			rr_header := dns.RR_Header{
				Name: q.Name,
				Rrtype: rrtype,
				Class: dns.ClassINET,
				Ttl: uint32(re.Ttl),
			}
			switch rrtype {
			case dns.TypeA:
				ip := net.ParseIP(strings.TrimSpace(re.Value))
				answer = &dns.A{rr_header, ip}
			case dns.TypeAAAA:
				ip := net.ParseIP(strings.TrimSpace(re.Value))
				answer = &dns.AAAA{rr_header, ip}				
			case dns.TypeCNAME:
				answer = &dns.CNAME{rr_header, strings.TrimSpace(re.Value)}
			case dns.TypeNS:
				answer = &dns.NS{rr_header, strings.TrimSpace(re.Value)}
			}
			
			m.Answer = append(m.Answer, answer)			
		}
	}else{
		logger.Debug("not found =%q", zck)
	}
	return 
}
