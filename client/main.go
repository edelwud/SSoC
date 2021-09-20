package main

import (
	"bufio"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
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

	a := app.New()
	w := a.NewWindow("TCP client")

	echoEntry := widget.NewMultiLineEntry()
	echoResult := widget.NewLabel("Result: ")
	currentTime := widget.NewLabel("Result: ")

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
			container.NewTabItem("ECHO", container.NewVBox(
				echoEntry,
				widget.NewButton("Send", func() {
					cmd := command.CreateEchoCommand(echoEntry.Text)
					err := tcpClient.Exec(cmd)
					if err != nil {
						topLevelLogger.Fatalf("cannot to exec echo command: %s", err)
					}

					conn := tcpClient.GetContext().GetConn()
					echo, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						topLevelLogger.Fatalf("cannot to read from session context: %s", err)
					}

					echoResult.SetText("Result: " + echo)
				}),
				echoResult,
			)),
			container.NewTabItem("TIME", container.NewVBox(
				widget.NewButton("Current time", func() {
					cmd := command.CreateTimeCommand()
					err := tcpClient.Exec(cmd)
					if err != nil {
						topLevelLogger.Fatalf("cannot to exec time command: %s", err)
					}

					conn := tcpClient.GetContext().GetConn()
					time, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						topLevelLogger.Fatalf("cannot to read from session context: %s", err)
					}

					currentTime.SetText("Result: " + time)
				}),
				currentTime,
			)),
			container.NewTabItem("UPLOAD", container.NewVBox(
				widget.NewButton("Upload file", func() {
					dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
						if err != nil {
							dialog.ShowError(err, w)
							return
						}

						if reader == nil {
							return
						}

						cmd := command.CreateUploadCommand(reader.URI().Name(), reader.URI().Path())
						err = tcpClient.Exec(cmd)

						if err != nil {
							topLevelLogger.Fatalf("cannot execute upload command: %s", err)
						}

						file := tcpClient.GetContext().FindUpload(reader.URI().Name())
						if file == nil {
							topLevelLogger.Fatalf("cannot find uploaded file")
						}

						topLevelLogger.Infof("spend %d seconds for uploading file %s, bitrate %d", file.Duration(), file.Filename, file.Bitrate())
					}, w)
				}),
			)),
		),
	))
	w.Resize(fyne.Size{Width: 550, Height: 400})
	w.ShowAndRun()
}
