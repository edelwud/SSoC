package server

import (
	"SSoC/internal/options"
	"github.com/sirupsen/logrus"
)

// Server declares generalized interface for server functionality
type Server interface {
	Run() error
	Close() error
	Write(payload, token string) error
}

var serverLogger = logrus.WithField("context", "server")
var serverTypes = map[string]func(options.Options) Server{
	"udp": func(o options.Options) Server {
		return CreateUDPServer(o)
	},
	"tcp": func(o options.Options) Server {
		return CreateTCPServer(o)
	},
}

func New(protocol string, options options.Options) Server {
	if server, ok := serverTypes[protocol]; !ok {
		return nil
	} else {
		return server(options)
	}
}
