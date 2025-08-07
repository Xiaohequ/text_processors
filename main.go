package main

import (
	"text_processors/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	myWindow := myApp.NewWindow("Text processors")
	myWindow.Resize(fyne.NewSize(800, 600))

	content := ui.MakeUI()
	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
