package server

import "github.com/sirupsen/logrus"

type Server interface {
	Run() error
	Close() error
}

type Options struct {
	Host string
	Port string
}

var serverLogger = logrus.WithField("context", "server")
