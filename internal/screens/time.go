package screens

import (
	"SSoC/internal/client"
	"SSoC/internal/requester"
	"bufio"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var (
	timeCurrent = widget.NewLabel("Result: ")
	timeLogger  = logrus.WithField("context", "time")
)

func getCurrentTimeCb(client client.Client) func() {
	return func() {
		cmd := requester.CreateTimeRequester()
		err := client.Exec(cmd)
		if err != nil {
			timeLogger.Fatalf("cannot to exec time command: %s", err)
		}

		conn := client.GetContext().GetConn()
		time, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			timeLogger.Fatalf("cannot to read from session context: %s", err)
		}

		timeLogger.Info("command executed")
		timeCurrent.SetText("Result: " + time)
	}
}

// CreateTimeTab creates time tab widget
func CreateTimeTab(client client.Client) *container.TabItem {
	c := container.NewTabItem("TIME", container.NewVBox(
		widget.NewButton("Current time", getCurrentTimeCb(client)),
		timeCurrent,
	))
	return c
}
