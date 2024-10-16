package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EventLog struct {
	TimeGenerated string `json:"TimeGenerated"`
	Message       string `json:"Message"`
}

var eventLogs []EventLog

func layoutDashboard() *fyne.Container {
	// Create a table widget to display the event logs
	table := widget.NewTable(
		func() (int, int) { return len(eventLogs) + 1, 2 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			label := co.(*widget.Label)
			if tci.Row == 0 {
				if tci.Col == 0 {
					label.SetText("Time")
				} else {
					label.SetText("Message")
				}
			} else {
				if tci.Col == 0 {
					label.SetText(eventLogs[tci.Row-1].TimeGenerated)
				} else {
					label.SetText(eventLogs[tci.Row-1].Message)
				}
			}
		},
	)

	table.SetColumnWidth(0, 300)
	table.SetColumnWidth(1, 500)

	content := container.NewStack(table)
	return content
}

func setupDashboard() (fyne.App, fyne.Window, error) {
	// Initialize the application
	myApp := app.New()

	myApp.Settings().SetTheme(theme.DarkTheme())
	// Set the application icon
	icon, err := fyne.LoadResourceFromPath(serviceRunningIcon)
	if err != nil {
		log.Printf("Failed to load application icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		myApp.SetIcon(icon)
	}

	myWindow := myApp.NewWindow(fmt.Sprintf("%s event viewer", serviceName))
	myWindow.Resize(fyne.NewSize(1000, 400))

	myWindow.SetContent(layoutDashboard())

	updateList := func() {
		// Define the PowerShell command
		cmd := exec.Command("powershell", "-Command", "Get-EventLog", "-LogName", "Application", "-Source", serviceName, "-Newest", lastEventCount, "| Select-Object @{Name='TimeGenerated';Expression={$_.TimeGenerated.ToString('yyyy-MM-dd HH:mm:ss')}}, Message | ConvertTo-Json")
		// cmd := exec.Command("powershell", "-Command", `
		// 	Get-EventLog -LogName Application -Source  -Newest 10 |
		// 	Select-Object @{Name='TimeGenerated';Expression={$_.TimeGenerated.ToString('yyyy-MM-dd HH:mm:ss')}}, Message |
		// 	ConvertTo-Json
		// `)

		// Capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
			return
		}

		// Parse the JSON output
		err = json.Unmarshal(output, &eventLogs)
		if err != nil {
			fmt.Printf("Error parsing JSON: %s\n", err)
			return
		}

		// Refresh the table to display the updated data
		myWindow.Content().Refresh()
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			updateList()
		}
	}()

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide() // Hide the window instead of closing it
	})

	// Load and set the running icon for the window
	runningIcon, err := fyne.LoadResourceFromPath(appIcon)
	if err != nil {
		log.Printf("Failed to load running icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		myWindow.SetIcon(runningIcon)
	}

	return myApp, myWindow, nil
}
