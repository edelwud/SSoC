package executor

import (
	"server/components/session"
	"strings"
)

// EchoExecutor responsible for executing "ECHO <client string>" command;
// should return <client string> to client back
type EchoExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e EchoExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

// Process writes to current client <client string>
func (e EchoExecutor) Process(session session.Session, params ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	_, err = s.GetConn().Write([]byte(strings.Join(params, " ") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// createEchoExecutor creates EchoExecutor with received context
func createEchoExecutor(ctx session.Storage) Executor {
	return &EchoExecutor{ctx: ctx}
}
