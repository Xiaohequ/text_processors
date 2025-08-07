package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// MakePipelineBuilderUI crée l'interface du constructeur de pipeline
func MakePipelineBuilderUI() fyne.CanvasObject {
	// Variable globale pour le pipeline actuel
	var CurrentPipeline = &Pipeline{
		Name:  "Mon Pipeline",
		Steps: []PipelineStep{},
	}

	// Pipeline en cours de construction
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

	// Sélecteur d'outil à ajouter
	toolOptions := []string{"JSON Formatter", "Text Splitter", "Text Joiner"}
	toolSelect := widget.NewSelect(toolOptions, nil)
	toolSelect.SetSelected("JSON Formatter")

	// Zone de configuration pour l'outil sélectionné
	configContainer := container.NewVBox()

	// Variables pour les configurations
	var jsonIndentSelect *widget.Select
	var splitterDelimiterEntry *widget.Entry
	var joinerDelimiterEntry *widget.Entry

	// Fonction pour mettre à jour la zone de configuration
	updateConfigDisplay := func(toolName string) {
		configContainer.Objects = nil

		switch toolName {
		case "JSON Formatter":
			configContainer.Add(widget.NewLabel("Configuration JSON Formatter:"))
			indentOptions := []string{"2 espaces", "4 espaces", "Tabulations"}
			jsonIndentSelect = widget.NewSelect(indentOptions, nil)
			jsonIndentSelect.SetSelected("2 espaces")
			configContainer.Add(container.NewHBox(
				widget.NewLabel("Indentation:"),
				jsonIndentSelect,
			))

		case "Text Splitter":
			configContainer.Add(widget.NewLabel("Configuration Text Splitter:"))
			splitterDelimiterEntry = widget.NewEntry()
			splitterDelimiterEntry.SetPlaceHolder("Délimiteur (ex: ,)")
			splitterDelimiterEntry.SetText(",")
			configContainer.Add(container.NewHBox(
				widget.NewLabel("Délimiteur:"),
				splitterDelimiterEntry,
			))

		case "Text Joiner":
			configContainer.Add(widget.NewLabel("Configuration Text Joiner:"))
			joinerDelimiterEntry = widget.NewEntry()
			joinerDelimiterEntry.SetPlaceHolder("Délimiteur (ex: ,)")
			joinerDelimiterEntry.SetText(", ")
			configContainer.Add(container.NewHBox(
				widget.NewLabel("Délimiteur:"),
				joinerDelimiterEntry,
			))
		}

		configContainer.Refresh()
	}

	// Initialiser l'affichage de configuration
	updateConfigDisplay("JSON Formatter")

	// Mettre à jour la configuration quand l'outil change
	toolSelect.OnChanged = updateConfigDisplay

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
		var config ToolConfig
		var err error

		switch toolSelect.Selected {
		case "JSON Formatter":
			if jsonIndentSelect.Selected == "" {
				showError(fmt.Errorf("veuillez sélectionner un type d'indentation"))
				return
			}
			config = JSONFormatterConfig{
				IndentType: jsonIndentSelect.Selected,
			}
		case "Text Splitter":
			config = TextSplitterConfig{
				Delimiter: splitterDelimiterEntry.Text,
			}
		case "Text Joiner":
			config = TextJoinerConfig{
				Delimiter: joinerDelimiterEntry.Text,
			}
		default:
			showError(fmt.Errorf("veuillez sélectionner un outil"))
			return
		}

		if err = config.Validate(); err != nil {
			showError(err)
			return
		}

		// Effacer les erreurs précédentes
		showError(nil)

		// Create processor instance
		var processor Processor
		switch toolSelect.Selected {
		case "JSON Formatter":
			processor = NewJSONFormatterUI()
		case "Text Splitter":
			processor = NewTextSplitterUI()
		case "Text Joiner":
			processor = NewTextJoinerUI()
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
						if processor.ViewModel().Validate() == nil {
							configDialog.Hide()
							step := PipelineStep{
								ID:        fmt.Sprintf("step_%d", len(currentPipeline.Steps)+1),
								Processor: processor,
								Name:      "",
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
		configDialog.Show()

		// Remove old configuration UI elements
		configContainer.Objects = nil
		configContainer.Refresh()
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
			configContainer,
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

	// Layout principal
	return container.NewHSplit(
		configSection,
		executionSection,
	)
}
