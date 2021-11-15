package executor

import (
	"SSoC/internal/session"
)

// Executor declares interfaces for command execution
type Executor interface {
	Process(session session.Session, params ...string) error
	CanAccess(accessToken string) bool
}