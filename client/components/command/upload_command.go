package command

import (
	"bufio"
	"io"
	"main/components/session"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type UploadCommand struct {
	Cmd      string
	Reader   io.Reader
	Filename string
	Size     int
}

func (c UploadCommand) Row() []byte {
	result := []byte(c.Cmd + " " + c.Filename + " " + strconv.Itoa(c.Size) + "\n")
	return result
}

func (c UploadCommand) CreateDatachannel(port string, reader io.Reader) (int, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return 0, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return 0, err
	}
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	dataWriter := bufio.NewWriter(conn)
	n, err := io.Copy(dataWriter, reader)
	if err != nil {
		return 0, err
	}

	err = dataWriter.Flush()
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func (c UploadCommand) Process(ctx session.Session) error {
	upload := ctx.RegisterUpload()
	upload.Filename = c.Filename
	upload.Size = c.Size
	upload.StartTime = time.Now()

	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}

	port, err := bufio.NewReader(ctx.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	port = strings.Trim(port, "\n")
	port = strings.Trim(port, " ")

	n, err := c.CreateDatachannel(port, c.Reader)
	if err != nil {
		return err
	}

	_, err = bufio.NewReader(ctx.GetConn()).ReadString('\n')
	if err != nil {
		return err
	}

	upload.Transferred = n
	upload.EndTime = time.Now()

	return nil
}

func CreateUploadCommand(filename string, path string) Command {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	reader := bufio.NewReader(file)

	return &UploadCommand{
		Cmd:      "UPLOAD",
		Filename: filename,
		Reader:   file,
		Size:     reader.Size(),
	}
}
