package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
	"main/components/token"
	"main/screens"
	"os"
)

// ConfigFilename path to config file
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

	accessToken, err := LoadSession()
	if err != nil {
		macToken, err := token.GenerateMACToken()
		if err != nil {
			topLevelLogger.Fatalf("cannot generate MAC access token: %q", err)
		}

		accessToken, err = macToken.Row()
		if err != nil {
			topLevelLogger.Fatalf("cannot stringify MAC access token: %q", err)
		}

		err = StoreSession(accessToken)
		if err != nil {
			topLevelLogger.Fatalf("cannot store MAC access token: %q", err)
		}
	}

	tcpClient := InitializeTCPClient(config, accessToken)
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
