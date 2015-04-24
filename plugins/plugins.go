package plugins

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

type Plugin interface {
	New(edns, remote net.IP, rr_header dns.RR_Header)
	Filter(conf map[string]interface{}) ([]dns.RR, error)
}

var plugins_type = make(map[string]uint16)
var plugins_list = make(map[uint16]func() interface{})

func RegisterPlugin(name string, rType uint16, plugin func() interface{}) {
	if plugin == nil {
		log.Fatalln("plugin: Register plugin is nil")
	}

	if _, ok := plugins_list[rType]; ok {
		log.Fatalln("plugin: Register called twice for plugin " + name)
	}
	log.Println("RegisterPlugin: ", name)

	plugins_type[name] = rType
	plugins_list[rType] = plugin
}

func DnsType(recordType string) (dType uint16, err error) {
	dType, ok := plugins_type[recordType]
	if !ok {
		return 0, fmt.Errorf("type not allowed: %s", recordType)
	}
	return dType, nil
}

func Get(recordType uint16) interface{} {
	if plugin, ok := plugins_list[recordType]; ok {

		plug := plugin()

		return plug.(Plugin)
	}
	return nil
}

func Filter(recordType uint16, config map[string]interface{}) (resp []dns.RR, err error) {
	if plugin, ok := plugins_list[recordType]; ok {

		plug := plugin()

		return plug.(Plugin).Filter(config)
	}
	return nil, fmt.Errorf("plugin: %d not register", recordType)
}
