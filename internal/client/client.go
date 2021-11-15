package client

import (
	"SSoC/internal/requester"
	"SSoC/internal/session"
)

// Client declares generalized interface for client functionality
type Client interface {
	Connect() error
	Disconnect() error
	Exec(cmd requester.Requester) error
	Write(cmd string) error
	GetContext() session.Session
}
