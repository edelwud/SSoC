package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type UDPDatachannelServer struct {
	Conn    *net.UDPConn
	Addr    net.Addr
	Port    string
	Options options.Options
}

func (d UDPDatachannelServer) Read(p []byte) (n int, err error) {
	n, _, err = d.Conn.ReadFromUDP(p)
	if binary.BigEndian.Uint64(p[n-8:n]) == EndOfFile {
		return 0, io.EOF
	}
	return n, err
}

func (d UDPDatachannelServer) Write(p []byte) (n int, err error) {
	n, err = d.Conn.WriteTo(p, d.Addr)
	if err != nil {
		return 0, err
	}

	return n, err
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
	acceptanceBuffer := make([]byte, 1024)
	n, addr, err := d.Conn.ReadFrom(acceptanceBuffer)
	if err != nil {
		return err
	}

	if string(acceptanceBuffer[:n-1]) != HelloMessage {
		return errors.New("acceptance message not valid")
	}

	d.Addr = addr
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
	_, err := io.Copy(d, file)
	if err != nil {
		return err
	}

	endOfFile := make([]byte, 8)
	binary.BigEndian.PutUint64(endOfFile, EndOfFile)

	_, err = io.Copy(d, bytes.NewReader(endOfFile))
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
	_, err := io.Copy(file, d)
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
