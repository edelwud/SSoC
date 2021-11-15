package executor

import (
	"SSoC/internal/server/datachannel"
	"SSoC/internal/session"
)

// DownloadExecutor responsible for executing "DOWNLOAD <filename>" command
type DownloadExecutor struct {
	File *session.File
	ctx  session.Storage
}

// DownloadFolder folder where stored upload files, which is able to download
const DownloadFolder = "files/uploads"

// CanAccess returns false if current client haven't access token
func (e DownloadExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process executes DOWNLOAD command; receives <filename> from client, registers download in client session,
// receives datachannel port, uses CreateDatachannel for datachannel initialization
func (e *DownloadExecutor) Process(session session.Session, params ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	filename := params[0]
	filepath := DownloadFolder + "/" + params[0]

	e.File, err = s.RegisterDownload(filename, filepath)
	if err != nil {
		return err
	}

	dc := datachannel.NewTCPDatachannel(session.GetOptions())
	err = dc.Listen()
	if err != nil {
		return err
	}

	_, err = session.GetConn().Write([]byte(dc.GetPort() + "\n"))
	if err != nil {
		return err
	}

	err = dc.Accept()
	if err != nil {
		return err
	}

	err = dc.Download(e.File)
	if err != nil {
		return err
	}

	err = dc.Close()
	if err != nil {
		return err
	}

	return nil
}

// createDownloadExecutor creates DownloadExecutor with received context
func createDownloadExecutor(ctx session.Storage) Executor {
	return &DownloadExecutor{
		ctx: ctx,
	}
}
