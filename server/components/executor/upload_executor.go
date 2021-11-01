package executor

import (
	"io"
	"math/rand"
	"net"
	"server/components/options"
	"server/components/session"
	"strconv"
	"time"
)

// UploadExecutor responsible for executing "UPLOAD <filename>" command
type UploadExecutor struct {
	File *session.File
	ctx  session.Storage
}

// UploadFolder folder where stores all clients uploads
const UploadFolder = "files/uploads"

// dataChannelReady channel witch indicates that server datachannel listener is ready
var dataChannelReady = make(chan bool, 10)

// CanAccess returns false if current client haven't access token
func (e UploadExecutor) CanAccess(accessToken string) bool {
	return accessToken != ""
}

// GeneratePort generates random port from 8000 to 9000
func GeneratePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

// CreateDatachannel creates datachannel between server and client;
// server performs datachannel listener with randomly generated port (from 8000 to 9000),
// client receives datachannel port and performs TCP connection to server
func (e UploadExecutor) CreateDatachannel(options options.Options, port string) error {
	addr, err := net.ResolveTCPAddr("tcp", options.Host+":"+port)
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

	port := GeneratePort()

	go func() {
		<-dataChannelReady
		_, err = s.GetConn().Write([]byte(port + "\n"))
		if err != nil {
			return
		}
	}()

	err = e.CreateDatachannel(session.GetOptions(), port)
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

// createUploadExecutor creates UploadExecutor with received context
func createUploadExecutor(ctx session.Storage) Executor {
	return &UploadExecutor{
		ctx: ctx,
	}
}
