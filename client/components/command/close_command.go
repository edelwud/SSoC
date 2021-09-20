package command

import (
	"main/components/session"
)

type CloseCommand struct {
	Cmd string
}

func (c CloseCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

func (c CloseCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

func CreateCloseCommand() Command {
	return &CloseCommand{Cmd: "CLOSE"}
}
