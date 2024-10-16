package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var serviceName = "SmartabaseSyncService"
var appIcon = "./images/app_icon.png"
var serviceRunningIcon = "./images/running.png"
var serviceStoppedIcon = "./images/not_running.png"
var serviceNotInstalledIcon = "./images/not_installed.png"
var lastEventCount = "50"

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func runningSystrayMenu(myWindow fyne.Window) *fyne.Menu {
	return fyne.NewMenu("MyApp",
		fyne.NewMenuItem(fmt.Sprintf("%s running", serviceName), func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show Dashboard", func() { myWindow.Show() }),
		fyne.NewMenuItem("Stop service", func() { stopService(serviceName) }),
	)
}
func stoppedSystrayMenu(myWindow fyne.Window) *fyne.Menu {
	return fyne.NewMenu("MyApp",
		fyne.NewMenuItem(fmt.Sprintf("%s stopped", serviceName), func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show Dashboard", func() { myWindow.Show() }),
		fyne.NewMenuItem("Start service", func() { startService(serviceName) }),
	)
}
func notinstalledSystrayMenu(myWindow fyne.Window) *fyne.Menu {
	return fyne.NewMenu("MyApp",
		fyne.NewMenuItem(fmt.Sprintf("%s not installed", serviceName), func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show Dashboard", func() { myWindow.Show() }),
	)
}

func setupSystemTray(myApp fyne.App, myWindow fyne.Window) error {
	running_icon, _ := fyne.LoadResourceFromPath(serviceRunningIcon)
	stopped_icon, _ := fyne.LoadResourceFromPath(serviceStoppedIcon)
	notinstalled_icon, _ := fyne.LoadResourceFromPath(serviceNotInstalledIcon)

	if desk, ok := myApp.(desktop.App); ok {
		// Set up initial icons and menu
		m := runningSystrayMenu(myWindow)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(running_icon)

		// Start the status checker
		go func() {
			for {
				status := isServiceRunning(serviceName)
				if status == 0 {
					m := runningSystrayMenu(myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(running_icon)
					myApp.SetIcon(running_icon)
				} else if status == 1 {
					m := stoppedSystrayMenu(myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(stopped_icon)
					myApp.SetIcon(stopped_icon)
				} else {
					m := notinstalledSystrayMenu(myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(notinstalled_icon)
					myApp.SetIcon(notinstalled_icon)
				}
				time.Sleep(10 * time.Second)
			}
		}()

	}
	return nil
}

func main() {

	// Setup the main dashboard window
	myApp, myWindow, err := setupDashboard()
	if err != nil {
		log.Fatalf("Failed to set up dashboard: %v", err)
	}

	// Setup the system tray icon and menu
	if err := setupSystemTray(myApp, myWindow); err != nil {
		log.Fatalf("Failed to set up system tray: %v", err)
	}

	// Display the main window and start the application event loop
	myWindow.ShowAndRun()
}
