package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"net"
)

type UDPDatachannelServer struct {
	Conn    *net.UDPConn
	Port    string
	Options options.Options
}

func (d *UDPDatachannelServer) Listen() error {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+d.Port)
	if err != nil {
		return err
	}

	d.Conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	return nil
}

func (d *UDPDatachannelServer) Accept() error {
	return nil
}

func (d UDPDatachannelServer) Close() error {
	err := d.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d UDPDatachannelServer) GetPort() string {
	return d.Port
}

func (d UDPDatachannelServer) Download(file *session.File) error {
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

func (d UDPDatachannelServer) Upload(file *session.File) error {
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

func (d UDPDatachannelServer) Connect() error {
	return nil
}

func NewUDPDatachannelServer(ops options.Options) Datachannel {
	return &UDPDatachannelServer{
		Port:    generatePort(),
		Options: ops,
	}
}
