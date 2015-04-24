package plugins

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

type Plugin interface {
	Filter(rr_header dns.RR_Header, conf map[string]interface{}) ([]dns.RR, error)
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

func Filter(recordType uint16, rr_header dns.RR_Header, config map[string]interface{}) (resp []dns.RR, err error) {
	if plugin, ok := plugins_list[recordType]; ok {

		plug := plugin()

		return plug.(Plugin).Filter(rr_header, config)
	}
	return nil, fmt.Errorf("plugin: %d not register", recordType)
}
