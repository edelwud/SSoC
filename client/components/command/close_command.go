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

func (c CloseCommand) AfterExec(_ session.Session) error {
	return nil
}

func CreateCloseCommand() Command {
	return &CloseCommand{Cmd: "CLOSE"}
}
