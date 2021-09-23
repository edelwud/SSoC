package executor

import (
	"io"
	"math/rand"
	"net"
	"server/components/session"
	"strconv"
	"time"
)

const UploadFolder = "files/uploads"

type UploadExecutor struct {
	File *session.File
	ctx  session.SessionStorage
}

var dataChannelReady = make(chan bool, 10)

func (e UploadExecutor) CanAccess(accessToken string) bool {
	if accessToken == "" {
		return false
	}
	return true
}

func GeneratePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

func (e UploadExecutor) CreateDatachannel(port string) error {
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer func(listener *net.TCPListener) {
		err := listener.Close()
		if err != nil {
			return
		}
	}(listener)

	dataChannelReady <- true

	conn, err := listener.AcceptTCP()
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

func createUploadExecutor(ctx session.SessionStorage) Executor {
	return &UploadExecutor{ctx: ctx}
}
