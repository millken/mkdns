package types

import (
	"net"

	"github.com/miekg/dns"
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

type Z struct {
	l string
	t uint16
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
