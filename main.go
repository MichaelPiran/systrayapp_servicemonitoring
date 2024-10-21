package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

//go:embed images/*.ico
var icons embed.FS

var svcName = "SmartCoachCloudSyncEventViewer"

const serviceMonitoredName = "SmartabaseSyncService"

const appIcon = "./images/app_icon.ico"
const serviceRunningIcon = "./images/running.ico"
const serviceStoppedIcon = "./images/not_running.ico"
const serviceNotInstalledIcon = "./images/not_installed.ico"
const lastEventCount = "50"

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {

	// set variable svcName from command line arguments
	flag.StringVar(&svcName, "name", svcName, "name of the service")
	flag.Parse()

	// check if the application is running as a service or not
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		runService(svcName, false)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		runService(svcName, true)
		return
	case "install":
		err = installService(svcName, mgr.Config{
			DisplayName: "Smartcoach-Sync-Event-Viewer",
			Description: "Smartcoach Cloud Sync Service Event Viewer",
			StartType:   mgr.StartAutomatic,
		})
		if err == nil {
			log.Printf("Service %s installed!", svcName)
		}
	case "remove":
		err = removeService(svcName)
		if err == nil {
			log.Printf("Service %s removed!", svcName)
		}
	case "start":
		err = startService(svcName)
		if err == nil {
			log.Printf("Service %s started!", svcName)
		}
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
		if err == nil {
			log.Printf("Service %s stopped!", svcName)
		}
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return
}
