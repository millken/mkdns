package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConf
	Log    LogConf
}

type ServerConf struct {
	Iface, Driver, Backend string
	//ReadTimeout  int `toml:"read_timeout"`
	//WriteTimeout int `toml:"write_timeout"`
}

type LogConf struct {
	File  string
	Level string
}

func LoadConfig(configPath string) (config *Config, err error) {

	p, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %s", err)
	}
	contents, err := ioutil.ReadAll(p)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %s", err)
	}
	/*
		config.Options = OptionsConf{
			ReadTimeout:  3,
			WriteTimeout: 3,
			Listen:       []string{"127.0.0.1:53"},
			Zones:        "mkdns.zones",
		}
	*/
	if _, err = toml.Decode(string(contents), &config); err != nil {
		return nil, fmt.Errorf("Error decoding config file: %s", err)
	}

	return config, nil
}
