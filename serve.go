package main

import (
	"github.com/miekg/dns"
)

func listenAndServe(ip string) {
	prots := []string{"udp", "tcp"}

	for _, prot := range prots {
		go func(p string) {
			server := &dns.Server{Addr: ip, Net: p}

			logger.Fine("Opening on %s %s", ip, p)
			if err := server.ListenAndServe(); err != nil {
				logger.Error("failed to setup %s %s: %s", ip, p, err)
			}
		}(prot)
	}

}
