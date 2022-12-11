package dns

import "net"

// SOA RR. See RFC 1035.
type SOA struct {
	Mname                                  net.NS
	Nname                                  net.NS
	Serial, Refresh, Retry, Expire, MinTTL uint32
}
