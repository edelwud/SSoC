package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"net"
	"time"
)

type TCPDatachannelServer struct {
	Listener *net.TCPListener
	Conn     *net.TCPConn
	Port     string
	Options  options.Options
}

func (d *TCPDatachannelServer) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+d.Port)
	if err != nil {
		return err
	}

	d.Listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	return nil
}

func (d *TCPDatachannelServer) Accept() error {
	var err error

	d.Conn, err = d.Listener.AcceptTCP()
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

func (d TCPDatachannelServer) Close() error {
	err := d.Listener.Close()
	if err != nil {
		return err
	}

	err = d.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d TCPDatachannelServer) GetPort() string {
	return d.Port
}

func (d TCPDatachannelServer) Download(file *session.File) error {
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

func (d TCPDatachannelServer) Upload(file *session.File) error {
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

func (d TCPDatachannelServer) Connect() error {
	return nil
}

func NewTCPDatachannelServer(ops options.Options) Datachannel {
	return &TCPDatachannelServer{
		Port:    generatePort(),
		Options: ops,
	}
}
