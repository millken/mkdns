package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
	"errors"
)

type TypeDef int

const (
	D_A TypeDef = iota
	D_AAAA
	D_CNAME
	D_NS
	D_TXT
	P_ALIAS
	P_VIEW
	P_WEIGHT
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

type Record struct {
	Label  string
	Ttl int
	Type TypeDef
	Value string
}

type Records []Record

type Zone struct {
	Records map[string]Records
	Soa *Soa
}

func NewZone() *Zone {
	zone := new(Zone)
	zone.Name = name
	zone.Soa = new(Soa)
	return zone
}

func (z *Zone) LoadFile(file string) (err error) {
	//var a, b []byte
	line := 1
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
			z.parseLine(line, a)
			logger.Debug("a=%s", a)
			a = ""
		}

		if c == "\"" && ahead != "\\" {
			b = math.Abs(b - 1)
		}

		if unicode.IsSpace([]rune(c)[0]) {
			//hack "
			if int(b) == 0 {
				c = "\t"
			}
			c = " "
		}
		ahead = c
		a = a + c

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}
	return nil
}

func (z *Zone) parseLine(line int, text string) (err error) {
	logger.Trace("line [%d] %s", line, text)
	text = strings.TrimSpace(text)
	if len(text) == 0 || text[0:1] == ";" {
		return nil
	}
	textlist := strings.Split(text, " ")
	if len(textlist) < 4 {
		return errors.New("line formart error")
	}
	if err = z.parseName(textlist[0]); err != nil {
		return fmt.Errorf("parsename error: %q", err)
	} 
	logger.Debug("%v", textlist)
	return nil
}

func (z *Zone) parseName(name string) error {
	
}