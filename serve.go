package main

import (
	"github.com/miekg/dns"
	"gopkg.in/millken/logger.v1"
)


func listenAndServe(ip string) {

	prots := []string{"udp", "tcp"}

	for _, prot := range prots {
		go func(p string) {
			server := &dns.Server{Addr: ip, Net: p}

			log.Printf("Opening on %s %s", ip, p)
			if err := server.ListenAndServe(); err != nil {
				logger.Error("geodns: failed to setup %s %s: %s", ip, p, err)
			}
			logger.Fatal("mkdns: ListenAndServe unexpectedly returned")
		}(prot)
	}

}

