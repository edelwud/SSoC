package executor

import (
	"bufio"
	"math/rand"
	"net"
	"os"
	"server/components/session"
	"strconv"
	"time"
)

const UploadFolder = "files/uploads"

type UploadExecutor struct {
	bufferSize int
	ctx        session.SessionStorage
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

func (e UploadExecutor) CreateDatachannel(port string, filename string) (int, error) {
	fileHandler, err := os.Create(UploadFolder + "/" + filename)
	if err != nil {
		return 0, err
	}

	defer func(fileHandler *os.File) {
		err := fileHandler.Close()
		if err != nil {
			return
		}
	}(fileHandler)

	writer := bufio.NewWriter(fileHandler)

	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	defer func(listener *net.TCPListener) {
		err := listener.Close()
		if err != nil {
			return
		}
	}(listener)
	if err != nil {
		return 0, err
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
		return 0, err
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		return 0, err
	}

	err = conn.SetKeepAlivePeriod(360 * time.Second)
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(conn)

	n, err := reader.WriteTo(writer)
	if err != nil {
		return 0, err
	}

	err = writer.Flush()
	if err != nil {
		return 0, err
	}

	err = fileHandler.Sync()
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func (e UploadExecutor) Process(session session.Session, params ...string) error {
	s, err := e.ctx.Find(session.GetAccessToken())
	if err != nil {
		return err
	}

	file := s.RegisterUpload()
	file.Filename = params[0]
	file.StartTime = time.Now()

	port := GeneratePort()

	go func() {
		<-dataChannelReady
		_, err = s.GetConn().Write([]byte(port + "\n"))
		if err != nil {
			return
		}
	}()

	n, err := e.CreateDatachannel(port, file.Filename)
	if err != nil {
		return err
	}

	_, err = s.GetConn().Write([]byte("RECEIVED\n"))
	if err != nil {
		return err
	}

	file.Transferred = n
	file.EndTime = time.Now()

	return nil
}

func createUploadExecutor(ctx session.SessionStorage) Executor {
	return &UploadExecutor{ctx: ctx}
}
