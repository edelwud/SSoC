package requester

import (
	"SSoC/internal/client/datachannel"
	"SSoC/internal/session"
	"bufio"
	"net"
	"strings"
)

// UploadRequester responding for construction "UPLOAD <filename>" command
type UploadRequester struct {
	Cmd      string
	Filename string
	Filepath string
}

// Row serializes command
func (r UploadRequester) Row() []byte {
	result := []byte(r.Cmd + " " + r.Filename + "\n")
	return result
}

func (r UploadRequester) ReceivePort(conn net.Conn) (string, error) {
	port, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	port = strings.Trim(port, "\n")
	port = strings.Trim(port, " ")

	return port, nil
}

// Process registers upload, writes command to server, receives datachannel port,
// initializes datachannel, writes received file to uploads
func (r *UploadRequester) Process(ctx session.Session) error {
	file, err := ctx.RegisterUpload(r.Filename, r.Filepath)
	if err != nil {
		return err
	}

	_, err = ctx.GetConn().Write(r.Row())
	if err != nil {
		return err
	}

	port, err := r.ReceivePort(ctx.GetConn())

	dc := datachannel.NewTCPDatachannel(port, ctx.GetOptions())
	err = dc.Connect()
	if err != nil {
		return err
	}

	err = dc.Upload(file)
	if err != nil {
		return err
	}

	err = dc.Close()
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadRequester constructs UploadRequester
func CreateUploadRequester(filename string, path string) Requester {
	return &UploadRequester{
		Cmd:      "UPLOAD",
		Filename: filename,
		Filepath: path,
	}
}
