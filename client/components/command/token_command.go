package command

import (
	"main/components/session"
	"main/components/token"
)

type TokenCommand struct {
	Cmd   string
	Token string
}

func (c TokenCommand) Row() []byte {
	return []byte(c.Cmd + " " + c.Token + "\n")
}

func (c TokenCommand) AfterExec(ctx session.Session) error {
	return nil
}

func CreateTokenCommand(payload token.Payload) Command {
	row, err := payload.Row()
	if err != nil {
		return nil
	}
	return &TokenCommand{Cmd: "TOKEN", Token: row}
}
