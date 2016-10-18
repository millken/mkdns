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
	WorkerNum              int    `toml:"worker_num"`
	StatsAddr              string `toml:"stats_addr"`
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
	if _, err = toml.Decode(string(contents), &config); err != nil {
		return nil, fmt.Errorf("Error decoding config file: %s", err)
	}

	return config, nil
}
