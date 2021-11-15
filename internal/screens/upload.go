package screens

import (
	"SSoC/internal/client"
	"SSoC/internal/requester"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var (
	uploadLogger = logrus.WithField("context", "upload")
)

func uploadProgressBar(client client.Client) *widget.ProgressBar {
	progress := widget.NewProgressBar()

	go func() {
		ctx := client.GetContext()
		for {
			status := ctx.UploadStatus()
			if status == 0 {
				progress.Hide()
			}
			if status != 0 {
				progress.SetValue(status)
				progress.Show()
			}
		}
	}()

	return progress
}

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
		cmd := requester.CreateUploadRequester(name, path)

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

// CreateUploadTab creates upload tab widget
func CreateUploadTab(w fyne.Window, client client.Client) *container.TabItem {
	c := container.NewTabItem("UPLOAD", container.NewVBox(
		widget.NewButton("Upload file", func() {
			dialog.ShowFileOpen(getUploadCb(w, client), w)
		}),
		uploadProgressBar(client),
	))
	return c
}
