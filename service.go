package main

import (
	"log"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// func isServiceRunning(name string) int {
// 	processes, err := process.Processes()
// 	if err != nil {
// 		log.Println("Error fetching processes:", err)
// 		return 2
// 	}
// 	for _, p := range processes {
// 		pName, err := p.Name()
// 		log.Println(pName)
// 		if err == nil && pName == name {
// 			return 0
// 		}
// 	}
// 	return 1
// }

func isServiceRunning(name string) int {
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
	cmd := exec.Command("sc", "start", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println("Error starting service:", err)
	} else {
		log.Println("Service started:", name)
	}
}

func stopService(name string) {
	cmd := exec.Command("sc", "stop", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println("Error stopping service:", err)
	} else {
		log.Println("Service stopped:", name)
	}
}

func installService(name string) {
	cmd := exec.Command("sc", "install", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println("Error installing service:", err)
	} else {
		log.Println("Service installed:", name)
	}
}
