package main

import (
	"bufio"
	"os"
	"math"
	"strings"
	"unicode"
	"fmt"
)

type Soa struct {
	Mname string
	Nname    string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minttl  uint32
}

type Zone struct {
	Soa *Soa
}

func NewZone(name string) *Zone {
	zone := new(Zone)
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

		if c == "\"" {
			b = math.Abs(b - 1)
		}

		if unicode.IsSpace([]rune(c)[0]) {
			c = " "
		}
			a = a + c
		
	}
 
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading zone file: %s, line : %d", scanner.Err(), line)
	}	
	return nil
}

func (z *Zone) parseLine(line int, text string) {
	logger.Trace("line [%d] %s", line, text)
	text = strings.TrimSpace(text)
	if len(text) == 0 || text[0:1] == ";" {
		return 
	}
	textlist := strings.Split(text, " ")
	logger.Debug("%q", textlist)
	return
}