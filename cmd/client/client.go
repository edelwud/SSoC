package main

import (
	"SSoC/internal/client"
	tcp "SSoC/internal/client/tcp_client"
	udp "SSoC/internal/client/udp_client"
	"SSoC/internal/options"
)

// InitializeClient initializes client from config
func InitializeClient(config options.Options, accessToken string) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	var c client.Client
	if config.Protocol == "tcp" {
		c = tcp.CreateTCPClient(config, accessToken)
	} else {
		c = udp.CreateUDPClient(config, accessToken)
	}

	err := c.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return c
}
