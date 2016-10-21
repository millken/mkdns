package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/VividCortex/godaemon"
	"github.com/google/gopacket/examples/util"
	"github.com/hashicorp/logutils"
	"github.com/millken/mkdns/backends"
	"github.com/millken/mkdns/stats"
)

var VERSION string = "2.0.0"

func main() {
	var err error
	var (
		configPath = flag.String("c", "config.toml", "config path")
		isDaemon   = flag.Bool("d", false, "backgroud running")
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
	if *isDaemon {
		defer godaemon.Daemonize()
	}
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Printf("[ERROR] LoadConfig : %s", err.Error())
		return
	}
	filter_writer := os.Stderr
	if config.Log.File != "" {
		filter_writer, err = os.Create(config.Log.File)
		if err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"FINE", "DEBUG", "TRACE", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(config.Log.Level),
		Writer:   filter_writer,
	}
	log.SetOutput(filter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("[INFO] Loading config : %s, version: %s", *configPath, VERSION)

	//load backend
	backend, err := backends.Open(config.Server.Backend)
	if err != nil {
		log.Fatalf("backend open error : %s", err)
	}
	go backend.Load()
	if config.Server.StatsAddr != "" {
		go func() {
			statsServer := stats.NewServer(config.Server.StatsAddr)
			if err := statsServer.Run(); err != nil {
				log.Fatalf("stats server run err: %s", err)
			}
		}()
	}

	server := NewServer(config)
	if err = server.Start(); err != nil {
		log.Printf("[ERROR] :%s", err)
	}
}
