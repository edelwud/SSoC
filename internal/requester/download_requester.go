package requester

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// DownloadRequester responding for construction "DOWNLOAD <filename>" command
type DownloadRequester struct {
	Cmd      string
	Filename string
	Filepath string
	File     *session.File
}

// DownloadFolder folder where stored all downloaded files from server
const DownloadFolder = "files/downloads"

// dataChannelReady channel witch indicates that server datachannel listener is ready
var dataChannelReady = make(chan bool, 10)

// GeneratePort generates random port from 8000 to 9000
func GeneratePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

// Row serializes command "DOWNLOAD <filename>"
func (c DownloadRequester) Row() []byte {
	result := []byte(c.Cmd + " " + c.Filename + "\n")
	return result
}

// CreateDatachannel creates datachannel between client and server;
// client acts as serverside with randomly generated port (from 8000 to 9000),
// server acts as clientside witch receives client port and connects to datachannel
func (c DownloadRequester) CreateDatachannel(options options.Options, port string) error {
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

	_, err = io.Copy(c.File, conn)
	if err != nil {
		return err
	}

	err = c.File.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Process registers download, generates port (8000-9000), creates datachannel via CreateDatachannel,
// sends port to server, receives file and stores them to DownloadFolder
func (c DownloadRequester) Process(ctx session.Session) error {
	var err error
	c.File, err = ctx.RegisterDownload(c.Filename, c.Filepath)

	_, err = ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}

	port := GeneratePort()

	go func() {
		<-dataChannelReady
		_, err = ctx.GetConn().Write([]byte(port + "\n"))
		if err != nil {
			return
		}
	}()

	err = c.CreateDatachannel(ctx.GetOptions(), port)
	if err != nil {
		return err
	}

	_, err = ctx.GetConn().Write([]byte("RECEIVED\n"))
	if err != nil {
		return err
	}

	err = c.File.Close()
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
