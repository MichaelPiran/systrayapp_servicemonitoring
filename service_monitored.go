package main

import (
	"log"
	"os/exec"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func isMonitoredServiceRunning(name string) int {
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

	for _, serviceMonitoredName := range services {
		// Open the service
		s, err := m.OpenService(serviceMonitoredName)
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
		if serviceMonitoredName == name {
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

func startServiceMonitored(name string) {
	err := exec.Command("sc", "start", name).Run()
	if err != nil {
		log.Println("Error starting service:", err)
	} else {
		log.Println("Service started:", name)
	}
}

func stopServiceMonitored(name string) {
	err := exec.Command("sc", "stop", name).Run()
	if err != nil {
		log.Println("Error stopping service:", err)
	} else {
		log.Println("Service stopped:", name)
	}
}
