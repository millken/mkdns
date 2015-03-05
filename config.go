package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type Config struct {
	Options OptionsConf
	Logging LoggingConf
}

type OptionsConf struct {
	Version      string
	Listen       []string
	ReadTimeout  int `toml:"read_timeout"`
	WriteTimeout int `toml:"write_timeout"`
}

type LoggingConf struct {
	File  string
	Debug string
}

func LoadConfig(configPath string) (config *Config, err error) {

	config = &Config{}
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
