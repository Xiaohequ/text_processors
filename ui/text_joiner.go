package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TextJoinerUI implémente Processor pour la jointure de texte
type TextJoinerUI struct {
	viewModel *TextJoinerViewModel
}

func NewTextJoinerUI() Processor {
	return &TextJoinerUI{
		viewModel: NewTextJoinerViewModel(),
	}
}

func (ui *TextJoinerUI) Name() string {
	return "Joigneur de Texte"
}

func (ui *TextJoinerUI) Description() string {
	return "Assemble des lignes de texte avec un délimiteur"
}

func (ui *TextJoinerUI) ViewModel() ViewModel {
	return ui.viewModel
}

func (ui *TextJoinerUI) CreateConfigurationUI() fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Entrez le texte à assembler...")
	input.Wrapping = fyne.TextWrapWord
	input.Resize(fyne.NewSize(0, 120))

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.Disable()

	delimiterEntry := widget.NewEntry()
	delimiterEntry.SetPlaceHolder("Délimiteur (ex: , )")

	processBtn := widget.NewButton("Assembler", func() {
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
		widget.NewLabel("Résultat assemblé:"),
	)

	return container.NewBorder(
		topSection,
		nil,
		nil,
		nil,
		container.NewVScroll(output),
	)
}

// TextJoinerViewModel implémente ViewModel pour la jointure
type TextJoinerViewModel struct {
	delimiter  string
	lastResult string
}

func NewTextJoinerViewModel() *TextJoinerViewModel {
	return &TextJoinerViewModel{
		delimiter: " ",
	}
}

func (vm *TextJoinerViewModel) Process(input string) (string, error) {
	lines := strings.Split(input, "\n")
	// Filtrer les lignes vides
	var nonEmptyLines []string
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			nonEmptyLines = append(nonEmptyLines, trimmed)
		}
	}
	vm.lastResult = strings.Join(nonEmptyLines, vm.delimiter)
	return vm.lastResult, nil
}

func (vm *TextJoinerViewModel) GetConfiguration() interface{} {
	return struct {
		Delimiter string
	}{
		Delimiter: vm.delimiter,
	}
}

func (vm *TextJoinerViewModel) LoadConfiguration(config interface{}) error {
	cfg, ok := config.(struct{ Delimiter string })
	if !ok {
		return fmt.Errorf("configuration invalide")
	}
	vm.delimiter = cfg.Delimiter
	return nil
}

func (vm *TextJoinerViewModel) Validate() error {
	return nil // Pas de validation nécessaire
}

func (vm *TextJoinerViewModel) GetLastResult() (string, error) {
	return vm.lastResult, nil
}
