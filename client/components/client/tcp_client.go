package client

import (
	"main/components/command"
	"net"
	"time"
)

type TcpClient struct {
	Conn    *net.TCPConn
	Options Options
}

func (c *TcpClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Options.Host+":"+c.Options.Port)
	if err != nil {
		return err
	}

	c.Conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	err = c.Conn.SetKeepAlive(c.Options.KeepAlive)
	if err != nil {
		return err
	}

	err = c.Conn.SetKeepAlivePeriod(time.Second * time.Duration(c.Options.KeepAlivePeriod))
	if err != nil {
		return err
	}

	clientLogger.Info("successfully connected to server")

	return nil
}

func (c TcpClient) Disconnect() error {
	if c.Conn == nil {
		return nil
	}

	err := c.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c TcpClient) Exec(cmd command.Command) error {
	w, err := c.Conn.Write(cmd.Row())
	if err != nil {
		return err
	}

	clientLogger.Infof("command executed, written %d bytes", w)

	return nil
}

func CreateTcpClient(options Options) Client {
	return &TcpClient{Options: options}
}
