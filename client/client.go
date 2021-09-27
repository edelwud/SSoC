package main

import "main/components/client"

func InitializeTCPClient(config client.Options) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	tcpClient = client.CreateTCPClient(config)

	err := tcpClient.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return tcpClient
}
