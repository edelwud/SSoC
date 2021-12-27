package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"net"
	"time"
)

type TCPDatachannelClient struct {
	Conn    *net.TCPConn
	Port    string
	Options options.Options
}

func (d *TCPDatachannelClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", d.Options.Host+":"+d.Port)
	if err != nil {
		return err
	}

	d.Conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	err = d.Conn.SetKeepAlive(d.Options.KeepAlive)
	if err != nil {
		return err
	}

	err = d.Conn.SetKeepAlivePeriod(time.Duration(d.Options.KeepAlivePeriod) * time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (d TCPDatachannelClient) Close() error {
	err := d.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d TCPDatachannelClient) Upload(file *session.File) error {
	_, err := io.Copy(d.Conn, file)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d TCPDatachannelClient) Download(file *session.File) error {
	_, err := io.Copy(file, d.Conn)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d TCPDatachannelClient) Accept() error {
	return nil
}

func (d TCPDatachannelClient) Listen() error {
	return nil
}

func (d TCPDatachannelClient) GetPort() string {
	return d.Port
}

func NewTCPDatachannelClient(port string, ops options.Options) Datachannel {
	return &TCPDatachannelClient{
		Port:    port,
		Options: ops,
	}
}
