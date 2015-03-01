package main

import (
	"expvar"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/millken/logger"
)

var VERSION string = "1.0.0"

var timeStarted = time.Now()
var qCounter = expvar.NewInt("qCounter")

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	var err error
	var configPath, debugLevel string
	flag.StringVar(&configPath, "c", "config.toml", "config path")
	flag.StringVar(&debugLevel, "debug", "INFO", "FINE|DEBUG|TRACE|INFO|ERROR")
	flag.Parse()

	logLevel := logger.INFO
	switch debugLevel {
	case "FINE":
		logLevel = logger.FINE
	case "DEBUG":
		logLevel = logger.DEBUG
	case "TRACE":
		logLevel = logger.TRACE
	case "ERROR":
		logLevel = logger.ERROR
	}
	logger.Global = logger.NewDefaultLogger(logLevel)
	logger.Info("Loading config : %s, version: %s", configPath, VERSION)
	config, err := LoadConfig(configPath)
	if err != nil {
		logger.Critical("Read config failed.Err = %s", err.Error())
	}

	logger.Finest("config= %v", config)

	LoadZones(config.Options.Zones)

	server := NewServer(config)
	server.Run()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	/*
		go func(c chan os.Signal) {
			// Wait for a signal:
			sig := <-c
			loggerger.Info("Caught signal '%s': shutting down.", sig)
			// Stop listening:

			// Delete the unix socket, if applicable:

			// And we're done:
			os.Exit(0)
		}(sigc)
		listen()
	*/
	<-sigc
	//logger.Printf("Bye bye :( %s", sigc)
	logger.Info("god")

	//os.Exit(0)

}
