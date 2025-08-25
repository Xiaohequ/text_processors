package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"text_processors/ui/processors"
)

// MakePipelineBuilderUI crée l'interface du constructeur de pipeline
func MakePipelineBuilderUI() fyne.CanvasObject {
	// Utiliser la variable globale du pipeline actuel
	currentPipeline := CurrentPipeline

	// Zone d'affichage des étapes
	stepsContainer := container.NewVBox()

	// Zone de texte d'entrée
	inputText := widget.NewMultiLineEntry()
	inputText.SetPlaceHolder("Entrez votre texte ici...")
	inputText.Wrapping = fyne.TextWrapWord
	inputText.Resize(fyne.NewSize(0, 100))

	// Zone de résultat
	outputText := widget.NewMultiLineEntry()
	outputText.Wrapping = fyne.TextWrapWord
	outputText.MultiLine = true

	var resultText string

	// Déclaration de la fonction pour mettre à jour l'affichage des étapes
	var updateStepsDisplay func()

	// Fonction pour mettre à jour l'affichage des étapes
	updateStepsDisplay = func() {
		stepsContainer.Objects = nil

		if len(currentPipeline.Steps) == 0 {
			stepsContainer.Add(widget.NewLabel("Aucune étape configurée"))
		} else {
			for i, step := range currentPipeline.Steps {
				stepIndex := i // Capture pour la closure

				// Conteneur pour une étape
				stepContainer := container.NewHBox()

				// Numéro et nom de l'étape
				stepLabel := widget.NewLabel(fmt.Sprintf("%d. %s", i+1, step.Processor.Name()))
				stepContainer.Add(stepLabel)

				// Bouton monter
				if i > 0 {
					upBtn := widget.NewButton("↑", func() {
						// Échanger avec l'étape précédente
						currentPipeline.Steps[stepIndex], currentPipeline.Steps[stepIndex-1] =
							currentPipeline.Steps[stepIndex-1], currentPipeline.Steps[stepIndex]
						updateStepsDisplay()
					})
					stepContainer.Add(upBtn)
				}

				// Bouton descendre
				if i < len(currentPipeline.Steps)-1 {
					downBtn := widget.NewButton("↓", func() {
						// Échanger avec l'étape suivante
						currentPipeline.Steps[stepIndex], currentPipeline.Steps[stepIndex+1] =
							currentPipeline.Steps[stepIndex+1], currentPipeline.Steps[stepIndex]
						updateStepsDisplay()
					})
					stepContainer.Add(downBtn)
				}

				// Bouton supprimer
				deleteBtn := widget.NewButton("×", func() {
					// Supprimer l'étape
					currentPipeline.Steps = append(
						currentPipeline.Steps[:stepIndex],
						currentPipeline.Steps[stepIndex+1:]...)
					updateStepsDisplay()
				})
				stepContainer.Add(deleteBtn)

				stepsContainer.Add(stepContainer)
			}
		}

		stepsContainer.Refresh()
	}

	// Fonction pour obtenir la liste des outils disponibles
	getToolOptions := func() []string {
		options := []string{"JSON Formatter", "Text Splitter", "Text Joiner"}
		// Ajouter les processeurs personnalisés
		for _, customProc := range GlobalCustomProcessorManager.GetProcessors() {
			options = append(options, "Custom: "+customProc.Name)
		}
		return options
	}

	// Sélecteur d'outil à ajouter
	toolSelect := widget.NewSelect(getToolOptions(), nil)
	toolSelect.SetSelected("JSON Formatter")

	// Zone d'affichage des erreurs
	errorLabel := widget.NewLabel("")
	errorLabel.Wrapping = fyne.TextWrapWord

	// Fonction pour afficher les erreurs
	showError := func(err error) {
		if err != nil {
			errorLabel.SetText(fmt.Sprintf("Erreur: %s", err.Error()))
			errorLabel.Importance = widget.HighImportance
		} else {
			errorLabel.SetText("")
		}
	}

	// Bouton pour ajouter l'étape au pipeline
	addStepBtn := widget.NewButton("Ajouter l'étape", func() {
		// Effacer les erreurs précédentes
		showError(nil)

		if toolSelect.Selected == "" {
			showError(fmt.Errorf("veuillez sélectionner un outil"))
			return
		}

		// Create processor instance
		var processor processors.Processor
		switch toolSelect.Selected {
		case "JSON Formatter":
			processor = processors.NewJSONFormatterUI()
		case "Text Splitter":
			processor = processors.NewTextSplitterUI()
		case "Text Joiner":
			processor = processors.NewTextJoinerUI()
		default:
			// Vérifier si c'est un processeur personnalisé
			if strings.HasPrefix(toolSelect.Selected, "Custom: ") {
				customName := strings.TrimPrefix(toolSelect.Selected, "Custom: ")
				// Trouver le processeur personnalisé
				for _, customProc := range GlobalCustomProcessorManager.GetProcessors() {
					if customProc.Name == customName {
						processor = processors.NewCustomProcessor(customProc.Name, customProc.Script)
						break
					}
				}
			}
		}

		// Declare dialog first
		var configDialog *dialog.CustomDialog

		// Create dialog content with buttons
		content := container.NewBorder(
			nil,
			container.NewCenter(
				container.NewHBox(
					widget.NewButton("Annuler", func() { configDialog.Hide() }),
					widget.NewButton("Valider", func() {
						// Valider la configuration du processeur
						if processor.ViewModel().Validate() == nil {
							configDialog.Hide()

							// Déterminer le type d'outil et la configuration
							var toolType ToolType
							var config ToolConfig

							switch toolSelect.Selected {
							case "JSON Formatter":
								toolType = JSONFormatterTool
								config = JSONFormatterConfig{IndentType: "2 espaces"} // Default value
							case "Text Splitter":
								toolType = TextSplitterTool  
								config = TextSplitterConfig{Delimiter: ","} // Default value
							case "Text Joiner":
								toolType = TextJoinerTool
								config = TextJoinerConfig{Delimiter: ", "} // Default value
							default:
								// Vérifier si c'est un processeur personnalisé
								if strings.HasPrefix(toolSelect.Selected, "Custom: ") {
									toolType = CustomProcessorTool
									customName := strings.TrimPrefix(toolSelect.Selected, "Custom: ")
									// Trouver le processeur personnalisé
									for _, customProc := range GlobalCustomProcessorManager.GetProcessors() {
										if customProc.Name == customName {
											config = CustomProcessorConfig{Name: customProc.Name, Script: customProc.Script}
											break
										}
									}
								}
							}

							step := PipelineStep{
								ID:        fmt.Sprintf("step_%d", len(currentPipeline.Steps)+1),
								Type:      toolType,
								Config:    config,
								Name:      "",
								Processor: processor,
							}
							currentPipeline.Steps = append(currentPipeline.Steps, step)
							updateStepsDisplay()
						} else {
							showError(fmt.Errorf("configuration invalide"))
						}
					}),
				),
			),
			nil,
			nil,
			processor.CreateConfigurationUI(),
		)

		// Initialize dialog and show
		configDialog = dialog.NewCustom(
			"Configuration du processeur",
			"Fermer",
			content,
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
		
		// Adapter la taille du dialog à la taille de la fenêtre
		window := fyne.CurrentApp().Driver().AllWindows()[0]
		windowSize := window.Canvas().Size()
		// Utiliser 80% de la largeur et 70% de la hauteur de la fenêtre
		dialogWidth := windowSize.Width * 0.8
		dialogHeight := windowSize.Height * 0.7
		configDialog.Resize(fyne.NewSize(dialogWidth, dialogHeight))
		
		configDialog.Show()
	})

	// Bouton pour exécuter le pipeline
	executeBtn := widget.NewButton("Exécuter le Pipeline", func() {
		// Effacer les erreurs précédentes
		showError(nil)

		if inputText.Text == "" {
			showError(fmt.Errorf("veuillez entrer du texte à traiter"))
			outputText.SetText("")
			resultText = ""
			return
		}

		if len(currentPipeline.Steps) == 0 {
			showError(fmt.Errorf("le pipeline doit contenir au moins une étape"))
			outputText.SetText("")
			resultText = ""
			return
		}

		executor := GetDefaultExecutor()
		result, err := executor.Execute(currentPipeline, inputText.Text)

		if err != nil {
			showError(err)
			resultText = ""
			outputText.SetText("")
		} else {
			resultText = result
			outputText.SetText(resultText)
		}
	})

	// Bouton pour copier le résultat
	copyBtn := widget.NewButton("Copier", func() {
		if resultText != "" {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(resultText)
		}
	})

	// Bouton pour vider le pipeline
	clearBtn := widget.NewButton("Vider le Pipeline", func() {
		currentPipeline.Steps = []PipelineStep{}
		updateStepsDisplay()
	})

	// Section de configuration du pipeline
	configSection := container.NewBorder(
		widget.NewCard("Configuration du Pipeline", "", container.NewVBox(
			container.NewHBox(
				widget.NewLabel("Ajouter un outil:"),
				toolSelect,
				addStepBtn,
			),
			errorLabel, // Zone d'affichage des erreurs
		)),
		nil, nil, nil,
		widget.NewCard("Étapes du Pipeline", "", container.NewBorder(
			container.NewHBox(clearBtn), nil, nil, nil,
			container.NewScroll(stepsContainer),
		)),
	)

	// Section d'exécution
	executionSection := container.NewVBox(
		widget.NewLabel("Texte d'entrée:"),
		inputText,
		container.NewHBox(executeBtn),
		container.NewHBox(
			widget.NewLabel("Résultat:"),
			copyBtn,
		),
		outputText,
	)

	// Initialiser l'affichage des étapes
	updateStepsDisplay()

	// Fonction pour rafraîchir la liste des outils
	refreshToolList := func() {
		currentSelected := toolSelect.Selected
		toolSelect.Options = getToolOptions()
		// Essayer de garder la sélection actuelle si elle existe encore
		found := false
		for _, option := range toolSelect.Options {
			if option == currentSelected {
				toolSelect.SetSelected(currentSelected)
				found = true
				break
			}
		}
		if !found && len(toolSelect.Options) > 0 {
			toolSelect.SetSelected(toolSelect.Options[0])
		}
		toolSelect.Refresh()
	}

	// Enregistrer le callback pour rafraîchir quand des processeurs personnalisés sont ajoutés
	GlobalCustomProcessorManager.SetUpdateCallback(refreshToolList)

	// Effacer les anciens callbacks et enregistrer le nouveau pour les importations
	ClearPipelineUpdateCallbacks()
	RegisterPipelineUpdateCallback(updateStepsDisplay)

	// Layout principal
	return container.NewHSplit(
		configSection,
		executionSection,
	)
}
