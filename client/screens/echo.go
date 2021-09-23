package screens

import (
	"bufio"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"main/components/client"
	"main/components/command"
)

var (
	echoEntry  = widget.NewMultiLineEntry()
	echoResult = widget.NewLabel("Result: ")
	echoLogger = logrus.WithField("context", "echo")
)

func getEchoCb(client client.Client) func() {
	return func() {
		cmd := command.CreateEchoCommand(echoEntry.Text)
		err := client.Exec(cmd)
		if err != nil {
			echoLogger.Fatalf("cannot to exec echo command: %s", err)
		}

		conn := client.GetContext().GetConn()
		echo, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			echoLogger.Fatalf("cannot to read from session context: %s", err)
		}

		echoLogger.Info("command executed")
		echoResult.SetText("Result: " + echo)
	}
}

func CreateEchoTab(client client.Client) *container.TabItem {
	c := container.NewTabItem("ECHO", container.NewVBox(
		echoEntry,
		widget.NewButton("Send", getEchoCb(client)),
		echoResult,
	))
	return c
}
