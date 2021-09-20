package command

import (
	"bufio"
	"main/components/session"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type DownloadCommand struct {
	Cmd      string
	Filename string
}

const DownloadFolder = "files/downloads"

var dataChannelReady = make(chan bool, 10)

func GeneratePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

func (c DownloadCommand) Row() []byte {
	result := []byte(c.Cmd + " " + c.Filename + "\n")
	return result
}

func (c DownloadCommand) CreateDatachannel(port string, filename string) (int, error) {
	fileHandler, err := os.Create(DownloadFolder + "/" + filename)
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

func (c DownloadCommand) Process(ctx session.Session) error {
	download := ctx.RegisterDownload()
	download.Filename = c.Filename
	download.StartTime = time.Now()

	_, err := ctx.GetConn().Write(c.Row())
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

	n, err := c.CreateDatachannel(port, c.Filename)
	if err != nil {
		return err
	}

	_, err = bufio.NewReader(ctx.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	download.Transferred = n
	download.EndTime = time.Now()

	return nil
}

func CreateDownloadCommand(filename string) Command {
	return &DownloadCommand{
		Cmd:      "DOWNLOAD",
		Filename: filename,
	}
}
