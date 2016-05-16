package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/gopacket/examples/util"
	"github.com/hashicorp/logutils"
)

var VERSION string = "2.0.0"

func main() {
	var err error
	var (
		configPath = flag.String("c", "config.toml", "config path")
	)
	defer util.Run()()

	if os.Geteuid() != 0 {
		log.Printf("requires root!")
		return
	}
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
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
	log.Printf("[INFO] Loading config : %s, version: %s", *configPath, VERSION)

	log.Printf("[DEBUG] config= %v , level=%s", config, config.Log.Level)

}
