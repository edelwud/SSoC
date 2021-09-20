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

func (c TimeCommand) AfterExec(_ session.Session) error {
	return nil
}

func CreateTimeCommand() Command {
	return &TimeCommand{Cmd: "TIME"}
}
