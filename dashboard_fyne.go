package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func showDashboard() {
	myApp := app.New()

	// Set the application icon
	icon, _ := fyne.LoadResourceFromPath("images/icon.png")
	myApp.SetIcon(icon)

	myWindow := myApp.NewWindow("Fyne App with Custom Icon")

	content := container.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Click Me!", func() {
			println("Button clicked!")
		}),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}
