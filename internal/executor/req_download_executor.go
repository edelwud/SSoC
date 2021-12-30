package executor

import (
	"SSoC/internal/session"
	"io"
	"strings"
)

// RequestDownloadExecutor responsible for executing "REQUEST_DOWNLOAD" command;
// should return unfinished server downloads for current access token
type RequestDownloadExecutor struct {
	ctx session.Storage
}

// CanAccess returns false if current client haven't access token
func (e RequestDownloadExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process returns user unfinished downloads
func (e RequestDownloadExecutor) Process(writer io.Writer, session session.Session, _ ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	downloads := s.ReceiveUnfinishedDownloads()
	_, err = writer.Write([]byte(strings.Join(downloads, ",") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// createRequestDownloadExecutor creates RequestDownloadExecutor with received context
func createRequestDownloadExecutor(ctx session.Storage) Executor {
	return &RequestDownloadExecutor{ctx: ctx}
}
