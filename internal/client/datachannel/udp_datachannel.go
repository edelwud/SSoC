package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"io"
	"net"
)

type UDPDatachannel struct {
	Conn    *net.UDPConn
	Port    string
	Options options.Options
}

func (d *UDPDatachannel) Connect() error {
	udpAddr, err := net.ResolveUDPAddr("udp", d.Options.Host+":"+d.Port)
	if err != nil {
		return err
	}

	d.Conn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	return nil
}

func (d UDPDatachannel) Close() error {
	err := d.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d UDPDatachannel) Upload(file *session.File) error {
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

func (d UDPDatachannel) Download(file *session.File) error {
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

func NewUDPDatachannel(port string, ops options.Options) Datachannel {
	return &UDPDatachannel{
		Port:    port,
		Options: ops,
	}
}
