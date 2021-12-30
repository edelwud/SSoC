package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type UDPDatachannelClient struct {
	Conn    *net.UDPConn
	Port    string
	Options options.Options
}

func (d UDPDatachannelClient) Read(p []byte) (n int, err error) {
	n, _, err = d.Conn.ReadFromUDP(p)
	if binary.BigEndian.Uint64(p[n-8:n]) == EndOfFile {
		return 0, io.EOF
	}
	return n, err
}

func (d *UDPDatachannelClient) Connect() error {
	UDPAddr, err := net.ResolveUDPAddr("udp", d.Options.Host+":"+d.Port)
	if err != nil {
		return err
	}

	d.Conn, err = net.DialUDP("udp", nil, UDPAddr)
	if err != nil {
		return err
	}

	_, err = d.Conn.Write([]byte(HelloMessage + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (d UDPDatachannelClient) Close() error {
	err := d.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d UDPDatachannelClient) Upload(file *session.File) error {
	_, err := io.Copy(d.Conn, file)
	if err != nil {
		return err
	}

	endOfFile := make([]byte, 8)
	binary.BigEndian.PutUint64(endOfFile, EndOfFile)

	_, err = io.Copy(d.Conn, bytes.NewReader(endOfFile))
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

func (d UDPDatachannelClient) Download(file *session.File) error {
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

func (d UDPDatachannelClient) Accept() error {
	return nil
}

func (d UDPDatachannelClient) Listen() error {
	return nil
}

func (d UDPDatachannelClient) GetPort() string {
	return d.Port
}

func NewUDPDatachannelClient(port string, ops options.Options) Datachannel {
	return &UDPDatachannelClient{
		Port:    port,
		Options: ops,
	}
}
