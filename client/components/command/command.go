package command

import "main/components/session"

type Command interface {
	Row() []byte
	Process(ctx session.Session) error
}
