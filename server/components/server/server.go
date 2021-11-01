package server

import (
	"github.com/sirupsen/logrus"
)

// Server declares generalized interface for server functionality
type Server interface {
	Run() error
	Close() error
}

// serverLogger logrus logger with context = server
var serverLogger = logrus.WithField("context", "server")
