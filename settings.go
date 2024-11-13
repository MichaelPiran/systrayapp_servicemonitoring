package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func openSettings() {
	// Create a new window for settings
	settingsWindow := myApp.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(400, 300))

	// Create checkboxes for each setting in the first tab
	option1Check := widget.NewCheck("Option 1", func(checked bool) {
		option1 = checked
	})
	option1Check.Checked = option1

	option2Check := widget.NewCheck("Option 2", func(checked bool) {
		option2 = checked
	})
	option2Check.Checked = option2

	option3Check := widget.NewCheck("Option 3", func(checked bool) {
		option3 = checked
	})
	option3Check.Checked = option3

	// Create containers to hold the checkboxes for each tab
	tab1Content := container.NewVBox(
		option1Check,
		option2Check,
		option3Check,
	)

	button := widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
		os.Exit(0)
	})
	button.Resize(fyne.NewSize(10, button.MinSize().Height)) // Set fixed width to 100

	tab2Content := container.NewVBox(
		widget.NewLabelWithStyle("App Version: "+appVersion, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		button,
	)

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Service", tab1Content),
		container.NewTabItem("About", tab2Content),
	)

	// Set the content of the settings window
	settingsWindow.SetContent(tabs)

	settingsWindow.SetCloseIntercept(func() {
		settingsWindow.Hide() // Hide the window instead of closing it
	})

	// Load and set the running icon for the window
	appIcon, err := icons.ReadFile(serviceAppIcon)
	if err != nil {
		log.Printf("Failed to load running icon: %v", err)
		// Consider whether you want to return the error or continue without the icon
	} else {
		settingsWindow.SetIcon(fyne.NewStaticResource("runningIcon", appIcon))
	}

	// Show the settings window
	settingsWindow.Show()
}
