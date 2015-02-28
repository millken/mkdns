package main

import (
	"expvar"
	"flag"
	//"gopkg.in/millken/logger.v1"
	//"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var VERSION string = "1.0.0"
var gitVersion string
var serverId string
var serverIP string
var serverGroups []string

var timeStarted = time.Now()
var qCounter = expvar.NewInt("qCounter")
var logger = NewLogger(os.Stderr, "", FINEST)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	var (
		//flagconfig = flag.String("config", "./", "directory of zone files")
		flaginter = flag.String("interface", "*", "set the listener address")
		flagport  = flag.String("port", "53", "default port number")
		//flaghttp   = flag.String("http", ":8053", "http listen address (:8053)")
		//flaglog    = flag.Bool("log", false, "be more verbose")
	)
	flag.Parse()

	logger.Info("Starting mkdns %s", VERSION)

	if *flaginter == "*" {
		addrs, _ := net.InterfaceAddrs()
		ips := make([]string, 0)
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if !(ip.IsLoopback() || ip.IsGlobalUnicast()) {
				continue
			}
			ips = append(ips, ip.String())
		}
		*flaginter = strings.Join(ips, ",")
	}

	inter := getInterfaces(*flaginter, *flagport)

	for _, host := range inter {
		go listenAndServe(host)
	}

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
