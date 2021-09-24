package executor

import (
	"server/components/session"
)

// CloseExecutor responsible for executing "CLOSE" command;
// should close current connection
type CloseExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e CloseExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process uses session.Storage for closing down client connection
func (e CloseExecutor) Process(session session.Session, _ ...string) error {
	err := e.ctx.Deregister(session.GetAccessToken())
	if err != nil {
		return err
	}

	return nil
}

// createCloseExecutor creates CloseExecutor with received context
func createCloseExecutor(ctx session.Storage) Executor {
	return &CloseExecutor{ctx: ctx}
}
