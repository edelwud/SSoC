package executor

import (
	"server/components/session"
	"time"
)

// TimeExecutor responsible for executing "TIME" command;
// should return current server time to client
type TimeExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e TimeExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process returns current server time to client
func (e TimeExecutor) Process(session session.Session, _ ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	_, err = s.GetConn().Write([]byte(time.Now().String() + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// createTimeExecutor creates TimeExecutor with received context
func createTimeExecutor(ctx session.Storage) Executor {
	return &TimeExecutor{ctx: ctx}
}
