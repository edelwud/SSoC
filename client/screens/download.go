package screens

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
)

var (
	downloadFilename = widget.NewEntry()
	downloadLogger   = logrus.WithField("context", "time")
)

func downloadCb(client client.Client) func() {
	return func() {
		cmd := command.CreateDownloadCommand(downloadFilename.Text)
		err := client.Exec(cmd)
		if err != nil {
			downloadLogger.Fatalf("cannot execute download command: %s", err)
		}
	}
}

func CreateDownloadTab(client client.Client) *container.TabItem {
	c := container.NewTabItem("DOWNLOAD", container.NewVBox(
		downloadFilename,
		widget.NewButton("Download file", downloadCb(client)),
	))
	return c
}
