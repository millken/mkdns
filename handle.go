package main

import (
	"github.com/miekg/dns"
	"github.com/rcrowley/go-metrics"
)

type Handler struct {
	responseTimer metrics.Timer
}

func (h *Handler) UDP(w dns.ResponseWriter, req *dns.Msg) {
}
