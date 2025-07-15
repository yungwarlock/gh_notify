package main

import (
	"os"
	"time"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/gen2brain/beeep"
)

var poller GithubPoller

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

	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mNotify := systray.AddMenuItem("Send Notification", "Send a notification")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mNotifications := systray.AddMenuItem("Notifications", "View notifications")
	mTogglePolling := systray.AddMenuItemCheckbox("Toggle Polling", "Toggle notification polling", true)

	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)

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
						beeep.Alert(notification.Reason, notification.Subject.Title, icon.Data)
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
				poller.checkIfNotificationsAsChanged()
			case <-mTogglePolling.ClickedCh:
				shouldPoll = !shouldPoll
				if mTogglePolling.Checked() {
					mTogglePolling.Uncheck()
				} else {
					mTogglePolling.Check()
				}
			case <-mNotify.ClickedCh:
				beeep.Notify("Awesome App", "Notification sent!", "assets/information.png")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// clean up here
}
