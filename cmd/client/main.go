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

	cl := InitializeClient(config, accessToken)
	defer func(cl client.Client) {
		err := cl.Disconnect()
		if err != nil {
			return
		}
	}(cl)

	a := app.New()
	w := a.NewWindow("TCP client")

	w.SetContent(container.NewVBox(
		widget.NewButton("Close connection", func() {
			req := requester.CreateCloseRequester()
			err := cl.Exec(req)
			if err != nil {
				topLevelLogger.Fatalf("cannot disconnect from tcp server: %q", err)
			}
			os.Exit(0)
		}),
		container.NewAppTabs(
			screens.CreateEchoTab(cl),
			screens.CreateTimeTab(cl),
			screens.CreateUploadTab(w, cl),
			screens.CreateDownloadTab(cl),
		),
	))
	w.Resize(fyne.Size{Width: 550, Height: 400})
	w.ShowAndRun()
}
