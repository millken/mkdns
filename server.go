package main

import (
	"github.com/miekg/dns"
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

type server struct {
	config     *Config
	group      *sync.WaitGroup
	udpHandler *dns.ServeMux
	//handler  *Handler
	//rTimeout time.Duration
	//wTimeout time.Duration
}

func NewServer(config *Config) *server {

	return &server{
		config: config,
		//handler: NewHandler(),
		group: new(sync.WaitGroup),
	}
}

func (s *server) Run() error {

	responseTimer := metrics.NewTimer()
	handler := &Handler{
		responseTimer: responseTimer,
	}
	s.udpHandler = dns.NewServeMux()
	s.udpHandler.HandleFunc(".", handler.UDP)
	for _, addr := range s.config.Options.Listen {
		udpServer := &dns.Server{Addr: addr,
			Net:          "udp",
			Handler:      s.udpHandler,
			UDPSize:      65535,
			ReadTimeout:  time.Duration(s.config.Options.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(s.config.Options.WriteTimeout) * time.Second,
		}
		go s.listenAndServe(addr, udpServer)
	}
	return nil
}

func (s *server) listenAndServe(addr string, ds *dns.Server) {
	logger.Info("Opening on %s", addr)
	if err := ds.ListenAndServe(); err != nil {
		logger.Exit("failed to setup %s: %s", addr, err)
	}

}
