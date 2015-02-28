package main

import (
	//"github.com/miekg/dns"
	"time"
)

type Server struct {
	addr     string
	rTimeout time.Duration
	wTimeout time.Duration
}

func (s *Server) Run() {
}
