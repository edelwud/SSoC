package command

import (
	"main/components/session"
)

type EchoCommand struct {
	Cmd  string
	Text string
}

func (c EchoCommand) Row() []byte {
	return []byte(c.Cmd + " " + c.Text + "\n")
}

func (c EchoCommand) AfterExec(_ session.Session) error {
	return nil
}

func CreateEchoCommand(text string) Command {
	return &EchoCommand{Cmd: "ECHO", Text: text}
}
