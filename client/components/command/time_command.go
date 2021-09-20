package command

import (
	"main/components/session"
)

type TimeCommand struct {
	Cmd string
}

func (c TimeCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

func (c TimeCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

func CreateTimeCommand() Command {
	return &TimeCommand{Cmd: "TIME"}
}
