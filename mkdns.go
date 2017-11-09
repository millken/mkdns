package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket/examples/util"
	"github.com/hashicorp/logutils"
	"github.com/millken/mkdns/backends"
	"github.com/millken/mkdns/ip"
	"github.com/millken/mkdns/stats"
)

var VERSION string = "2.0.0"

func main() {
	var err error
	var (
		configPath = flag.String("c", "config.toml", "config path")
	)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic ->>>> %s", err)
		}
	}()
	defer util.Run()()

	if os.Geteuid() != 0 {
		log.Printf("requires root!")
		return
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Printf("[ERROR] LoadConfig : %s", err.Error())
		return
	}
	filterWriter := os.Stderr
	if config.Log.File != "" {
		filterWriter, err = os.Create(config.Log.File)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"FINE", "DEBUG", "TRACE", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(config.Log.Level),
		Writer:   filterWriter,
	}
	log.SetOutput(filter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("[INFO] Loading config : %s, version: %s", *configPath, VERSION)

	//load cnip db
	if config.Server.IPDBPath != "" {
		err = ip.LoadCNIpDB(config.Server.IPDBPath)
		if err != nil {
			log.Printf("[ERROR] LoadCNIpDB : %s", err.Error())
			return
		}
	}
	//load backend
	backend, err := backends.Open(config.Server.Backend)
	if err != nil {
		log.Fatalf("backend open error : %s", err)
	}
	go backend.Load()
	if config.Stats.Addr != "" {
		go func() {
			statsServer := stats.NewServer(config.Stats.Addr)
			if err := statsServer.Run(); err != nil {
				log.Fatalf("stats server run err: %s", err)
			}
		}()
	}

	if config.Stats.AutoReport {
		autoreport := stats.NewAutoReport(config.Stats.Url, config.Stats.Schedule)
		autoreport.Start()
	}

	server := NewServer(config)
	if err = server.Start(); err != nil {
		log.Printf("[ERROR] :%s", err)
	}
}
