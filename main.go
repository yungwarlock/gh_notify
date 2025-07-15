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
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mNotify := systray.AddMenuItem("Send Notification", "Send a notification")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mNotifications := systray.AddMenuItem("Notifications", "View notifications")

	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)

	go func() {
		for {
			hasChanged, err := poller.checkIfNotificationsAsChanged()
			if err != nil {
				panic(err)
			}

			if hasChanged {
				poller.getNotifications()
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-mNotifications.ClickedCh:
				poller.checkIfNotificationsAsChanged()
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
