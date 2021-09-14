package main

import (
	s "server/components/server"
)

func main() {
	cfg := LoadConfig()

	server := s.CreateTcpServer(s.Options{Host: cfg.Server.Host, Port: cfg.Server.Port})
	err := server.Run()
	if err != nil {
		panic(err)
	}

	defer server.Close()
}
