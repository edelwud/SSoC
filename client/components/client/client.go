package client

import (
	"github.com/sirupsen/logrus"
	"main/components/command"
)

type Client interface {
	Connect() error
	Disconnect() error
	Exec(cmd command.Command) error
}

type Options struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	KeepAlive       bool   `yaml:"keepAlive"`
	KeepAlivePeriod int    `yaml:"keepAlivePeriod"`
}

var clientLogger = logrus.WithField("context", "client")
