package processors

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dop251/goja"
)

// CustomProcessor implémente Processor pour les processeurs personnalisés JavaScript
type CustomProcessor struct {
	viewModel *CustomProcessorViewModel
}

func NewCustomProcessor(name, script string) Processor {
	return &CustomProcessor{
		viewModel: NewCustomProcessorViewModel(name, script),
	}
}

func (cp *CustomProcessor) Name() string {
	if cp.viewModel.name != "" {
		return cp.viewModel.name
	}
	return "Processeur Personnalisé"
}

func (cp *CustomProcessor) Description() string {
	return "Processeur personnalisé utilisant JavaScript"
}

func (cp *CustomProcessor) ViewModel() ViewModel {
	return cp.viewModel
}

func (cp *CustomProcessor) CreateConfigurationUI() fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Entrez le texte à traiter...")
	input.Wrapping = fyne.TextWrapWord
	input.Resize(fyne.NewSize(0, 120))

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.Disable()

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom du processeur")
	nameEntry.SetText(cp.viewModel.name)

	scriptEntry := widget.NewMultiLineEntry()
	scriptEntry.SetPlaceHolder("Script JavaScript (fonction process(input) { return input; })")
	scriptEntry.SetText(cp.viewModel.script)
	scriptEntry.Wrapping = fyne.TextWrapWord
	scriptEntry.Resize(fyne.NewSize(0, 100))

	processBtn := widget.NewButton("Traiter", func() {
		// Mettre à jour le nom et le script
		cp.viewModel.name = nameEntry.Text
		cp.viewModel.script = scriptEntry.Text
		
		result, err := cp.viewModel.Process(input.Text)
		if err != nil {
			output.SetText(fmt.Sprintf("Erreur: %s", err.Error()))
		} else {
			output.SetText(result)
		}
	})

	copyBtn := widget.NewButton("Copier", func() {
		if result, _ := cp.viewModel.GetLastResult(); result != "" {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(result)
		}
	})

	// Mettre à jour le modèle quand les champs changent
	nameEntry.OnChanged = func(s string) {
		cp.viewModel.name = s
	}
	
	scriptEntry.OnChanged = func(s string) {
		cp.viewModel.script = s
	}

	topSection := container.NewVBox(
		widget.NewLabel("Configuration du Processeur:"),
		container.NewHBox(
			widget.NewLabel("Nom:"),
			nameEntry,
		),
		widget.NewLabel("Script JavaScript:"),
		scriptEntry,
		widget.NewLabel("Texte d'entrée:"),
		input,
		container.NewHBox(
			processBtn,
			copyBtn,
		),
		widget.NewLabel("Résultat:"),
	)

	return container.NewBorder(
		topSection,
		nil,
		nil,
		nil,
		container.NewVScroll(output),
	)
}

// CustomProcessorViewModel implémente ViewModel pour les processeurs personnalisés
type CustomProcessorViewModel struct {
	name       string
	script     string
	lastResult string
}

func NewCustomProcessorViewModel(name, script string) *CustomProcessorViewModel {
	return &CustomProcessorViewModel{
		name:   name,
		script: script,
	}
}

func (vm *CustomProcessorViewModel) Process(input string) (string, error) {
	if vm.script == "" {
		return "", fmt.Errorf("aucun script défini")
	}

	// Créer un runtime JavaScript
	jsRuntime := goja.New()
	
	// Définir la fonction d'entrée
	jsRuntime.Set("input", input)
	
	// Préparer le script avec une fonction wrapper si nécessaire
	script := vm.script
	if !strings.Contains(script, "function") && !strings.Contains(script, "=>") {
		// Si ce n'est pas une fonction, on l'enveloppe
		script = fmt.Sprintf("function process(input) { %s }", script)
	}
	
	// Ajouter une fonction process par défaut si elle n'existe pas
	if !strings.Contains(script, "process") {
		script = script + "\nfunction process(input) { return input; }"
	}
	
	// Exécuter le script
	_, err := jsRuntime.RunString(script)
	if err != nil {
		return "", fmt.Errorf("erreur dans le script: %v", err)
	}
	
	// Appeler la fonction process
	processFunc, ok := goja.AssertFunction(jsRuntime.Get("process"))
	if !ok {
		return "", fmt.Errorf("fonction 'process' non trouvée dans le script")
	}
	
	result, err := processFunc(goja.Undefined(), jsRuntime.ToValue(input))
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'exécution: %v", err)
	}
	
	resultStr := result.String()
	vm.lastResult = resultStr
	return resultStr, nil
}

func (vm *CustomProcessorViewModel) GetConfiguration() interface{} {
	return struct {
		Name   string
		Script string
	}{
		Name:   vm.name,
		Script: vm.script,
	}
}

func (vm *CustomProcessorViewModel) LoadConfiguration(config interface{}) error {
	cfg, ok := config.(struct{ Name, Script string })
	if !ok {
		return fmt.Errorf("configuration invalide")
	}
	vm.name = cfg.Name
	vm.script = cfg.Script
	return nil
}

func (vm *CustomProcessorViewModel) Validate() error {
	if vm.name == "" {
		return fmt.Errorf("le nom du processeur ne peut pas être vide")
	}
	if vm.script == "" {
		return fmt.Errorf("le script ne peut pas être vide")
	}
	return nil
}

func (vm *CustomProcessorViewModel) GetLastResult() (string, error) {
	return vm.lastResult, nil
}
