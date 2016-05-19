package protocols

import (

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/miekg/dns"
)

func DNSParser(layer gopacket.Layer) (*dns.Msg, error) {

	dnsLayer := layer.(*layers.DNS)

	contents := dnsLayer.BaseLayer.LayerContents()

	dnsMsg := new(dns.Msg)
	if err := dnsMsg.Unpack(contents); err != nil {
		return dnsMsg, err
	}

	return dnsMsg, nil
}
