package executor

import (
	"server/components/session"
	"strings"
)

type EchoExecutor struct {
	ctx session.SessionStorage
}

func (e EchoExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

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

func createEchoExecutor(ctx session.SessionStorage) Executor {
	return &EchoExecutor{ctx: ctx}
}
