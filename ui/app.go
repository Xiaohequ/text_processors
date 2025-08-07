package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
		// Créer une boîte de dialogue pour sélectionner l'emplacement de sauvegarde
		fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				fmt.Printf("Erreur lors de la création du fichier : %v\n", err)
				errorDialog := dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				errorDialog.Show()
				return
			}

			if writer == nil {
				// L'utilisateur a annulé
				return
			}
			defer writer.Close()

			// Obtenir le chemin du fichier
			filePath := writer.URI().Path()

			// Exporter le pipeline vers le fichier sélectionné
			pipeline := CurrentPipeline // Utiliser le pipeline global du package ui
			err = pipeline.SaveToFile(filePath)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation vers %s : %v\n", filePath, err)
				// Afficher une boîte de dialogue d'erreur
				errorDialog := dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				errorDialog.Show()
			} else {
				fmt.Printf("Pipeline exporté avec succès vers %s\n", filePath)
				// Afficher une confirmation
				infoDialog := dialog.NewInformation("Export réussi",
					fmt.Sprintf("Pipeline exporté avec succès vers :\n%s", filePath),
					fyne.CurrentApp().Driver().AllWindows()[0])
				infoDialog.Show()
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])

		// Configurer le filtre pour les fichiers JSON et le nom par défaut
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.SetFileName("mon_pipeline.json")

		// Afficher la boîte de dialogue
		fileDialog.Show()
	})
	importBtn = widget.NewButton("Import", func() {
		// Créer une boîte de dialogue pour sélectionner le fichier à importer
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Printf("Erreur lors de l'ouverture du fichier : %v\n", err)
				errorDialog := dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				errorDialog.Show()
				return
			}

			if reader == nil {
				// L'utilisateur a annulé
				return
			}
			defer reader.Close()

			// Obtenir le chemin du fichier
			filePath := reader.URI().Path()

			// Charger le pipeline depuis le fichier sélectionné
			err = CurrentPipeline.LoadFromFile(filePath)
			if err != nil {
				fmt.Printf("Erreur lors de l'importation depuis %s : %v\n", filePath, err)
				// Afficher une boîte de dialogue d'erreur
				errorDialog := dialog.NewError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				errorDialog.Show()
			} else {
				fmt.Printf("Pipeline importé avec succès depuis %s\n", filePath)
				// Notifier que le pipeline a été mis à jour
				NotifyPipelineUpdated()
				// Afficher une confirmation
				infoDialog := dialog.NewInformation("Import réussi",
					fmt.Sprintf("Pipeline importé avec succès depuis :\n%s", filePath),
					fyne.CurrentApp().Driver().AllWindows()[0])
				infoDialog.Show()
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])

		// Configurer le filtre pour les fichiers JSON
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))

		// Afficher la boîte de dialogue
		fileDialog.Show()
	})

	// Créer le layout final avec scroll
	finalContent := container.NewScroll(mainContent)
	mainContent = container.NewStack(finalContent)

	// Commencer par afficher la grille
	showToolsGrid()

	return mainContent
}
