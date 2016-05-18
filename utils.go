package main

import "net"

func grabAddresses(iface string) (macAddr net.HardwareAddr, ipAddr net.IP) {

	netInterface, err := net.InterfaceByName(iface)
	if err != nil {
		panic(err)
	}

	macAddr = netInterface.HardwareAddr
	addrs, _ := netInterface.Addrs()
	ipAddr, _, err = net.ParseCIDR(addrs[0].String())
	if err != nil {
		panic(err)
	}

	return macAddr, ipAddr
}
