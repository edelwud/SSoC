package executor

import (
	"server/components/session"
	"strings"
)

// RequestUploadExecutor responsible for executing "REQUEST_UPLOAD" command;
// should return unfinished server uploads for current access token
type RequestUploadExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e RequestUploadExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process returns user unfinished uploads
func (e RequestUploadExecutor) Process(session session.Session, _ ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	uploads := s.ReceiveUnfinishedUploads()
	_, err = s.GetConn().Write([]byte(strings.Join(uploads, ",") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// createRequestUploadExecutor creates RequestUploadExecutor with received context
func createRequestUploadExecutor(ctx session.Storage) Executor {
	return &RequestUploadExecutor{ctx: ctx}
}
