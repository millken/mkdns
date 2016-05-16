package main

import "sync"

const BPFFilter = "udp port 53"

type server struct {
	config *Config
	group  *sync.WaitGroup
	//handler  *Handler
	//rTimeout time.Duration
	//wTimeout time.Duration
}

func NewServer(config *Config) *server {
	return &server{
		config: config,
		group:  new(sync.WaitGroup),
	}
}

func (s *server) Run() error {
	return nil
}
