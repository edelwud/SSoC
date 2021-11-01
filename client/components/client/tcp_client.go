package client

import (
	"main/components/command"
	"main/components/options"
	"main/components/session"
	"main/components/token"
	"net"
	"time"
)

// TCPClient implementation of Client interface based on TCP protocol
type TCPClient struct {
	Session session.Session
	Options options.Options
}

// Connect resolves server options from Options, dials via net.DialTCP with TCPv4 background,
// enables/disables keep alive and sets keep alive period from Options, generates and sends access token via Auth
func (c *TCPClient) Connect() error {
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

// Auth generates MAC address based access token, creates client session, sends access token to server
func (c *TCPClient) Auth(conn *net.TCPConn) error {
	macToken, err := token.GenerateMACToken()
	if err != nil {
		return err
	}

	t, err := macToken.Row()
	if err != nil {
		return err
	}

	c.Session = session.CreateClientSession(conn, c.Options, t)

	cmd := command.CreateTokenCommand(macToken)

	err = c.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect closes TCP connection
func (c TCPClient) Disconnect() error {
	if c.Session.GetConn() == nil {
		return nil
	}

	err := c.Session.GetConn().Close()
	if err != nil {
		return err
	}

	return nil
}

// Exec executes received command
func (c TCPClient) Exec(cmd command.Command) error {
	err := cmd.Process(c.Session)
	if err != nil {
		return err
	}

	return nil
}

// Write writes message for server
func (c TCPClient) Write(cmd string) error {
	_, err := c.Session.GetConn().Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}

// GetContext receives client context
func (c TCPClient) GetContext() session.Session {
	return c.Session
}

// CreateTCPClient constructs TCPClient with received Options
func CreateTCPClient(options options.Options) Client {
	return &TCPClient{Options: options}
}
