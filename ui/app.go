package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeUI() fyne.CanvasObject {
	// Déclarer d'abord les boutons
	var exportBtn *widget.Button
	var importBtn *widget.Button

	// Container principal qui contiendra soit la grille, soit l'outil sélectionné
	mainContent := container.NewStack()

	var showToolsGrid func()
	var backBtn *widget.Button

	// Fonction pour afficher la grille des outils
	showToolsGrid = func() {
		toolsGrid := MakeToolsGrid(func(toolUI fyne.CanvasObject) {
			// Quand un outil est sélectionné, l'afficher avec le bouton retour
			mainContent.Objects = []fyne.CanvasObject{
				container.NewBorder(
					backBtn,
					container.NewVBox(
						widget.NewSeparator(),
						container.NewCenter(
							container.NewHBox(
								exportBtn,
								importBtn,
							),
						),
					),
					nil,
					nil,
					toolUI,
				),
			}
			mainContent.Refresh()
		})

		// Afficher la grille sans le bouton retour
		mainContent.Objects = []fyne.CanvasObject{toolsGrid}
		mainContent.Refresh()
	}

	// Bouton retour
	backBtn = widget.NewButton("← Retour", showToolsGrid)

	// Initialiser les boutons
	exportBtn = widget.NewButton("Export", func() {
		// Logique pour exporter la configuration
		pipeline := CurrentPipeline // Utiliser le pipeline global du package ui
		err := pipeline.SaveToFile("exported_pipeline.json")
		if err != nil {
			fmt.Println("Erreur lors de l'exportation :", err)
		} else {
			fmt.Println("Pipeline exporté avec succès.")
		}
	})
	importBtn = widget.NewButton("Import", func() {
		// Logique pour importer la configuration
		pipeline := &Pipeline{}
		err := pipeline.LoadFromFile("imported_pipeline.json")
		if err != nil {
			fmt.Println("Erreur lors de l'importation :", err)
		} else {
			fmt.Println("Pipeline importé avec succès.")
		}
	})

	// Créer le layout final avec scroll
	finalContent := container.NewScroll(mainContent)
	mainContent = container.NewStack(finalContent)

	// Commencer par afficher la grille
	showToolsGrid()

	return mainContent
}
