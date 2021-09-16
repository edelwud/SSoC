package server

import (
	"github.com/sirupsen/logrus"
)

type Server interface {
	Run() error
	Close() error
}

type Options struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	KeepAlive       bool   `yaml:"keepAlive"`
	KeepAlivePeriod int    `yaml:"keepAlivePeriod"`
}

var serverLogger = logrus.WithField("context", "server")
