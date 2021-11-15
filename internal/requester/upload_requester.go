package requester

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// UploadRequester responding for construction "UPLOAD <filename>" command
type UploadRequester struct {
	Cmd      string
	Filename string
	Filepath string
	File     *session.File
}

// Row serializes command
func (c UploadRequester) Row() []byte {
	result := []byte(c.Cmd + " " + c.Filename + " " + strconv.Itoa(int(c.File.Size)) + "\n")
	return result
}

// CreateDatachannel creates datachannel between server and client;
// server performs datachannel listener with randomly generated port (from 8000 to 9000),
// client receives datachannel port and performs TCP connection to server
func (c UploadRequester) CreateDatachannel(options options.Options, port string) error {
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

// Process registers upload, writes command to server, receives datachannel port,
// initializes datachannel, writes received file to uploads
func (c *UploadRequester) Process(ctx session.Session) error {
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

	err = c.CreateDatachannel(ctx.GetOptions(), port)
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

// CreateUploadRequester constructs UploadRequester
func CreateUploadRequester(filename string, path string) Requester {
	return &UploadRequester{
		Cmd:      "UPLOAD",
		Filename: filename,
		Filepath: path,
	}
}
