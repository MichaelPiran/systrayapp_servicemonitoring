package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type EventLog struct {
	TimeGenerated string `json:"TimeGenerated"`
	Message       string `json:"Message"`
}

var eventLogs []EventLog
var stopUpdating bool

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

func openDashboard() {
	stopUpdating = false
	ctx, cancel := context.WithCancel(context.Background())

	// Create a new window
	myWindow := myApp.NewWindow(fmt.Sprintf("%s event viewer", serviceName))
	myWindow.Resize(fyne.NewSize(1000, 400))

	myWindow.SetContent(layoutDashboard())

	// Set the content of the window
	updateList := func() {
		if stopUpdating {
			return
		}

		// Define the PowerShell command
		cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "Get-EventLog", "-LogName", "Application", "-Source", serviceName, "-Newest", lastEventCount, "| Select-Object @{Name='TimeGenerated';Expression={$_.TimeGenerated.ToString('yyyy-MM-dd HH:mm:ss')}}, Message | ConvertTo-Json")

		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateList()
			case <-ctx.Done():
				return
			}
		}
	}()

	myWindow.SetCloseIntercept(func() {
		stopUpdating = true
		cancel()        // Cancel the context to stop the goroutine
		myWindow.Hide() // Hide the window instead of closing it
	})

	// Load and set the running icon for the window
	appIcon, err := icons.ReadFile(serviceAppIcon)
	if err != nil {
		log.Printf("Failed to load running icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		myWindow.SetIcon(fyne.NewStaticResource("runningIcon", appIcon))
	}

	// Show the window
	myWindow.Show()
}
