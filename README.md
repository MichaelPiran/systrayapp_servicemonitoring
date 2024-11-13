### Service-monitoring systray app
This is a simple system tray application for monitoring services on Windows. It allows you to start, stop, and check the status of specified services directly from the system tray.

### Features
- Monitor log service
- Start and stop service
- View service status

### Prerequisites
- Go programming language installed
- Windows operating system


### Build
```sh
go build -ldflags="-H windowsgui" -o SERVICENAME.exe .
```

