package command

import (
	"bufio"
	"errors"
	"main/components/session"
	"strings"
)

// TokenCommand responding for construction "TOKEN <access token>" command
type TokenCommand struct {
	Cmd   string
	Token string
}

const SuccessResult = "SUCCESS"

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

	buf, _, err := bufio.NewReader(ctx.GetConn()).ReadLine()
	if err != nil {
		return err
	}

	result := strings.Trim(string(buf), "\n")
	result = strings.Trim(result, " ")

	if result != SuccessResult {
		return errors.New("token was expired")
	}

	return nil
}

// CreateTokenCommand constructs TokenCommand
func CreateTokenCommand(accessToken string) Command {
	return &TokenCommand{Cmd: "TOKEN", Token: accessToken}
}
