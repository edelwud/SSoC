package main

import (
	"github.com/sirupsen/logrus"
	s "server/components/server"
)

const ConfigFilename = "config.yaml"

var topLevelLogger = logrus.WithField("context", "main")

func main() {
	cfg, err := LoadConfig(ConfigFilename)
	if err != nil {
		topLevelLogger.Fatalf("cannot read config: %s", err)
	}

	topLevelLogger.Infof("config loaded: %+v", cfg)

	server := s.CreateTcpServer(s.Options{Host: cfg.Server.Host, Port: cfg.Server.Port})
	err = server.Run()
	if err != nil {
		topLevelLogger.Fatalf("server running error: %s", err)
	}

	defer func(server s.Server) {
		err := server.Close()
		if err != nil {
			topLevelLogger.Fatalf("closing server socket error: %s", err)
		}
	}(server)
}
