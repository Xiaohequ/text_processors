package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Tool struct {
	Name        string
	Description string
	Icon        fyne.Resource
	MakeUI      func() fyne.CanvasObject
}

func MakeToolsGrid(onToolSelect func(fyne.CanvasObject)) fyne.CanvasObject {
	tools := []Tool{
		{
			Name:        "Pipeline Builder",
			Description: "Enchaîne plusieurs outils de traitement",
			MakeUI:      MakePipelineBuilderUI,
		},
		{
			Name:        "JSON Formatter",
			Description: "Formate et valide du JSON",
			MakeUI:      MakeJSONFormatterUI,
		},
		{
			Name:        "Text Splitter",
			Description: "Divise du texte selon un délimiteur",
			MakeUI:      MakeTextSplitterUI,
		},
		{
			Name:        "Text Joiner",
			Description: "Joint du texte avec un délimiteur",
			MakeUI:      MakeTextJoinerUI,
		},
	}

	// Créer une grille qui s'adapte à l'espace disponible
	// Utiliser 2 colonnes pour que les cartes soient bien disposées
	grid := container.NewGridWithColumns(2)

	for _, tool := range tools {
		toolCopy := tool // Capture pour la closure

		card := widget.NewCard(toolCopy.Name, toolCopy.Description,
			widget.NewButton("Ouvrir", func() {
				onToolSelect(toolCopy.MakeUI())
			}))

		grid.Add(card)
	}

	// Titre centré
	title := widget.NewLabel("Sélectionnez un outil")
	title.Alignment = fyne.TextAlignCenter

	// Ajouter du padding autour de la grille pour une meilleure présentation
	paddedGrid := container.NewPadded(grid)

	// Utiliser NewBorder pour que la grille prenne tout l'espace disponible
	return container.NewBorder(
		container.NewPadded(title), // top avec padding
		nil,                        // bottom
		nil,                        // left
		nil,                        // right
		paddedGrid,                 // center - prend tout l'espace restant
	)
}
