package command

import (
	"main/components/session"
)

// CloseCommand responding for construction "CLOSE" command
type CloseCommand struct {
	Cmd string
}

// Row serializes command "CLOSE"
func (c CloseCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c CloseCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateCloseCommand constructs CloseCommand with CloseCommand.Cmd == "CLOSE"
func CreateCloseCommand() Command {
	return &CloseCommand{Cmd: "CLOSE"}
}
