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

// var serviceRunningIcon = getIcon("./images/icon1.ico")
// var serviceStoppedIcon = getIcon("./images/icon1.ico")
// var serviceNotInstalledIcon = getIcon("./images/icon1.ico")

var appIcon = "./images/app_icon.png"
var serviceRunningIcon = "./images/running.png"
var serviceStoppedIcon = "./images/not_running.png"
var serviceNotInstalledIcon = "./images/not_installed.png"

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// func onReady() {
// 	// myWindow = setupDashboard()
// 	// Set up initial icons and menu
// 	systray.SetIcon(serviceRunningIcon)
// 	systray.SetTooltip(fmt.Sprintf("%s Running", serviceName))
// 	mDashboard := systray.AddMenuItem("Show Dashboard", "Show synchronizer dashboard")
// 	mStart := systray.AddMenuItem("Start Service", "Start the service")
// 	mStop := systray.AddMenuItem("Stop Service", "Stop the service")

// 	// Start the status checker
// 	go func() {
// 		for {
// 			status := isServiceRunning(serviceName)
// 			if status == 0 {
// 				systray.SetIcon(serviceRunningIcon)
// 				systray.SetTooltip(fmt.Sprintf("%s running", serviceName))

// 				mStart.Hide()
// 				mStop.Enable()
// 			} else if status == 1 {
// 				systray.SetIcon(serviceStoppedIcon)
// 				systray.SetTooltip(fmt.Sprintf("%s stopped", serviceName))

// 				mStart.Enable()
// 				mStop.Hide()
// 			} else {
// 				systray.SetIcon(serviceNotInstalledIcon)
// 				systray.SetTooltip(fmt.Sprintf("%s not installed", serviceName))

// 				mStart.Hide()
// 				mStop.Hide()
// 			}
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()

// 	// Menu item actions
// 	go func() {
// 		for {
// 			select {
// 			case <-mDashboard.ClickedCh:
// 				// showDashboard()
// 				mainQueue <- func() {
// 					myWindow.Show()
// 					// myWindow.ShowAndRun()
// 				}
// 			case <-mStart.ClickedCh:
// 				startService(serviceName)
// 			case <-mStop.ClickedCh:
// 				stopService(serviceName)
// 			}
// 		}
// 	}()
// }

// func onExit() {
// 	// Clean up here if needed
// }

func runningSystrayMenu(myApp fyne.App, myWindow fyne.Window) *fyne.Menu {
	return fyne.NewMenu("MyApp",
		fyne.NewMenuItem(fmt.Sprintf("%s running", serviceName), func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show Dashboard", func() { myWindow.Show() }),
		fyne.NewMenuItem("Stop service", func() { stopService(serviceName) }),
	)
}
func stoppedSystrayMenu(myApp fyne.App, myWindow fyne.Window) *fyne.Menu {
	return fyne.NewMenu("MyApp",
		fyne.NewMenuItem(fmt.Sprintf("%s stopped", serviceName), func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show Dashboard", func() { myWindow.Show() }),
		fyne.NewMenuItem("Start service", func() { startService(serviceName) }),
	)
}
func notinstalledSystrayMenu(myApp fyne.App, myWindow fyne.Window) *fyne.Menu {
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
		m := runningSystrayMenu(myApp, myWindow)
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(running_icon)

		// Start the status checker
		go func() {
			for {
				status := isServiceRunning(serviceName)
				if status == 0 {
					m := runningSystrayMenu(myApp, myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(running_icon)
					myApp.SetIcon(running_icon)
				} else if status == 1 {
					m := stoppedSystrayMenu(myApp, myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(stopped_icon)
					myApp.SetIcon(stopped_icon)
				} else {
					m := notinstalledSystrayMenu(myApp, myWindow)
					desk.SetSystemTrayMenu(m)
					desk.SetSystemTrayIcon(notinstalled_icon)
					myApp.SetIcon(notinstalled_icon)
				}
				time.Sleep(30 * time.Second)
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
