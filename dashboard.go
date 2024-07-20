package main

func showDashboard_gioui() {}

// import (
// 	"bytes"
// 	"fmt"
// 	"image/color"
// 	"log"
// 	"os/exec"
// 	"time"

// 	"gioui.org/app"
// 	"gioui.org/layout"
// 	"gioui.org/op"
// 	"gioui.org/widget"
// 	"gioui.org/widget/material"
// )

// var (
// 	button1 widget.Clickable
// 	button2 widget.Clickable
// 	logs    = []string{"Log 1: Initialization complete.", "Log 2: Connection established.", "Log 3: Data received."}
// )

// func showDashboard() {
// 	go func() {
// 		window := new(app.Window)
// 		window.Option(app.Title("Synchronizer Dashboard"))
// 		setWindowIcon(window)
// 		// window.Option(app.Icon("images/icon1.ico"))
// 		err := run(window)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// os.Exit(0)
// 	}()
// 	// app.Main()

// 	// // PowerShell command to get the last 10 events from the specified service
// 	// psCommand := fmt.Sprintf("Get-EventLog -LogName Application -Source %s -Newest 10", serviceName)

// 	// // Execute the PowerShell command
// 	// cmd := exec.Command("powershell", "-Command", psCommand)

// 	// // Capture the output
// 	// output, err := cmd.CombinedOutput()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to execute PowerShell command: %v", err)
// 	// }

// 	// // Print the output
// 	// fmt.Println(string(output))
// }

// func run(window *app.Window) error {
// 	theme := material.NewTheme()
// 	var ops op.Ops
// 	var eventLogs string

// 	fetchLogs := func() string {
// 		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		// defer cancel()

// 		psCommand := fmt.Sprintf("Get-EventLog -LogName Application -Source %s -Newest 10", serviceName)
// 		cmd := exec.Command("powershell", "-Command", psCommand)

// 		var out bytes.Buffer
// 		cmd.Stdout = &out
// 		if err := cmd.Run(); err != nil {
// 			return fmt.Sprintf("Error fetching logs: %v", err)
// 		}
// 		return out.String()
// 	}

// 	// Ticker to fetch logs every 30 seconds
// 	ticker := time.NewTicker(30 * time.Second)
// 	defer ticker.Stop()

// 	// Goroutine to fetch logs periodically
// 	go func() {
// 		for range ticker.C {
// 			eventLogs = fetchLogs()
// 			window.Invalidate()
// 			log.Printf("fetching logs: %s", eventLogs)
// 		}
// 	}()

// 	for {
// 		switch e := window.Event().(type) {
// 		case app.DestroyEvent:
// 			return e.Err
// 		case app.FrameEvent:
// 			gtx := app.NewContext(&ops, e)
// 			layoutUI(gtx, theme)
// 			e.Frame(gtx.Ops)
// 		}
// 	}
// }

// func layoutUI(gtx layout.Context, th *material.Theme) {
// 	layout.Flex{
// 		Axis: layout.Vertical,
// 	}.Layout(gtx,
// 		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
// 			return layoutButtons(gtx, th)
// 		}),
// 		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
// 			return layoutLogs(gtx, th)
// 		}),
// 	)
// }

// func layoutButtons(gtx layout.Context, th *material.Theme) layout.Dimensions {
// 	return layout.Flex{
// 		Axis:      layout.Horizontal,
// 		Spacing:   layout.SpaceBetween,
// 		Alignment: layout.Middle,
// 	}.Layout(gtx,
// 		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
// 			for button1.Clicked(gtx) {
// 				logs = append(logs, "Button 1 clicked")
// 			}
// 			btn := material.Button(th, &button1, "Button 1")
// 			return btn.Layout(gtx)
// 		}),
// 		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
// 			for button2.Clicked(gtx) {
// 				logs = append(logs, "Button 2 clicked")
// 			}
// 			btn := material.Button(th, &button2, "Button 2")
// 			return btn.Layout(gtx)
// 		}),
// 	)
// }

// func layoutLogs(gtx layout.Context, th *material.Theme) layout.Dimensions {
// 	list := layout.List{Axis: layout.Vertical}
// 	return list.Layout(gtx, len(logs), func(gtx layout.Context, i int) layout.Dimensions {
// 		log := logs[i]
// 		lbl := material.Body1(th, log)
// 		lbl.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
// 		return lbl.Layout(gtx)
// 	})
// }
