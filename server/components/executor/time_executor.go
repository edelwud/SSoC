package executor

import (
	"server/components/session"
	"time"
)

type TimeExecutor struct {
	ctx session.SessionStorage
}

func (e TimeExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

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

func createTimeExecutor(ctx session.SessionStorage) Executor {
	return &TimeExecutor{ctx: ctx}
}
