package main

import (
	"encoding/hex"
	"testing"

	"github.com/miekg/dns"
)

func TestParse(t *testing.T) {
	payload, _ := hex.DecodeString("4ffd8500000100020000000105617874717303636f6d0000010001c00c0001000100000258000401010101c00c000100010000025800040303030300002904d0000000000000")
	msg := new(dns.Msg)
	if err := msg.Unpack(payload); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", msg)
}
