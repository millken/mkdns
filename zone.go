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

	"github.com/miekg/dns"
	"github.com/millken/mkdns/plugins"
	"github.com/ugorji/go/codec"
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
	Soa     *Soa
	Options ZoneOptions
}

func NewZone() *Zone {
	zone := new(Zone)
	zone.Soa = new(Soa)
	zone.Records = make(map[Zck]*Record)
	zone.Options.EdnsAddr = nil
	zone.Options.RemoteAddr = nil
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
			zck := Zck{
				L: r.Label,
				T: r.Type,
			}
			records[zck] = &Record{
				Ttl:   r.Ttl,
				Value: r.Value,
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
	logger.Debug("zone(%s)\n %s", file, z)
	return nil
}

func (z *Zone) String() string {
	ret := fmt.Sprintf("SOA { %s %s\n %d %d %d %d %d\n }", z.Soa.Mname, z.Soa.Nname, z.Soa.Serial, z.Soa.Refresh, z.Soa.Retry, z.Soa.Expire, z.Soa.Minttl)
	for zck, r := range z.Records {
		ret = fmt.Sprintf("%s\n (label)=%s, (type)=%s, (ttl)=%d, (value)=%s", ret, zck.L, zck.T, r.Ttl, r.Value)
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
	var dtype uint16
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

	rtype := strings.ToUpper(textlist[2])
	if dtype, err = plugins.DnsType(rtype); err != nil {
		return nil, fmt.Errorf("parsetype error: %q", err)
	}

	value, err := z.parseValue(textlist[3])
	if err != nil {
		return nil, fmt.Errorf("%s : %s", textlist[3], err)
	}
	record = &ORecord{
		Label: label,
		Ttl:   ttl,
		Type:  dtype,
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

func (z *Zone) parseValue(value string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	_codec := codec.JsonHandle{}
	_codec.MapType = reflect.TypeOf(map[string]interface{}(nil))
	dec := codec.NewDecoderBytes([]byte(value), &_codec)
	err := dec.Decode(&ret)

	return ret, err
}

func (z *Zone) FindRecord(req *dns.Msg) (m *dns.Msg, err error) {
	//var answer dns.RR
	var slab string
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
	zck := Zck{L: slab, T: q.Qtype}
	if r, ok := z.Records[zck]; ok {

		rr_header := dns.RR_Header{
			Name:   q.Name,
			Rrtype: q.Qtype,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.Ttl),
		}
		plugin := plugins.Get(q.Qtype).(plugins.Plugin)
		if plugin == nil {
			err = fmt.Errorf("plugin: %d not register", q.Qtype)
			return
		}
		plugin.New(z.Options.EdnsAddr, z.Options.RemoteAddr, rr_header)
		m.Answer, err = plugin.Filter(r.Value)
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
		logger.Debug("record not found :%+v", zck)
		err = fmt.Errorf("record not foud :%s", q.Name)
	}
	return
}
