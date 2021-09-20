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

func (c EchoCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

func CreateEchoCommand(text string) Command {
	return &EchoCommand{Cmd: "ECHO", Text: text}
}
