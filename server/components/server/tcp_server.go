package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type TcpServer struct {
	Host     string
	Port     string
	Listener net.Listener
	Session  SessionStorage
}

func (s *TcpServer) Run() error {
	var err error
	s.Listener, err = net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}

	err = s.AcceptLoop()
	if err != nil {
		return err
	}

	return nil
}

func (s *TcpServer) AcceptLoop() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		go s.HandleConnection(conn)
	}
}

func (s *TcpServer) HandleConnection(conn net.Conn) {
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		_, err = conn.Write([]byte(myTime))
		if err != nil {
			fmt.Println("kekw")
		}
	}
}

func (s *TcpServer) Close() error {
	if s.Listener == nil {
		return nil
	}

	err := s.Listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreateTCPServer(options Options) Server {
	return &TcpServer{Host: options.Host, Port: options.Port}
}
