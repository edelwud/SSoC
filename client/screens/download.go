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
	downloadLogger   = logrus.WithField("context", "download")
)

func downloadCb(client client.Client) func() {
	return func() {
		name := downloadFilename.Text

		cmd := command.CreateDownloadCommand(name)
		err := client.Exec(cmd)
		if err != nil {
			downloadLogger.Fatalf("cannot execute download command: %s", err)
		}

		file := client.GetContext().FindDownload(name)
		if file == nil {
			downloadLogger.Fatalf("cannot find downloaded file")
		}

		duration, filename, bitrate := file.Duration(), file.Filename, file.Bitrate()
		downloadLogger.WithFields(logrus.Fields{
			"duration (ms)":  duration / 1000000,
			"bitrate (MB/s)": bitrate,
		}).Infof("file \"%s\" successfully downloaded", filename)
	}
}

func CreateDownloadTab(client client.Client) *container.TabItem {
	c := container.NewTabItem("DOWNLOAD", container.NewVBox(
		downloadFilename,
		widget.NewButton("Download file", downloadCb(client)),
	))
	return c
}
