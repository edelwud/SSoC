package requester

import (
	"SSoC/internal/session"
	"bufio"
	"errors"
	"strings"
)

// TokenRequester responding for construction "TOKEN <access token>" command
type TokenRequester struct {
	Cmd   string
	Token string
}

const SuccessResult = "SUCCESS"

// Row serializes command
func (c TokenRequester) Row() []byte {
	return []byte(c.Cmd + " " + c.Token + "\n")
}

// Process writes serialized command to server
func (c TokenRequester) Process(ctx session.Session) error {
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

// CreateTokenRequester constructs TokenRequester
func CreateTokenRequester(accessToken string) Requester {
	return &TokenRequester{Cmd: "TOKEN", Token: accessToken}
}
