package tcp_client

import (
	"SSoC/internal/client"
	"SSoC/internal/options"
	"SSoC/internal/requester"
	"SSoC/internal/session"
	clientSession "SSoC/internal/session/client_session"
	"github.com/sirupsen/logrus"
	"net"
)

// UDPClient implementation of Client interface based on TCP protocol
type UDPClient struct {
	AccessToken string
	Session     session.Session
	Options     options.Options
}

var clientLogger = logrus.WithField("context", "client")

// Connect resolves server options from Options, dials via net.DialTCP with TCPv4 background,
// enables/disables keep alive and sets keep alive period from Options, generates and sends access token via Auth
func (c *UDPClient) Connect() error {
	udpAddr, err := net.ResolveUDPAddr("udp", c.Options.Host+":"+c.Options.Port)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
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
func (c *UDPClient) Auth(conn *net.UDPConn) error {
	c.Session = clientSession.CreateClientSession(conn, c.Options, c.AccessToken)
	req := requester.CreateTokenRequester(c.AccessToken)

	err := c.Exec(req)
	if err != nil {
		return err
	}

	reqDownload := requester.CreateRequestDownloadRequester()
	err = c.Exec(reqDownload)
	if err != nil {
		return err
	}

	reqUpload := requester.CreateRequestUploadRequester()
	err = c.Exec(reqUpload)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect closes TCP connection
func (c UDPClient) Disconnect() error {
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
func (c UDPClient) Exec(cmd requester.Requester) error {
	err := cmd.Process(c.Session)
	if err != nil {
		return err
	}

	return nil
}

// Write writes message for server
func (c UDPClient) Write(cmd string) error {
	_, err := c.Session.GetConn().Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}

// GetContext receives client context
func (c UDPClient) GetContext() session.Session {
	return c.Session
}

// CreateUDPClient constructs UDPClient with received Options
func CreateUDPClient(options options.Options, accessToken string) client.Client {
	return &UDPClient{
		Options:     options,
		AccessToken: accessToken,
	}
}
