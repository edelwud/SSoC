package executor

import (
	"server/components/command"
	"server/components/session"
)

type ExecService interface {
	Process(ctx session.Session, cmd command.Command) error
}
