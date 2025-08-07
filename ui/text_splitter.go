package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TextSplitterUI implémente Processor pour le découpage de texte
type TextSplitterUI struct {
	viewModel *TextSplitterViewModel
}

func NewTextSplitterUI() Processor {
	return &TextSplitterUI{
		viewModel: NewTextSplitterViewModel(),
	}
}

func (ui *TextSplitterUI) Name() string {
	return "Découpeur de Texte"
}

func (ui *TextSplitterUI) Description() string {
	return "Découpe le texte selon un délimiteur personnalisable"
}

func (ui *TextSplitterUI) ViewModel() ViewModel {
	return ui.viewModel
}

func (ui *TextSplitterUI) CreateConfigurationUI() fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Entrez le texte à découper...")
	input.Wrapping = fyne.TextWrapWord
	input.Resize(fyne.NewSize(0, 120))

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.Disable()

	delimiterEntry := widget.NewEntry()
	delimiterEntry.SetPlaceHolder("Délimiteur (laisser vide pour \\n)")

	processBtn := widget.NewButton("Découper", func() {
		result, err := ui.viewModel.Process(input.Text)
		if err != nil {
			output.SetText(PrettyValidationError(err))
		} else {
			output.SetText(result)
		}
	})

	copyBtn := widget.NewButton("Copier", func() {
		if result, _ := ui.viewModel.GetLastResult(); result != "" {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(result)
		}
	})

	delimiterEntry.OnChanged = func(s string) {
		ui.viewModel.delimiter = s
	}

	topSection := container.NewVBox(
		widget.NewLabel("Entrée Texte:"),
		input,
		container.NewHBox(
			widget.NewLabel("Délimiteur:"),
			delimiterEntry,
			processBtn,
			copyBtn,
		),
		widget.NewLabel("Résultat découpé:"),
	)

	return container.NewBorder(
		topSection,
		nil,
		nil,
		nil,
		container.NewVScroll(output),
	)
}

// TextSplitterViewModel implémente ViewModel pour le découpage
type TextSplitterViewModel struct {
	delimiter  string
	lastResult string
}

func NewTextSplitterViewModel() *TextSplitterViewModel {
	return &TextSplitterViewModel{
		delimiter: "\n",
	}
}

func (vm *TextSplitterViewModel) Process(input string) (string, error) {
	delim := vm.delimiter
	if delim == "" {
		delim = "\n"
	}

	parts := strings.Split(input, delim)
	return strings.Join(parts, "\n"), nil
}

func (vm *TextSplitterViewModel) GetConfiguration() interface{} {
	return struct {
		Delimiter string
	}{
		Delimiter: vm.delimiter,
	}
}

func (vm *TextSplitterViewModel) LoadConfiguration(config interface{}) error {
	cfg, ok := config.(struct{ Delimiter string })
	if !ok {
		return fmt.Errorf("configuration invalide")
	}
	vm.delimiter = cfg.Delimiter
	return nil
}

func (vm *TextSplitterViewModel) Validate() error {
	return nil // Pas de validation nécessaire
}

func (vm *TextSplitterViewModel) GetLastResult() (string, error) {
	return vm.lastResult, nil
}
