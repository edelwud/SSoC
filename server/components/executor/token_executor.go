package executor

import (
	"server/components/session"
	t "server/components/token"
)

type TokenExecutor struct {
	ctx session.SessionStorage
}

func (e TokenExecutor) CanAccess(_ string) bool {
	return true
}

func (e TokenExecutor) Process(session session.Session, params ...string) error {
	token := params[0]

	payload, err := t.ParseToken([]byte(token))
	if err != nil {
		return err
	}

	err = t.ValidateToken(payload)
	if err != nil {
		return err
	}

	if _, err := e.ctx.Find(token); err == nil {
		err := e.ctx.Deregister(token)
		if err != nil {
			return err
		}
	}

	session.SetAccessToken(token)
	e.ctx.Register(session)
	return nil
}

func createTokenExecutor(ctx session.SessionStorage) Executor {
	return &TokenExecutor{ctx: ctx}
}
