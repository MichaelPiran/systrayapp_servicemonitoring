package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func setupDashboard(myApp fyne.App) (fyne.Window, error) {

	// Set the application icon
	icon, _ := fyne.LoadResourceFromPath("images/icon.png")
	myApp.SetIcon(icon)

	myWindow := myApp.NewWindow("Fyne App with Custom Icon")
	myWindow.SetContent(widget.NewLabel("Fyne System Tray"))
	myWindow.SetCloseIntercept(func() {
		myWindow.Hide() // Hide the window instead of closing it
	})
	running_icon, _ := fyne.LoadResourceFromPath(serviceRunningIcon)
	myWindow.SetIcon(running_icon)
	return myWindow, nil
}

// func setupDashboard() fyne.Window {
// 	myApp := app.New()

// 	// Set the application icon
// 	icon, _ := fyne.LoadResourceFromPath("images/icon.png")
// 	myApp.SetIcon(icon)

// 	myWindow := myApp.NewWindow("Fyne App with Custom Icon")

// 	content := container.NewVBox(
// 		widget.NewLabel("Hello Fyne!"),
// 		widget.NewButton("Click Me!", func() {
// 			println("Button clicked!")
// 		}),
// 	)

// 	myWindow.SetContent(content)
// 	myWindow.Resize(fyne.NewSize(300, 200))

// 	return myWindow
// }

func showDashboard() {

	// myWindow.ShowAndRun()
}
