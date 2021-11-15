package executor

import (
	"SSoC/internal/command"
	"SSoC/internal/session"
)

// Service declares interface for commands processing
type Service interface {
	Process(ctx session.Session, cmd command.Command) error
}
