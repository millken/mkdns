package main

import "github.com/miekg/dns"

func main() {
	domain := dns.Fqdn("baidu.com")
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 2)
	msg.Question[0] = dns.Question{
		Name:   domain,
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	msg.Question[1] = dns.Question{
		Name:   domain,
		Qtype:  dns.TypeAAAA,
		Qclass: dns.ClassINET,
	}
	r, err := dns.Exchange(msg, "114.114.114.114:53")
	if err != nil {
		panic(err)
	}
	if r.Rcode != dns.RcodeSuccess {
		panic(r.Rcode)
	}
	for _, ans := range r.Answer {
		println(ans.String())
	}
}
