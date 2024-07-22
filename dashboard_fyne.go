package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var reportData = []string{}
var reportList *widget.List

func layoutDashboard() *fyne.Container {

	// Crea una list widget per visualizzare i dati del report
	reportList = widget.NewList(
		func() int {
			return len(reportData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(reportData[i])
		},
	)

	content := container.NewStack(
		reportList,
	)
	return content
}

func setupDashboard() (fyne.App, fyne.Window, error) {
	// Initialize the application
	myApp := app.New()

	// Set the application icon
	icon, err := fyne.LoadResourceFromPath(appIcon)
	if err != nil {
		log.Printf("Failed to load application icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		myApp.SetIcon(icon)
	}

	myWindow := myApp.NewWindow(fmt.Sprintf("%s dashboard", serviceName))
	myWindow.Resize(fyne.NewSize(600, 400))

	myWindow.SetContent(layoutDashboard())

	updateList := func() {
		// Prepare the PowerShell command
		cmdStr := fmt.Sprintf("Get-EventLog -LogName Application -Source %s -Newest 10", serviceName)
		cmd := exec.Command("powershell", "-Command", cmdStr)

		// Execute the command and capture the output
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}

		// Convert output to string and split into lines if necessary
		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")

		// Append each line to reportData
		for _, line := range lines {
			if line != "" { // Avoid appending empty lines
				reportData = append(reportData, line)
			}
		}
		// reportData = append(reportData, fmt.Sprintf("Item %d", len(reportData)+1))
		reportList.Refresh()
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			myWindow.Content().Refresh()
			updateList()
		}
	}()

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide() // Hide the window instead of closing it
	})

	// Load and set the running icon for the window
	runningIcon, err := fyne.LoadResourceFromPath(serviceRunningIcon)
	if err != nil {
		log.Printf("Failed to load running icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		myWindow.SetIcon(runningIcon)
	}

	return myApp, myWindow, nil
}
