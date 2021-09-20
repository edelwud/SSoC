package executor

import (
	"server/components/session"
	"strings"
)

type DownloadExecutor struct {
	bufferSize int
	ctx        session.SessionStorage
}

func (e DownloadExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

func (e DownloadExecutor) Process(session session.Session, params ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	//filename := params[0]
	//
	//file := s.RegisterDownload()

	_, err = s.GetConn().Write([]byte(strings.Join(params, " ") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func createDownloadExecutor(ctx session.SessionStorage) Executor {
	return &DownloadExecutor{ctx: ctx}
}
