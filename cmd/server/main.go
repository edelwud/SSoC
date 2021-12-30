package main

import (
	c "SSoC/internal/config/server_config"
	"SSoC/internal/server"
	"github.com/sirupsen/logrus"
)

// ConfigFilename path to config file
const ConfigFilename = "config.yaml"

var topLevelLogger = logrus.WithField("context", "main")

func main() {
	cfg, err := c.LoadServerConfig(ConfigFilename)
	if err != nil {
		topLevelLogger.Fatalf("cannot read config: %s", err)
	}

	topLevelLogger.Infof("config loaded: %+v", cfg)

	s := server.New(cfg.Protocol, cfg)
	err = s.Run()
	if err != nil {
		topLevelLogger.Fatalf("server running error: %s", err)
	}

	defer func(s server.Server) {
		err := s.Close()
		if err != nil {
			topLevelLogger.Fatalf("closing server socket error: %s", err)
		}
	}(s)
}
