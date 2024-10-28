package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2/app"
	"github.com/getlantern/systray"
)

//go:embed images/*.ico
var icons embed.FS

var serviceName = "smartcoach-sync"
var appIcon = "images/app_icon.ico"
var serviceRunningIcon = "images/running.ico"
var serviceStoppedIcon = "images/not_running.ico"
var serviceNotInstalledIcon = "images/not_installed.ico"
var lastEventCount = "50"

var myApp = app.New()

func onReady() {

	log.Println("onReady called")

	runningIcon, err := icons.ReadFile(serviceRunningIcon)
	if err != nil {
		log.Fatal(err)
	}
	stoppedIcon, err := icons.ReadFile(serviceStoppedIcon)
	if err != nil {
		log.Fatal(err)
	}
	notInstalledIcon, err := icons.ReadFile(serviceNotInstalledIcon)
	if err != nil {
		log.Fatal(err)
	}

	title := fmt.Sprintf("%s event viewer", serviceName)
	runningStatus := fmt.Sprintf("%s running", serviceName)
	stoppedStatus := fmt.Sprintf("%s stopped", serviceName)
	notInstalledStatus := fmt.Sprintf("%s not installed", serviceName)

	systray.SetTitle(title)
	systray.SetTooltip(title)
	systray.SetIcon(runningIcon)

	// Add menu items
	mServiceStatus := systray.AddMenuItem(runningStatus, "Service status")
	systray.AddSeparator()
	mDashboard := systray.AddMenuItem("Show dashboard", "Open the dashboard")
	mStartService := systray.AddMenuItem("Start service", "Start the service")
	mStopService := systray.AddMenuItem("Stop service", "Stop the service")

	mStartService.Hide()

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mDashboard.ClickedCh:
				openDashboard()
			case <-mStartService.ClickedCh:
				startService(serviceName)
			case <-mStopService.ClickedCh:
				stopService(serviceName)
			case <-mServiceStatus.ClickedCh:

			}

		}
	}()

	// Monitor service status
	go func() {
		var previousStatus int
		for {
			status := isServiceRunning(serviceName)
			if status != previousStatus {
				previousStatus = status
				switch status {
				case 0: // Running
					systray.SetIcon(runningIcon)
					mServiceStatus.SetTitle(runningStatus)
					mStartService.Hide()
					mStopService.Show()
				case 1: // Stopped
					systray.SetIcon(stoppedIcon)
					mServiceStatus.SetTitle(stoppedStatus)
					mStartService.Show()
					mStopService.Hide()
				default: // Not installed
					systray.SetIcon(notInstalledIcon)
					mServiceStatus.SetTitle(notInstalledStatus)
					mStartService.Hide()
					mStopService.Hide()
					mDashboard.Hide()
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func onExit() {
	log.Println("Chiusura dell'applicazione...")
}

func main() {
	log.Println("Application starting")
	go func() {
		systray.Run(onReady, onExit)
	}()
	myApp.Run()
	log.Println("Application running")
}
