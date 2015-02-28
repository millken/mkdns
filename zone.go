package main

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"os"
	"strings"
	"unicode"
	"errors"
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
	Label  string
	Ttl int
	Type TypeDef
	Value string
}
type Record struct {
	Ttl int
	Value string
}

type Zck struct {
	L string
	T TypeDef
}

type Records map[Zck][]*Record

type Zone struct {
	Records Records
	Soa *Soa
}

func NewZone() *Zone {
	zone := new(Zone)
	zone.Soa = new(Soa)
	zone.Records = make(map[Zck][]*Record)
	return zone
}

func (z *Zone) LoadFile(file string) (err error) {
	//var a, b []byte
	records	:= make(Records)
	line := 0
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	inputFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Error opening zone file: %s", err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanBytes)
	a := ""
	ahead := ""
	b := 1.0

	for scanner.Scan() {
		//logger.Debug("%v %s", scanner.Bytes(), scanner.Text())		
		c := scanner.Text()
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
			records[zck] = append(records[zck], &Record{
				Ttl: r.Ttl,
				Value: r.Value,
			})

		}

		if c == "\"" && ahead != "\\" {
			b = math.Abs(b - 1)
		}

		if unicode.IsSpace([]rune(c)[0]) {
			if c == " " && ahead == " " {
				c = ""
			}else
			//hack "
			if int(b) == 0 {
				c = "\t"
			}else{
				c = " "
			}
		}
		ahead = c
		a = a + c

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}
	z.Records = records
	logger.Debug("zone = %v", z)
	return nil
}

func (z *Zone) parseLine(line int, text string) (record *ORecord, err error) {


	textlist := strings.Split(text, " ")
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
		Ttl: ttl,
		Type: typedef,
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
	switch typedef {
	case "A": ret = D_A
	case "AAAA": ret = D_AAAA
	case "CNAME": ret = D_CNAME
	case "NS": ret = D_NS
	case "TXT": ret = D_TXT
	case "SOA": ret = D_SOA
	default: ret = P_ERROR
	}
	if ret == P_ERROR {
		return ret, fmt.Errorf("unknown type: %s", typedef)
	}
	return ret, nil
}

func (z *Zone) parseValue(value string) (ret string, err error) {
	return value, nil
}