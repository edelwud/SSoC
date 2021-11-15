package executor

import (
	"SSoC/internal/server/datachannel"
	"SSoC/internal/session"
)

// UploadExecutor responsible for executing "UPLOAD <filename>" command
type UploadExecutor struct {
	File *session.File
	ctx  session.Storage
}

// UploadFolder folder where stores all clients uploads
const UploadFolder = "files/uploads"

// CanAccess returns false if current client haven't access token
func (e UploadExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// Process executes UPLOAD command; receives <filename> from client, registers upload in client session,
// sends generated port, performs datachannel connection with client, writes file to UploadFolder
func (e *UploadExecutor) Process(session session.Session, params ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	filename := params[0]
	filepath := UploadFolder + "/" + filename

	e.File, err = s.RegisterUpload(filename, filepath)
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

	err = dc.Upload(e.File)
	if err != nil {
		return err
	}

	err = dc.Close()
	if err != nil {
		return err
	}

	return nil
}

// createUploadExecutor creates UploadExecutor with received context
func createUploadExecutor(ctx session.Storage) Executor {
	return &UploadExecutor{
		ctx: ctx,
	}
}
