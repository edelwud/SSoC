package executor

import (
	"bufio"
	"io"
	"net"
	"server/components/options"
	"server/components/session"
	"strings"
	"time"
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

// CreateDatachannel creates datachannel between client and server;
// client acts as serverside with randomly generated port (from 8000 to 9000),
// server acts as clientside witch receives client port and connects to datachannel
func (e DownloadExecutor) CreateDatachannel(options options.Options, port string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", options.Host+":"+port)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	err = conn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlivePeriod(360 * time.Second)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, e.File)
	if err != nil {
		return err
	}

	err = e.File.Sync()
	if err != nil {
		return err
	}

	return nil
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

	port, err := bufio.NewReader(s.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	port = strings.Trim(port, "\n")
	port = strings.Trim(port, " ")

	err = e.CreateDatachannel(session.GetOptions(), port)
	if err != nil {
		return err
	}

	_, err = bufio.NewReader(s.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	err = e.File.Close()
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
