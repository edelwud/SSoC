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

func (c TokenCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

func CreateTokenCommand(payload token.Payload) Command {
	row, err := payload.Row()
	if err != nil {
		return nil
	}
	return &TokenCommand{Cmd: "TOKEN", Token: row}
}
