package executor

import (
	"SSoC/internal/session"
	"io"
	"strings"
)

// EchoExecutor responsible for executing "ECHO <client string>" command;
// should return <client string> to client back
type EchoExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e EchoExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process writes to current client <client string>
func (e EchoExecutor) Process(writer io.Writer, session session.Session, params ...string) error {
	_, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(strings.Join(params, " ") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// createEchoExecutor creates EchoExecutor with received context
func createEchoExecutor(ctx session.Storage) Executor {
	return &EchoExecutor{ctx: ctx}
}
