package main

import (
	_ "embed"
	"os"
	"time"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
)

var poller GithubPoller

//go:embed github-mark/github-mark.png
var icon []byte

func init() {
	token, exists := os.LookupEnv("GITHUB_TOKEN")

	if !exists {
		ghToken, err := getTokenFromGithubCLI()
		if err != nil {
			panic(err)
		}

		token = ghToken
	}

	poller = *NewGithubPoller(token)
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	shouldPoll := true

	systray.SetIcon(icon)
	systray.SetTitle("GH Notify")
	systray.SetTooltip("Get github notifications on desktop")
	mNotifications := systray.AddMenuItem("Notifications", "View notifications")
	mTogglePolling := systray.AddMenuItemCheckbox("Toggle Polling", "Toggle notification polling", true)
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			if shouldPoll {
				hasChanged, err := poller.checkIfNotificationsAsChanged()
				if err != nil {
					panic(err)
				}

				if hasChanged {
					notifications, err := poller.getNotifications()
					if err != nil {
						panic(err)
					}

					for _, notification := range *notifications {
						beeep.Notify(notification.Reason, notification.Subject.Title, icon)
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-mNotifications.ClickedCh:
				openBrowser("https://github.com/notifications")
			case <-mTogglePolling.ClickedCh:
				shouldPoll = !shouldPoll
				if mTogglePolling.Checked() {
					mTogglePolling.Uncheck()
				} else {
					mTogglePolling.Check()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// clean up here
}
