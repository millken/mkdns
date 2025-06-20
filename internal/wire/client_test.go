package wire

import (
	"net/netip"
	"testing"
	"time"

	"github.com/millken/mkdns/dns"
)

func testIPv6Network() error {
	client := &Client{
		AddrPort:    netip.AddrPortFrom(netip.MustParseAddr("2606:4700:4700::1111"), 53),
		ReadTimeout: 1 * time.Second,
		MaxConns:    5,
	}

	req, resp := new(dns.Request), new(dns.Response)
	req.SetQuestion(".", dns.TypeNS, dns.ClassINET)
	return client.Exchange(req, resp)
}

func TestClientIPv6(t *testing.T) {
	if err := testIPv6Network(); err != nil {
		t.Errorf("testIPv6Network error: %v", err)
	} else {
		t.Logf("testIPv6Network success")
	}
}

func TestClientExchange(t *testing.T) {
	var cases = []struct {
		Domain string
		Class  dns.Class
		Type   dns.Type
	}{
		{"www.google.com", dns.ClassINET, dns.TypeA},
	}

	client := &Client{
		AddrPort:    netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 53),
		ReadTimeout: 1 * time.Second,
		MaxConns:    1000,
	}

	for _, c := range cases {
		req, resp := dns.AcquireRequest(), dns.AcquireResponse()
		req.SetQuestion(c.Domain, c.Type, c.Class)
		err := client.Exchange(req, resp)
		if err != nil {
			t.Errorf("client=%+v exchange(%v) error: %+v\n", client, c.Domain, err)
		}
		t.Logf("%s: CLASS %s TYPE %s\n", resp.Domain, resp.Question.Class, resp.Question.Type)
		dns.ReleaseRequest(req)
		defer dns.ReleaseResponse(resp)
		if resp.Header.Rcode() != dns.RcodeSuccess {
			t.Errorf("client=%+v exchange(%v) error: response rcode %s != %s\n", client, c.Domain, resp.Header.Rcode(), dns.RcodeSuccess)
			continue
		}
		if resp.Header.Ancount == 0 {
			t.Errorf("client=%+v exchange(%v) error: response ancount is 0\n", client, c.Domain)
			continue
		}
		for _, rr := range resp.Answer {
			if rr.Header().Rrtype != c.Type {
				t.Errorf("client=%+v exchange(%v) error: response rr type %s != %s\n", client, c.Domain, rr.Header().Rrtype, c.Type)
			}
			if rr.Header().Class != c.Class {
				t.Errorf("client=%+v exchange(%v) error: response rr class %s != %s\n", client, c.Domain, rr.Header().Class, c.Class)
			}
			t.Logf("response rr: %+v\n", rr)
		}
	}
}
