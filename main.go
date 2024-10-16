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
	runningIcon, err := fyne.LoadResourceFromPath(serviceRunningIcon)
	if err != nil {
		return err
	}
	stoppedIcon, err := fyne.LoadResourceFromPath(serviceStoppedIcon)
	if err != nil {
		return err
	}
	notInstalledIcon, err := fyne.LoadResourceFromPath(serviceNotInstalledIcon)
	if err != nil {
		return err
	}

	if desk, ok := myApp.(desktop.App); ok {
		// Set up initial icons and menu
		updateSystemTray(desk, myApp, myWindow, runningIcon, runningSystrayMenu(myWindow))

		// Start the status checker
		go func() {
			var previousStatus int
			for {
				status := isServiceRunning(serviceName)
				if status != previousStatus {
					log.Printf("Service status changed")
					log.Printf("%d", previousStatus)
					log.Printf("%d", status)
					previousStatus = status
					switch status {
					case 0:
						updateSystemTray(desk, myApp, myWindow, runningIcon, runningSystrayMenu(myWindow))
					case 1:
						updateSystemTray(desk, myApp, myWindow, stoppedIcon, stoppedSystrayMenu(myWindow))
					default:
						updateSystemTray(desk, myApp, myWindow, notInstalledIcon, notinstalledSystrayMenu(myWindow))
					}
				}
				time.Sleep(10 * time.Second)
			}
		}()
	}
	return nil
}

func updateSystemTray(desk desktop.App, myApp fyne.App, myWindow fyne.Window, icon fyne.Resource, menu *fyne.Menu) {
	desk.SetSystemTrayMenu(menu)
	desk.SetSystemTrayIcon(icon)
	myApp.SetIcon(icon)
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
