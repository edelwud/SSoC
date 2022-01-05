package server

import (
	"SSoC/internal/command"
	"SSoC/internal/executor"
	"SSoC/internal/options"
	"SSoC/internal/session"
	"errors"
	"io"
	"net"
	"reflect"
	"syscall"
)

// UDPServer implementation of Server interfaces based on UDP protocol
type UDPServer struct {
	Options     options.Options
	Conn        *net.UDPConn
	Context     session.Storage
	ExecService executor.Service
	Clients     map[string]session.Session
	UDPWriter   *UDPWriter
	Epoll       int
	Event       syscall.EpollEvent
	Events      [MaxEvents]syscall.EpollEvent
}

type UDPWriter struct {
	FD       int
	Sockaddr syscall.Sockaddr
}

const (
	MaxEvents         = 128
	CommandBufferSize = 1024 * 1024
	EPOLLET           = 1 << 31
)

// Run resolves server options from Options
// creates net.Listener with UDPv4 background
// executes AcceptLoop
func (s *UDPServer) Run() error {
	addr, err := net.ResolveUDPAddr("udp", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}

	s.Conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	file, err := s.Conn.File()
	if err != nil {
		return err
	}

	fd := int(file.Fd())
	err = syscall.SetNonblock(fd, true)
	if err != nil {
		return err
	}

	s.Epoll, err = syscall.EpollCreate1(0)
	if err != nil {
		return err
	}

	s.Event.Events = syscall.EPOLLIN
	s.Event.Fd = int32(fd)
	err = syscall.EpollCtl(s.Epoll, syscall.EPOLL_CTL_ADD, fd, &s.Event)
	if err != nil {
		return err
	}

	serverLogger.Infof("server started on port %s", s.Options.Port)

	err = s.AcceptLoop()
	if err != nil {
		return err
	}

	return nil
}

func (s UDPServer) Auth(writer io.Writer, addr net.Addr, c string) error {
	cmd, err := command.ParseCommand(c)
	if err != nil {
		return err
	}

	if cmd.Cmd != "TOKEN" {
		return errors.New("cannot authenticate new client")
	}

	sess := session.CreateServerSession(s.Conn, s.Options, cmd.Parameters[0])

	err = s.AddClient(addr.String(), sess)
	if err != nil {
		return err
	}

	err = s.ExecService.Process(writer, sess, cmd)
	if err != nil {
		return err
	}

	return nil
}

// AcceptLoop accepts client connection, sets keep alive and keep alive period options from Options, handles connection
func (s *UDPServer) AcceptLoop() error {
	for {
		n, err := syscall.EpollWait(s.Epoll, s.Events[:], -1)
		if err != nil {
			serverLogger.Fatalf("epoll wait error: %q", err)
			return err
		}

		for i := 0; i < n; i++ {
			fd := int(s.Events[i].Fd)
			buf := make([]byte, CommandBufferSize)
			n, from, err := syscall.Recvfrom(fd, buf, 0)
			if err != nil {
				serverLogger.Warnf("client disconnected: %q", err)
				break
			}

			v := reflect.ValueOf(from).Elem()
			addr := v.FieldByName("Addr").Interface().([16]byte)
			port := v.FieldByName("Port").Interface().(int)

			s.HandleClient(fd, &net.UDPAddr{IP: addr[12:16], Port: port}, from, buf[:n])
			if err != nil {
				return err
			}
		}
	}
}

func (s UDPServer) HandleClient(fd int, addr net.Addr, to syscall.Sockaddr, buf []byte) {
	defer syscall.Close(fd)

	sess, err := s.FindClient(addr.String())
	writer := CreateUDPWriter(fd, to)
	if err != nil {
		err := s.Auth(writer, addr, string(buf))
		if err != nil {
			serverLogger.Warnf("client disconnected")
			return
		}
		return
	}

	cmd, err := command.ParseCommand(string(buf))
	if err != nil {
		serverLogger.Warnf("sent undefined command")
		return
	}

	err = s.ExecService.Process(writer, sess, cmd)
	if err != nil {
		serverLogger.Warnf("command execution error: %s", err)
		return
	}

	serverLogger.Infof("command %s successfully executed", cmd.Cmd)
}

func (s UDPServer) FindClient(address string) (session.Session, error) {
	if _, ok := s.Clients[address]; !ok {
		return nil, errors.New("client not found")
	}

	return s.Clients[address], nil
}

func (s UDPServer) AddClient(address string, session session.Session) error {
	if _, ok := s.Clients[address]; ok {
		return errors.New("client already exists")
	}

	s.Clients[address] = session
	return nil
}

// Close closes net.Listener
func (s *UDPServer) Close() error {
	if s.Conn == nil {
		return nil
	}

	err := syscall.Close(s.Epoll)
	if err != nil {
		return err
	}

	err = s.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (w UDPWriter) Write(p []byte) (n int, err error) {
	err = syscall.Sendmsg(w.FD, p, nil, w.Sockaddr, 0)
	if err != nil {
		return 0, err
	}

	return len(p), err
}

func CreateUDPWriter(fd int, to syscall.Sockaddr) *UDPWriter {
	return &UDPWriter{
		Sockaddr: to,
		FD:       fd,
	}
}

// CreateUDPServer creates UDPServer with initialized session.ServerStorage and executor.ServerExecutorService
func CreateUDPServer(options options.Options) Server {
	ctx := session.CreateServerSessionStorage()
	executorService := executor.RegisterServerExecutorService(ctx)
	return &UDPServer{
		Options:     options,
		Context:     ctx,
		ExecService: executorService,
		Clients:     make(map[string]session.Session),
	}
}
