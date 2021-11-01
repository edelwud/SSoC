package main

import (
	"main/components/client"
	"main/components/options"
)

// InitializeTCPClient initializes TCP client from config
func InitializeTCPClient(config options.Options) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	tcpClient = client.CreateTCPClient(config)

	err := tcpClient.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return tcpClient
}
