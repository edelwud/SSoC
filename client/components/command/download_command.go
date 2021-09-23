package command

import (
	"bufio"
	"main/components/session"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type DownloadCommand struct {
	Cmd      string
	Filename string
	Filepath string
	File     *session.File
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

func (c DownloadCommand) CreateDatachannel(port string) error {
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

	reader := bufio.NewReader(conn)

	_, err = reader.WriteTo(c.File)
	if err != nil {
		return err
	}

	err = c.File.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (c DownloadCommand) Process(ctx session.Session) error {
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

	err = c.CreateDatachannel(port)
	if err != nil {
		return err
	}

	_, err = bufio.NewReader(ctx.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	err = c.File.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreateDownloadCommand(filename string) Command {
	return &DownloadCommand{
		Cmd:      "DOWNLOAD",
		Filename: filename,
		Filepath: DownloadFolder + "/" + filename,
	}
}
