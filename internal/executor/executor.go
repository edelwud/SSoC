package executor

import (
	"SSoC/internal/session"
	"io"
)

// Executor declares interfaces for command execution
type Executor interface {
	Process(writer io.Writer, session session.Session, params ...string) error
	CanAccess(accessToken string) bool
}
