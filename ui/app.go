package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeUI() fyne.CanvasObject {
	// Container principal qui contiendra soit la grille, soit l'outil sélectionné
	mainContent := container.NewStack()

	var showToolsGrid func()
	var backBtn *widget.Button

	// Fonction pour afficher la grille des outils
	showToolsGrid = func() {
		toolsGrid := MakeToolsGrid(func(toolUI fyne.CanvasObject) {
			// Quand un outil est sélectionné, l'afficher avec le bouton retour
			mainContent.Objects = []fyne.CanvasObject{
				container.NewBorder(backBtn, nil, nil, nil, toolUI),
			}
			mainContent.Refresh()
		})

		// Afficher la grille sans le bouton retour
		mainContent.Objects = []fyne.CanvasObject{toolsGrid}
		mainContent.Refresh()
	}

	// Bouton retour
	backBtn = widget.NewButton("← Retour", showToolsGrid)

	// Commencer par afficher la grille
	showToolsGrid()

	return mainContent
}
