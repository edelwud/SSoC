package main

import (
	"SSoC/internal/client"
	tcp "SSoC/internal/client/tcp_client"
	"SSoC/internal/options"
)

// InitializeTCPClient initializes TCP client from config
func InitializeTCPClient(config options.Options, accessToken string) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	c := tcp.CreateTCPClient(config, accessToken)

	err := c.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return c
}