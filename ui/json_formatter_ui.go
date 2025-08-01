package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeJSONFormatterUI() fyne.CanvasObject {
	// Zone de texte pour l'entrée JSON (taille limitée)
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Entrez votre JSON mal formaté ici...")
	input.Wrapping = fyne.TextWrapWord
	input.Resize(fyne.NewSize(0, 120))

	// Zone de texte pour le résultat formaté
	output := widget.NewRichTextFromMarkdown("")
	output.Wrapping = fyne.TextWrapWord
	output.Scroll = container.ScrollBoth

	var formattedText string

	// Options d'indentation
	indentOptions := []string{"2 espaces", "4 espaces", "Tabulations"}
	indentSelect := widget.NewSelect(indentOptions, func(s string) {})
	indentSelect.SetSelected("2 espaces")

	// Bouton de formatage
	formatBtn := widget.NewButton("Formater", func() {
		inputText := input.Text
		if inputText == "" {
			formattedText = ""
			output.ParseMarkdown("```\n\n```")
			return
		}

		formatter := NewFormatter(indentSelect.Selected)
		formatted, err := formatter.FormatJSON(inputText)
		if err != nil {
			formattedText = PrettyValidationError(err)
			output.ParseMarkdown("```\n" + formattedText + "\n```")
		} else {
			formattedText = formatted
			output.ParseMarkdown("```json\n" + formattedText + "\n```")
		}
	})

	// Bouton pour copier le résultat
	copyBtn := widget.NewButton("Copier", func() {
		if formattedText != "" {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(formattedText)
		}
	})

	// Rafraîchir le formatage si l'indentation change
	indentSelect.OnChanged = func(s string) {
		if input.Text != "" && formattedText != "" {
			formatBtn.OnTapped()
		}
	}

	topSection := container.NewVBox(
		widget.NewLabel("Entrée JSON:"),
		input,
		container.NewHBox(
			formatBtn,
			widget.NewLabel("Indentation:"),
			indentSelect,
		),
		container.NewHBox(
			widget.NewLabel("Résultat formaté:"),
			copyBtn,
		),
	)

	return container.NewBorder(
		topSection,
		nil,
		nil,
		nil,
		output,
	)
}
