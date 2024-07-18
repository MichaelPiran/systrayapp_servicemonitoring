package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
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

func isServiceRunning(name string) int {
	processes, err := process.Processes()
	if err != nil {
		log.Println("Error fetching processes:", err)
		return 2
	}
	for _, p := range processes {
		pName, err := p.Name()
		log.Println(pName)
		if err == nil && pName == name {
			return 0
		}
	}
	return 1
}

func isServiceRunning_1(name string) int {
	// Open the Windows Service Control Manager
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()

	// Get a list of all services
	services, err := m.ListServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, serviceName := range services {
		// Open the service
		s, err := m.OpenService(serviceName)
		if err != nil {
			continue
		}
		defer s.Close()

		// Get the service status
		status, err := s.Query()
		if err != nil {
			continue
		}

		// Check if the service is running or stopped
		if serviceName == name {
			switch status.State {
			case svc.Running: // Running
				// fmt.Printf("Service is running")
				return 0
			case svc.Stopped: // Stopped
				// fmt.Printf("Service is stopped")
				return 1
			default: // Not installed
				// fmt.Printf("Service is not installed")
				return 2
			}
		}
	}
	return 2

}

func startService(name string) {
	err := exec.Command("sc", "start", name).Run()
	if err != nil {
		log.Println("Error starting service:", err)
	} else {
		log.Println("Service started:", name)
	}
}

func stopService(name string) {
	err := exec.Command("sc", "stop", name).Run()
	if err != nil {
		log.Println("Error stopping service:", err)
	} else {
		log.Println("Service stopped:", name)
	}
}

func showDashboard() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		// os.Exit(0)
	}()
	// app.Main()

	// // PowerShell command to get the last 10 events from the specified service
	// psCommand := fmt.Sprintf("Get-EventLog -LogName Application -Source %s -Newest 10", serviceName)

	// // Execute the PowerShell command
	// cmd := exec.Command("powershell", "-Command", psCommand)

	// // Capture the output
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Failed to execute PowerShell command: %v", err)
	// }

	// // Print the output
	// fmt.Println(string(output))
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	var eventLogs string

	fetchLogs := func() string {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		psCommand := fmt.Sprintf("Get-EventLog -LogName Application -Source %s -Newest 10", serviceName)
		cmd := exec.Command("powershell", "-Command", psCommand)

		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return fmt.Sprintf("Error fetching logs: %v", err)
		}
		return out.String()
	}

	// Ticker to fetch logs every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Goroutine to fetch logs periodically
	go func() {
		for range ticker.C {
			eventLogs = fetchLogs()
			window.Invalidate()
			log.Printf("fetching logs: %s", eventLogs)
		}
	}()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			title := material.H1(theme, "Hello, Gio")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func onReady() {
	// Set up initial icons and menu
	systray.SetIcon(serviceRunningIcon)
	systray.SetTooltip(fmt.Sprintf("%s Running", serviceName))
	mDashboard := systray.AddMenuItem("Show Dashboard", "Show synchronizer dashboard")
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start Service", "Start the service")
	mStop := systray.AddMenuItem("Stop Service", "Stop the service")

	// Start the status checker
	go func() {
		for {
			status := isServiceRunning_1(serviceName)
			if status == 0 {
				systray.SetIcon(serviceRunningIcon)
				systray.SetTooltip(fmt.Sprintf("%s running", serviceName))

				mStart.Disable()
				mStop.Enable()
			} else if status == 1 {
				systray.SetIcon(serviceStoppedIcon)
				systray.SetTooltip(fmt.Sprintf("%s stopped", serviceName))

				mStart.Enable()
				mStop.Disable()
			} else {
				systray.SetIcon(serviceNotInstalledIcon)
				systray.SetTooltip(fmt.Sprintf("%s not installed", serviceName))

				mStart.Disable()
				mStop.Disable()
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
