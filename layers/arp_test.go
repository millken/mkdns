package layers

import (
	"fmt"
	"testing"
)

func Test_GetAll(t *testing.T) {
	p := []byte{
		0x00, 0x01, 0x08, 0x00, 0x06, 0x04, 0x00, 0x01,
	}
	arp := *(*ARP)(&p)
	fmt.Println(Swap16(arp.GetLinkType()))
	fmt.Println(Swap16(arp.GetProtocolType()))
	fmt.Println(arp.GetLinkAddressLength())
	fmt.Println(arp.GetProtocolAddressLength())
	fmt.Println(Swap16(arp.GetOpCode()))
}

func Test_SetAll(t *testing.T) {
	p := make([]byte, 8)

	arp := *(*ARP)(&p)
	arp.SetLinkType(Swap16(LinkTypeEthernet))
	arp.SetProtocolType(Swap16(EthernetTypeIPv4))
	arp.SetLinkAddressLength(6)
	arp.SetProtocolAddressLength(4)
	arp.SetOpCode(Swap16(ARPRequest))

	fmt.Printf("%x\n", arp)
}
