package executor

import (
	"server/components/session"
)

type CloseExecutor struct {
	ctx session.SessionStorage
}

func (e CloseExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

func (e CloseExecutor) Process(session session.Session, _ ...string) error {
	err := e.ctx.Deregister(session.GetAccessToken())
	if err != nil {
		return err
	}

	return nil
}

func createCloseExecutor(ctx session.SessionStorage) Executor {
	return &CloseExecutor{ctx: ctx}
}
