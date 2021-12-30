package requester

import (
	"SSoC/internal/datachannel"
	"SSoC/internal/session"
	"bufio"
	"net"
	"strings"
)

// DownloadRequester responding for construction "DOWNLOAD <filename>" command
type DownloadRequester struct {
	Cmd      string
	Filename string
	Filepath string
}

// DownloadFolder folder where stored all downloaded files from server
const DownloadFolder = "files/downloads"

// Row serializes command "DOWNLOAD <filename>"
func (r DownloadRequester) Row() []byte {
	result := []byte(r.Cmd + " " + r.Filename + "\n")
	return result
}

func (r DownloadRequester) ReceivePort(conn net.Conn) (string, error) {
	port, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	port = strings.Trim(port, "\n")
	port = strings.Trim(port, " ")

	return port, nil
}

// Process registers download, generates port (8000-9000), creates datachannel via CreateDatachannel,
// sends port to server, receives file and stores them to DownloadFolder
func (r DownloadRequester) Process(ctx session.Session) error {
	file, err := ctx.RegisterDownload(r.Filename, r.Filepath)

	_, err = ctx.GetConn().Write(r.Row())
	if err != nil {
		return err
	}

	port, err := r.ReceivePort(ctx.GetConn())

	dc := datachannel.New("client", ctx.GetOptions().Protocol, port, ctx.GetOptions())
	err = dc.Connect()
	if err != nil {
		return err
	}

	err = dc.Download(file)
	if err != nil {
		return err
	}

	err = dc.Close()
	if err != nil {
		return err
	}

	return nil
}

// CreateDownloadRequester constructs DownloadRequester
func CreateDownloadRequester(filename string) Requester {
	return &DownloadRequester{
		Cmd:      "DOWNLOAD",
		Filename: filename,
		Filepath: DownloadFolder + "/" + filename,
	}
}
