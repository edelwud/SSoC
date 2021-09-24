package server

import (
	"github.com/sirupsen/logrus"
)

// Server declares generalized interface for server functionality
type Server interface {
	Run() error
	Close() error
}

// Options declares generalized structure of server parameters
type Options struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	KeepAlive       bool   `yaml:"keepAlive"`
	KeepAlivePeriod int    `yaml:"keepAlivePeriod"`
}

// serverLogger logrus logger with context = server
var serverLogger = logrus.WithField("context", "server")
