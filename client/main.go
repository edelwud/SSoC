package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
	"main/screens"
	"os"
)

const ConfigFilename = "config.yaml"

var (
	topLevelLogger = logrus.WithField("context", "main")
	tcpClient      client.Client
)

func main() {
	config, err := LoadClientConfig(ConfigFilename)
	if err != nil {
		topLevelLogger.Fatalf("cannot read config: %s", err)
	}

	tcpClient := InitializeTCPClient(config)
	defer func(tcpClient client.Client) {
		err := tcpClient.Disconnect()
		if err != nil {
			return
		}
	}(tcpClient)

	a := app.New()
	w := a.NewWindow("TCP client")

	w.SetContent(container.NewVBox(
		widget.NewButton("Close connection", func() {
			cmd := command.CreateCloseCommand()
			err := tcpClient.Exec(cmd)
			if err != nil {
				topLevelLogger.Fatalf("cannot disconnect from tcp server: %q", err)
			}
			os.Exit(0)
		}),
		container.NewAppTabs(
			screens.CreateEchoTab(tcpClient),
			screens.CreateTimeTab(tcpClient),
			screens.CreateUploadTab(w, tcpClient),
			screens.CreateDownloadTab(tcpClient),
		),
	))
	w.Resize(fyne.Size{Width: 550, Height: 400})
	w.ShowAndRun()
}
