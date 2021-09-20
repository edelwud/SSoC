package client

import (
	"github.com/sirupsen/logrus"
	"main/components/command"
	"main/components/session"
)

type Client interface {
	Connect() error
	Disconnect() error
	Exec(cmd command.Command) error
	Write(cmd string) error
	GetContext() session.Session
}

type Options struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	KeepAlive       bool   `yaml:"keepAlive"`
	KeepAlivePeriod int    `yaml:"keepAlivePeriod"`
}

var clientLogger = logrus.WithField("context", "client")
