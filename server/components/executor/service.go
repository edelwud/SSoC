package executor

import (
	"server/components/command"
	"server/components/session"
)

// Service declares interface for commands processing
type Service interface {
	Process(ctx session.Session, cmd command.Command) error
}
