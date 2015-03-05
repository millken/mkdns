package main

import (
	"expvar"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var VERSION string = "1.0.0"

var timeStarted = time.Now()
var qCounter = expvar.NewInt("qCounter")
var logger = NewLogger(os.Stderr, "", FINEST)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	var err error
	var configPath string
	flag.StringVar(&configPath, "c", "config.toml", "config path")
	flag.Parse()

	logger.Info("Loading config : %s, version: %s", configPath, VERSION)
	config, err := LoadConfig(configPath)
	if err != nil {
		logger.Exit("Read config failed.Err = %s", err.Error())
	}

	logger.Debug("config= %v", config)
	server := NewServer(config)
	server.Run()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	/*
		go func(c chan os.Signal) {
			// Wait for a signal:
			sig := <-c
			logger.Info("Caught signal '%s': shutting down.", sig)
			// Stop listening:

			// Delete the unix socket, if applicable:

			// And we're done:
			os.Exit(0)
		}(sigc)
		listen()
	*/
	<-sigc
	//log.Printf("Bye bye :( %s", sigc)
	logger.Info("god")

	//os.Exit(0)

}
