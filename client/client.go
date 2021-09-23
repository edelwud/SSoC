package main

import "main/components/client"

func InitializeTcpClient(config client.Options) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	tcpClient = client.CreateTcpClient(config)

	err := tcpClient.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return tcpClient
}
