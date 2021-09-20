package main

import (
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
	"os"

	g "github.com/AllenDang/giu"
)

const ConfigFilename = "config.yaml"

var (
	topLevelLogger = logrus.WithField("context", "main")
	tcpClient      client.Client
	echoText       string
)

func loop() {
	g.SingleWindow().Layout(
		g.PrepareMsgbox(),
		g.Row(
			g.Button("Close connection").Size(g.Auto, 30).OnClick(func() {
				cmd := command.CreateCloseCommand()
				err := tcpClient.Exec(cmd)
				if err != nil {
					topLevelLogger.Fatalf("cannot disconnect from tcp server: %q", err)
					g.Msgbox("Error", "Cannot disconnect from tcp server")
				}
				os.Exit(0)
			}),
		),
		g.TabBar().TabItems(
			g.TabItem("ECHO").Layout(
				g.Label("Text"),
				g.Row(
					g.InputTextMultiline(&echoText),
					g.Button("ECHO").OnClick(func() {
						cmd := command.CreateEchoCommand(echoText)
						err := tcpClient.Exec(cmd)
						if err != nil {
							topLevelLogger.Fatalf("cannot to exec echo command: %s", err)
							g.Msgbox("Error", "While execution ECHO command: "+err.Error())
						}
					}),
				),
			),
			g.TabItem("TIME").Layout(
				g.Button("TIME").Size(g.Auto, g.Auto).OnClick(func() {
					cmd := command.CreateTimeCommand()
					err := tcpClient.Exec(cmd)
					if err != nil {
						topLevelLogger.Fatalf("cannot to exec time command: %s", err)
						g.Msgbox("Error", "While execution TIME command: "+err.Error())
					}
				}),
			),
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

	wnd := g.NewMasterWindow("TCP client", 400, 250, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
