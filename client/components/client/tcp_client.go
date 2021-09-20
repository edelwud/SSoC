package client

import (
	"main/components/command"
	"main/components/session"
	"main/components/token"
	"net"
	"time"
)

type TcpClient struct {
	Session session.Session
	Options Options
}

func (c *TcpClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Options.Host+":"+c.Options.Port)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlive(c.Options.KeepAlive)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlivePeriod(time.Second * time.Duration(c.Options.KeepAlivePeriod))
	if err != nil {
		return err
	}

	err = c.Auth(conn)
	if err != nil {
		return err
	}

	clientLogger.Info("successfully connected to server")

	return nil
}

func (c *TcpClient) Auth(conn *net.TCPConn) error {
	macToken, err := token.GenerateMACToken()
	if err != nil {
		return err
	}

	t, err := macToken.Row()
	if err != nil {
		return err
	}

	c.Session = session.CreateBasicSession(conn, t)

	cmd := command.CreateTokenCommand(macToken)

	err = c.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c TcpClient) Disconnect() error {
	if c.Session.GetConn() == nil {
		return nil
	}

	err := c.Session.GetConn().Close()
	if err != nil {
		return err
	}

	return nil
}

func (c TcpClient) Exec(cmd command.Command) error {
	err := cmd.Process(c.Session)
	if err != nil {
		return err
	}

	clientLogger.Info("command executed")

	return nil
}

func (c TcpClient) Write(cmd string) error {
	_, err := c.Session.GetConn().Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}

func (c TcpClient) GetContext() session.Session {
	return c.Session
}

func CreateTcpClient(options Options) Client {
	return &TcpClient{Options: options}
}
