package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"main/components/options"
	"os"
)

const SessionFilePath = "./session.data"

// configLogger logger for config utils
var configLogger = logrus.WithField("context", "config")

// LoadClientConfig loads client config from yaml file
func LoadClientConfig(configPath string) (options.Options, error) {
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

	var cfg options.Options
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return options.Options{}, err
	}
	return cfg, nil
}

func LoadSession() (string, error) {
	session, err := os.Open(SessionFilePath)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 1024)
	_, err = session.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func StoreSession(accessKey string) error {
	session, err := os.OpenFile(SessionFilePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	_, err = session.WriteString(accessKey)
	if err != nil {
		return err
	}

	return nil
}
