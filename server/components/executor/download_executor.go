package executor

import (
	"io"
	"net"
	"server/components/session"
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
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	defer func(listener *net.TCPListener) {
		err := listener.Close()
		if err != nil {
			return
		}
	}(listener)
	if err != nil {
		return err
	}

	dataChannelReady <- true

	conn, err := listener.AcceptTCP()
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlivePeriod(360 * time.Second)
	if err != nil {
		return err
	}

	_, err = io.Copy(e.File, conn)
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

	port := GeneratePort()

	go func() {
		<-dataChannelReady
		_, err = s.GetConn().Write([]byte(port + "\n"))
		if err != nil {
			return
		}
	}()

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
