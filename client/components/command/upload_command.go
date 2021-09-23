package command

import (
	"bufio"
	"io"
	"main/components/session"
	"net"
	"strconv"
	"strings"
	"time"
)

type UploadCommand struct {
	Cmd      string
	Filename string
	Filepath string
	File     *session.File
}

func (c UploadCommand) Row() []byte {
	result := []byte(c.Cmd + " " + c.Filename + " " + strconv.Itoa(int(c.File.Size)) + "\n")
	return result
}

func (c UploadCommand) CreateDatachannel(port string) error {
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

	_, err = io.Copy(conn, c.File)
	if err != nil {
		return err
	}

	err = c.File.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (c *UploadCommand) Process(ctx session.Session) error {
	var err error
	c.File, err = ctx.RegisterUpload(c.Filename, c.Filepath)
	if err != nil {
		return err
	}

	_, err = ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}

	port, err := bufio.NewReader(ctx.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	port = strings.Trim(port, "\n")
	port = strings.Trim(port, " ")

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

func CreateUploadCommand(filename string, path string) Command {

	return &UploadCommand{
		Cmd:      "UPLOAD",
		Filename: filename,
		Filepath: path,
	}
}
