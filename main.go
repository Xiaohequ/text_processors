package main

import (
	"jsonformatter/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	myWindow := myApp.NewWindow("JSON Formatter")
	myWindow.Resize(fyne.NewSize(800, 600))

	content := ui.MakeUI()
	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
