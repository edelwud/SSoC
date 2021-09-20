package command

import "main/components/session"

type Command interface {
	Row() []byte
	AfterExec(ctx session.Session) error
}
