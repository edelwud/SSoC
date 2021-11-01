package client

import (
	"github.com/sirupsen/logrus"
	"main/components/command"
	"main/components/session"
)

// Client declares generalized interface for client functionality
type Client interface {
	Connect() error
	Disconnect() error
	Exec(cmd command.Command) error
	Write(cmd string) error
	GetContext() session.Session
}

// clientLogger logrus logger with context = client
var clientLogger = logrus.WithField("context", "client")
