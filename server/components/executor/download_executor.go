package executor

import (
	"bufio"
	"io"
	"net"
	"server/components/session"
	"strings"
	"time"
)

const DownloadFolder = "files/uploads"

type DownloadExecutor struct {
	File *session.File
	ctx  session.SessionStorage
}

func (e DownloadExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

func (e DownloadExecutor) CreateDatachannel(port string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
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

	err = e.CreateDatachannel(port)
	if err != nil {
		return err
	}

	_, err = s.GetConn().Write([]byte("RECEIVED\n"))
	if err != nil {
		return err
	}

	err = e.File.Close()
	if err != nil {
		return err
	}

	return nil
}

func createDownloadExecutor(ctx session.SessionStorage) Executor {
	return &DownloadExecutor{ctx: ctx}
}
