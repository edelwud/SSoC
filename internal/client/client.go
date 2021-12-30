package client

import (
	"SSoC/internal/options"
	"SSoC/internal/requester"
	"SSoC/internal/session"
	"github.com/sirupsen/logrus"
)

// Client declares generalized interface for client functionality
type Client interface {
	Connect() error
	Disconnect() error
	Exec(cmd requester.Requester) error
	Write(cmd string) error
	GetContext() session.Session
}

var clientLogger = logrus.WithField("context", "client")
var clientTypes = map[string]func(options.Options, string) Client{
	"udp": func(o options.Options, token string) Client {
		return CreateUDPClient(o, token)
	},
	"tcp": func(o options.Options, token string) Client {
		return CreateTCPClient(o, token)
	},
}

func New(protocol string, options options.Options, token string) Client {
	if client, ok := clientTypes[protocol]; !ok {
		return nil
	} else {
		return client(options, token)
	}
}
