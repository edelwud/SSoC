package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type TCPDatachannel struct {
	Listener *net.TCPListener
	Conn     *net.TCPConn
	Port     string
	Options  options.Options
}

func generatePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

func (d *TCPDatachannel) Listen() error {
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

func (d *TCPDatachannel) Accept() error {
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

func (d TCPDatachannel) Close() error {
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

func (d TCPDatachannel) GetPort() string {
	return d.Port
}

func (d TCPDatachannel) Download(file *session.File) error {
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

func (d TCPDatachannel) Upload(file *session.File) error {
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

func NewTCPDatachannel(ops options.Options) Datachannel {
	return &TCPDatachannel{
		Port:    generatePort(),
		Options: ops,
	}
}
