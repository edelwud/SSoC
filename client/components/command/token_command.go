package command

import (
	"main/components/session"
	"main/components/token"
)

// TokenCommand responding for construction "TOKEN <access token>" command
type TokenCommand struct {
	Cmd   string
	Token string
}

// Row serializes command
func (c TokenCommand) Row() []byte {
	return []byte(c.Cmd + " " + c.Token + "\n")
}

// Process writes serialized command to server
func (c TokenCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateTokenCommand constructs TokenCommand
func CreateTokenCommand(payload token.Payload) Command {
	row, err := payload.Row()
	if err != nil {
		return nil
	}
	return &TokenCommand{Cmd: "TOKEN", Token: row}
}
