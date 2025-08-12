package main

import (
	"fmt"
	"text_processors/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	// Charger les processeurs personnalisés depuis conf/ au démarrage
	if err := ui.GlobalCustomProcessorManager.LoadAll(); err != nil {
		fmt.Printf("[WARN] Échec du chargement des processeurs personnalisés: %v\n", err)
	}

	myWindow := myApp.NewWindow("Text processors")
	myWindow.Resize(fyne.NewSize(800, 600))

	content := ui.MakeUI()
	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
