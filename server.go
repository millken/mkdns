package main

import (
	"time"
	"github.com/miekg/dns"

)

type Server struct {
	addr string
	rTimeout time.Duration
	wTimeout time.Duration
}

func (s *Server) Run() {
}
