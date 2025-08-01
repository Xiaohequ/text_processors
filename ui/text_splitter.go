package ui

import (
    "strings"
    
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func MakeTextSplitterUI() fyne.CanvasObject {
    // Zone de texte d'entrée
    input := widget.NewMultiLineEntry()
    input.SetPlaceHolder("Entrez le texte à diviser...")
    input.Wrapping = fyne.TextWrapWord
    input.Resize(fyne.NewSize(0, 120))

    // Champ délimiteur
    delimiter := widget.NewEntry()
    delimiter.SetPlaceHolder("Délimiteur (ex: ,)")
    delimiter.SetText(",")

    // Zone de résultat
    output := widget.NewMultiLineEntry()
    output.Wrapping = fyne.TextWrapWord
    output.MultiLine = true

    var resultText string

    // Bouton diviser
    splitBtn := widget.NewButton("Diviser", func() {
        inputText := input.Text
        delimiterText := delimiter.Text
        
        if inputText == "" {
            output.SetText("")
            resultText = ""
            return
        }

        if delimiterText == "" {
            delimiterText = "\n"
        }

        parts := strings.Split(inputText, delimiterText)
        resultText = strings.Join(parts, "\n")
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
        widget.NewLabel("Texte à diviser:"),
        input,
        container.NewHBox(
            widget.NewLabel("Délimiteur:"),
            delimiter,
            splitBtn,
        ),
        container.NewHBox(
            widget.NewLabel("Résultat (une ligne par partie):"),
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