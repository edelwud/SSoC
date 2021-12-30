package main

import (
	"SSoC/internal/client"
	"SSoC/internal/options"
)

// InitializeClient initializes client from config
func InitializeClient(config options.Options, accessToken string) client.Client {
	topLevelLogger.Infof("config loaded: %+v", config)

	c := client.New(config.Protocol, config, accessToken)
	err := c.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	return c
}
