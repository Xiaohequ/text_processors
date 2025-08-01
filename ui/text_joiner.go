package ui

import (
    "strings"
    
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func MakeTextJoinerUI() fyne.CanvasObject {
    // Zone de texte d'entrée
    input := widget.NewMultiLineEntry()
    input.SetPlaceHolder("Entrez le texte à joindre (une ligne par élément)...")
    input.Wrapping = fyne.TextWrapWord
    input.Resize(fyne.NewSize(0, 120))

    // Champ délimiteur
    delimiter := widget.NewEntry()
    delimiter.SetPlaceHolder("Délimiteur (ex: ,)")
    delimiter.SetText(", ")

    // Zone de résultat
    output := widget.NewMultiLineEntry()
    output.Wrapping = fyne.TextWrapWord
    output.MultiLine = true

    var resultText string

    // Bouton joindre
    joinBtn := widget.NewButton("Joindre", func() {
        inputText := input.Text
        delimiterText := delimiter.Text
        
        if inputText == "" {
            output.SetText("")
            resultText = ""
            return
        }

        lines := strings.Split(inputText, "\n")
        // Filtrer les lignes vides
        var nonEmptyLines []string
        for _, line := range lines {
            if strings.TrimSpace(line) != "" {
                nonEmptyLines = append(nonEmptyLines, strings.TrimSpace(line))
            }
        }

        resultText = strings.Join(nonEmptyLines, delimiterText)
        output.SetText(resultText)
    })

    // Bouton copier
    copyBtn := widget.NewButton("Copier", func() {
        if resultText != "" {
            clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
            clipboard.SetContent(resultText)
        }
    })

    topSection := container.NewVBox(
        widget.NewLabel("Texte à joindre (une ligne par élément):"),
        input,
        container.NewHBox(
            widget.NewLabel("Délimiteur:"),
            delimiter,
            joinBtn,
        ),
        container.NewHBox(
            widget.NewLabel("Résultat:"),
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