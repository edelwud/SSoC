package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"

	g "github.com/AllenDang/giu"
)

const ConfigFilename = "../server/config.yaml"

var (
	topLevelLogger = logrus.WithField("context", "main")
	tcpClient      client.Client
)

func onEchoButton() {
	cmd := command.CreateEchoCommand("KEKW")
	err := tcpClient.Exec(cmd)
	if err != nil {
		topLevelLogger.Fatalf("cannot to exec echo command: %s", err)
	}
}

func onTimeButton() {
	fmt.Println("Im sooooooo cute!!")
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("Hello world from giu"),
		g.Row(
			g.Button("Send echo command").OnClick(onEchoButton),
			g.Button("Send time command").OnClick(onTimeButton),
		),
	)
}

func main() {
	config, err := LoadClientConfig(ConfigFilename)
	if err != nil {
		topLevelLogger.Fatalf("cannot read config: %s", err)
	}

	topLevelLogger.Infof("config loaded: %+v", config)

	tcpClient = client.CreateTcpClient(config)

	err = tcpClient.Connect()
	if err != nil {
		topLevelLogger.Fatalf("cannot connect to server: %s", err)
	}

	defer func(tcpClient client.Client) {
		err := tcpClient.Disconnect()
		if err != nil {
			return
		}
	}(tcpClient)

	wnd := g.NewMasterWindow("Hello world", 400, 200, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
