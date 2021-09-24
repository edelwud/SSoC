package command

import (
	"main/components/session"
)

// EchoCommand responding for construction "ECHO <client text>" command
type EchoCommand struct {
	Cmd  string
	Text string
}

// Row serializes command
func (c EchoCommand) Row() []byte {
	return []byte(c.Cmd + " " + c.Text + "\n")
}

// Process writes serialized command to server
func (c EchoCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateEchoCommand constructs EchoCommand with EchoCommand.Cmd == "ECHO"
func CreateEchoCommand(text string) Command {
	return &EchoCommand{Cmd: "ECHO", Text: text}
}
