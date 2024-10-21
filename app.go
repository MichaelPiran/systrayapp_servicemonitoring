package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/getlantern/systray"
// )

// func onReady() {
// 	elog.Info(1, "setup systray")

// 	runningIcon, err := os.ReadFile(serviceRunningIcon)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	stoppedIcon, err := os.ReadFile(serviceStoppedIcon)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	notInstalledIcon, err := os.ReadFile(serviceNotInstalledIcon)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	title := fmt.Sprintf("%s event viewer", serviceMonitoredName)
// 	runningStatus := fmt.Sprintf("%s running", serviceMonitoredName)
// 	stoppedStatus := fmt.Sprintf("%s stopped", serviceMonitoredName)
// 	notInstalledStatus := fmt.Sprintf("%s not installed", serviceMonitoredName)

// 	systray.SetTitle(title)
// 	systray.SetTooltip(title)
// 	systray.SetIcon(runningIcon)

// 	// Add menu items
// 	mServiceStatus := systray.AddMenuItem(runningStatus, "Service status")
// 	systray.AddSeparator()
// 	mDashboard := systray.AddMenuItem("Show dashboard", "Open the dashboard")
// 	mStartService := systray.AddMenuItem("Start service", "Start the service")
// 	mStopService := systray.AddMenuItem("Stop service", "Stop the service")

// 	mStartService.Hide()

// 	// Handle menu item clicks
// 	go func() {
// 		elog.Info(1, "onReady")
// 		for {
// 			select {
// 			case <-mDashboard.ClickedCh:
// 				openDashboard()
// 			case <-mStartService.ClickedCh:
// 				startServiceMonitored(serviceMonitoredName)
// 			case <-mStopService.ClickedCh:
// 				stopServiceMonitored(serviceMonitoredName)
// 			case <-mServiceStatus.ClickedCh:

// 			}

// 		}
// 	}()

// 	// Monitor service status
// 	go func() {
// 		var previousStatus int
// 		for {
// 			status := isMonitoredServiceRunning(serviceMonitoredName)
// 			if status != previousStatus {
// 				previousStatus = status
// 				switch status {
// 				case 0: // Running
// 					systray.SetIcon(runningIcon)
// 					mServiceStatus.SetTitle(runningStatus)
// 					mStartService.Hide()
// 					mStopService.Show()
// 				case 1: // Stopped
// 					systray.SetIcon(stoppedIcon)
// 					mServiceStatus.SetTitle(stoppedStatus)
// 					mStartService.Show()
// 					mStopService.Hide()
// 				default: // Not installed
// 					systray.SetIcon(notInstalledIcon)
// 					mServiceStatus.SetTitle(notInstalledStatus)
// 					mStartService.Hide()
// 					mStopService.Hide()
// 					mDashboard.Hide()
// 				}
// 			}
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }

// func onExit() {
// 	fmt.Println("Chiusura dell'applicazione...")
// }

// func runApp() {
// 	elog.Info(1, "runApp0")
// 	go func() {
// 		systray.Run(onReady, onExit)
// 	}()
// 	myApp.Run()
// }
