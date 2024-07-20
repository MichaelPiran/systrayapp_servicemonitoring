package main

import (
	"fmt"
	"log"

	"os"
	"time"

	"github.com/getlantern/systray"
)

var serviceName = "SmartabaseSyncService"
var serviceRunningIcon = getIcon("./images/icon1.ico")
var serviceStoppedIcon = getIcon("./images/icon1.ico")
var serviceNotInstalledIcon = getIcon("./images/icon1.ico")

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func onReady() {
	// Set up initial icons and menu
	systray.SetIcon(serviceRunningIcon)
	systray.SetTooltip(fmt.Sprintf("%s Running", serviceName))
	mDashboard := systray.AddMenuItem("Show Dashboard", "Show synchronizer dashboard")
	mStart := systray.AddMenuItem("Start Service", "Start the service")
	mStop := systray.AddMenuItem("Stop Service", "Stop the service")

	// Start the status checker
	go func() {
		for {
			status := isServiceRunning(serviceName)
			if status == 0 {
				systray.SetIcon(serviceRunningIcon)
				systray.SetTooltip(fmt.Sprintf("%s running", serviceName))

				mStart.Hide()
				mStop.Enable()
			} else if status == 1 {
				systray.SetIcon(serviceStoppedIcon)
				systray.SetTooltip(fmt.Sprintf("%s stopped", serviceName))

				mStart.Enable()
				mStop.Hide()
			} else {
				systray.SetIcon(serviceNotInstalledIcon)
				systray.SetTooltip(fmt.Sprintf("%s not installed", serviceName))

				mStart.Hide()
				mStop.Hide()
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Menu item actions
	go func() {
		for {
			select {
			case <-mDashboard.ClickedCh:
				showDashboard()
			case <-mStart.ClickedCh:
				startService(serviceName)
			case <-mStop.ClickedCh:
				stopService(serviceName)
			}
		}
	}()
}

func onExit() {
	// Clean up here if needed
}

func main() {
	systray.Run(onReady, onExit)
}
