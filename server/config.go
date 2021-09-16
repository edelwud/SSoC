package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"server/components/server"
)

var configLogger = logrus.WithField("context", "config")

func LoadServerConfig(configPath string) (server.Options, error) {
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

	var cfg server.Options
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return server.Options{}, err
	}
	return cfg, nil
}
