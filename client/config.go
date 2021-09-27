package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"main/components/client"
	"os"
)

// configLogger logger for config utils
var configLogger = logrus.WithField("context", "config")

// LoadClientConfig loads client config from yaml file
func LoadClientConfig(configPath string) (client.Options, error) {
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			configLogger.Fatalf("cannot close config file: %s", err)
		}
	}(f)

	var cfg client.Options
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return client.Options{}, err
	}
	return cfg, nil
}
