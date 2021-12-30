package executor

import (
	"SSoC/internal/command"
	"SSoC/internal/session"
	"io"
)

// Service declares interface for commands processing
type Service interface {
	Process(writer io.Writer, ctx session.Session, cmd command.Command) error
}
