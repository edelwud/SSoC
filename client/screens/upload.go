package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
)

var (
	uploadLogger = logrus.WithField("context", "upload")
)

func getUploadCb(w fyne.Window, client client.Client) func(reader fyne.URIReadCloser, err error) {
	return func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if reader == nil {
			return
		}

		name, path := reader.URI().Name(), reader.URI().Path()
		cmd := command.CreateUploadCommand(name, path)

		err = client.Exec(cmd)
		if err != nil {
			uploadLogger.Fatalf("cannot execute upload command: %s", err)
		}

		file := client.GetContext().FindUpload(name)
		if file == nil {
			uploadLogger.Fatalf("cannot find uploaded file")
		}

		duration, filename, bitrate := file.Duration(), file.Filename, file.Bitrate()
		uploadLogger.WithFields(logrus.Fields{
			"duration (ms)":  duration / 1000000,
			"bitrate (MB/s)": bitrate,
		}).Infof("file \"%s\" successfully uploaded", filename)
	}
}

func CreateUploadTab(w fyne.Window, client client.Client) *container.TabItem {
	c := container.NewTabItem("UPLOAD", container.NewVBox(
		widget.NewButton("Upload file", func() {
			dialog.ShowFileOpen(getUploadCb(w, client), w)
		}),
	))
	return c
}
