package executor

import (
	"SSoC/internal/session"
	t "SSoC/internal/token"
	"io"
)

// TokenExecutor responsible for executing "TOKEN <access token>" command;
// should parse and validate client token, should deregister previous client with equal token from session.Storage,
// should set current access token as primary token for current client session
type TokenExecutor struct {
	ctx session.Storage
}

const (
	SuccessResult = "SUCCESS\n"
	FatalResult   = "FATAL\n"
)

// CanAccess always returns true
func (e TokenExecutor) CanAccess(_ string) bool {
	return true
}

// Process receives <access token> from client, parses and validates token,
// deregister previous user with equal token and sets current token for client session
func (e TokenExecutor) Process(writer io.Writer, session session.Session, params ...string) error {
	token := params[0]

	payload, err := t.ParseToken([]byte(token))
	if err != nil {
		return err
	}

	err = t.ValidateToken(payload)
	if err != nil {
		_, err = writer.Write([]byte(FatalResult))
		if err != nil {
			return err
		}

		return err
	}

	if _, err := e.ctx.Find(token); err == nil {
		err := e.ctx.Deregister(token)
		if err != nil {
			_, err = writer.Write([]byte(FatalResult))
			if err != nil {
				return err
			}

			return err
		}
	}

	_, err = writer.Write([]byte(SuccessResult))
	if err != nil {
		return err
	}

	session.SetAccessToken(token)
	e.ctx.Register(session)
	return nil
}

// createTokenExecutor creates TokenExecutor with received context
func createTokenExecutor(ctx session.Storage) Executor {
	return &TokenExecutor{ctx: ctx}
}
