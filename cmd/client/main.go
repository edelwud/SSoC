package main

import (
	"SSoC/internal/client"
	c "SSoC/internal/config/client_config"
	"SSoC/internal/requester"
	"SSoC/internal/screens"
	"SSoC/internal/token"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"os"
)

// ConfigFilename path to config file
const ConfigFilename = "config.yaml"

var (
	topLevelLogger = logrus.WithField("context", "main")
)

func main() {
	config, err := c.LoadClientConfig(ConfigFilename)
	if err != nil {
		topLevelLogger.Fatalf("cannot read config: %s", err)
	}

	accessToken, err := c.LoadSession()
	if err != nil {
		macToken, err := token.GenerateMACToken()
		if err != nil {
			topLevelLogger.Fatalf("cannot generate MAC access token: %q", err)
		}

		accessToken, err = macToken.Row()
		if err != nil {
			topLevelLogger.Fatalf("cannot stringify MAC access token: %q", err)
		}

		err = c.StoreSession(accessToken)
		if err != nil {
			topLevelLogger.Fatalf("cannot store MAC access token: %q", err)
		}
	}

	clientChannel := InitializeTCPClient(config, accessToken)
	defer func(tcpClient client.Client) {
		err := tcpClient.Disconnect()
		if err != nil {
			return
		}
	}(clientChannel)

	a := app.New()
	w := a.NewWindow("TCP client")

	w.SetContent(container.NewVBox(
		widget.NewButton("Close connection", func() {
			req := requester.CreateCloseRequester()
			err := clientChannel.Exec(req)
			if err != nil {
				topLevelLogger.Fatalf("cannot disconnect from tcp server: %q", err)
			}
			os.Exit(0)
		}),
		container.NewAppTabs(
			screens.CreateEchoTab(clientChannel),
			screens.CreateTimeTab(clientChannel),
			screens.CreateUploadTab(w, clientChannel),
			screens.CreateDownloadTab(clientChannel),
		),
	))
	w.Resize(fyne.Size{Width: 550, Height: 400})
	w.ShowAndRun()
}
