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
	Stats  StatsConf
}

type ServerConf struct {
	Iface, Driver, Backend string
	WorkerNum              int    `toml:"worker_num"`
	IPDBPath               string `toml:"ip_db_path"`
	//ReadTimeout  int `toml:"read_timeout"`
	//WriteTimeout int `toml:"write_timeout"`
}

type LogConf struct {
	File  string
	Level string
}

type StatsConf struct {
	Addr, Url, Schedule string
	AutoReport          bool `toml:"auto_report"`
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
