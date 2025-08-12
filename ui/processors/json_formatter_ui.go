package processors

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// JSONFormatterUI implémente l'interface Processor pour le formateur JSON
type JSONFormatterUI struct {
	viewModel *JSONFormatterViewModel
}

func NewJSONFormatterUI() Processor {
	return &JSONFormatterUI{
		viewModel: NewJSONFormatterViewModel(),
	}
}

func (ui *JSONFormatterUI) Name() string {
	return "Formateur JSON"
}

func (ui *JSONFormatterUI) Description() string {
	return "Formate les documents JSON avec indentation personnalisée"
}

func (ui *JSONFormatterUI) ViewModel() ViewModel {
	return ui.viewModel
}

func (ui *JSONFormatterUI) CreateConfigurationUI() fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Entrez votre JSON mal formaté ici...")
	input.Wrapping = fyne.TextWrapWord
	input.Resize(fyne.NewSize(0, 120))

	output := widget.NewRichTextFromMarkdown("")
	output.Wrapping = fyne.TextWrapWord
	output.Scroll = container.ScrollBoth

	indentOptions := []string{"2 espaces", "4 espaces", "Tabulations"}
	indentSelect := widget.NewSelect(indentOptions, func(s string) {
		ui.viewModel.indentType = s
	})
	indentSelect.SetSelected(ui.viewModel.indentType)

	formatBtn := widget.NewButton("Formater", func() {
		result, err := ui.viewModel.Process(input.Text)
		if err != nil {
			output.ParseMarkdown("```\n" + PrettyValidationError(err) + "\n```")
		} else {
			output.ParseMarkdown("```json\n" + result + "\n```")
		}
	})

	copyBtn := widget.NewButton("Copier", func() {
		if result, _ := ui.viewModel.GetLastResult(); result != "" {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(result)
		}
	})

	indentSelect.OnChanged = func(s string) {
		if input.Text != "" {
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

// JSONFormatterViewModel implémente ViewModel pour le formateur JSON
type JSONFormatterViewModel struct {
	indentType string
	lastResult string
}

func NewJSONFormatterViewModel() *JSONFormatterViewModel {
	return &JSONFormatterViewModel{
		indentType: "2 espaces",
	}
}

func (vm *JSONFormatterViewModel) Process(input string) (string, error) {
	if input == "" {
		vm.lastResult = ""
		return "", nil
	}

	formatter := NewFormatter(vm.indentType)
	formatted, err := formatter.FormatJSON(input)
	if err != nil {
		return "", err
	}

	vm.lastResult = formatted
	return formatted, nil
}

func (vm *JSONFormatterViewModel) GetConfiguration() interface{} {
	return struct {
		IndentType string
	}{
		IndentType: vm.indentType,
	}
}

func (vm *JSONFormatterViewModel) LoadConfiguration(config interface{}) error {
	cfg, ok := config.(struct{ IndentType string })
	if !ok {
		return fmt.Errorf("configuration invalide")
	}
	vm.indentType = cfg.IndentType
	return nil
}

func (vm *JSONFormatterViewModel) Validate() error {
	validTypes := []string{"2 espaces", "4 espaces", "Tabulations"}
	for _, valid := range validTypes {
		if vm.indentType == valid {
			return nil
		}
	}
	return fmt.Errorf("type d'indentation invalide: %s", vm.indentType)
}

func (vm *JSONFormatterViewModel) GetLastResult() (string, error) {
	return vm.lastResult, nil
}
