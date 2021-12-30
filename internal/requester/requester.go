package requester

import (
	"SSoC/internal/session"
)

// Requester declares interface for client commands
type Requester interface {
	Row() []byte
	Process(ctx session.Session) error
}
