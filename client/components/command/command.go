package command

import "main/components/session"

// Command declares interface for client commands
type Command interface {
	Row() []byte
	Process(ctx session.Session) error
}
