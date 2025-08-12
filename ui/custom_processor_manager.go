package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"text_processors/ui/processors"
)

// CustomProcessorDefinition représente la définition d'un processeur personnalisé
type CustomProcessorDefinition struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

// computeConfDir retourne le chemin du dossier conf/ à côté de l'exécutable
func computeConfDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	confDir := filepath.Join(exeDir, "conf")
	return confDir, nil
}

// processorsConfDir retourne le sous-dossier conf/custom_processors
func processorsConfDir() (string, error) {
	confDir, err := computeConfDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(confDir, "custom_processors"), nil
}

// sanitizeFileName crée un nom de fichier sûr à partir d'un nom de processeur
func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	// remplacer tout ce qui n'est pas lettre, chiffre, tiret ou underscore par un tiret
	re := regexp.MustCompile("[^a-z0-9-_]+")

	name = re.ReplaceAllString(name, "-")
	if name == "" {
		name = "custom"
	}
	return name
}

// SaveAll sauvegarde chaque processeur dans conf/custom_processors/<slug>.json
func (cpm *CustomProcessorManager) SaveAll() error {
	dir, err := processorsConfDir()
	if err != nil {
		return fmt.Errorf("impossible de déterminer le dossier: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("impossible de créer le dossier: %w", err)
	}
	// Purge simple: on peut laisser les fichiers orphelins, ou nettoyer (optionnel). Ici on laisse.
	for _, def := range cpm.definitions {
		file := filepath.Join(dir, sanitizeFileName(def.Name)+".json")
		data, err := json.MarshalIndent(def, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(file, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

// LoadAll charge depuis conf/custom_processors/ (ou migre depuis l'ancien fichier unique)
func (cpm *CustomProcessorManager) LoadAll() error {
	var defs []CustomProcessorDefinition
	// 1) tenter nouveau format dossier
	dir, err := processorsConfDir()
	if err == nil {
		if entries, err2 := os.ReadDir(dir); err2 == nil {
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if !strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
					continue
				}
				b, err3 := os.ReadFile(filepath.Join(dir, e.Name()))
				if err3 != nil {
					continue
				}
				var def CustomProcessorDefinition
				if json.Unmarshal(b, &def) == nil && def.Name != "" && def.Script != "" {
					defs = append(defs, def)
				}
			}
		}
	}
	// 2) fallback ancien fichier unique conf/custom_processors.json
	if len(defs) == 0 {
		confDir, err4 := computeConfDir()
		if err4 == nil {
			legacy := filepath.Join(confDir, "custom_processors.json")
			if b, err5 := os.ReadFile(legacy); err5 == nil {
				var payload struct {
					Processors []CustomProcessorDefinition `json:"processors"`
				}
				if json.Unmarshal(b, &payload) == nil {
					defs = payload.Processors
					// migrer en nouveau format
					cpm.definitions = defs
					_ = cpm.SaveAll()
				}
			}
		}
	}
	cpm.definitions = defs
	cpm.fireUpdate()
	return nil
}

// CustomProcessorManager gère les processeurs personnalisés
type CustomProcessorManager struct {
	definitions     []CustomProcessorDefinition
	onUpdate        func()   // Callback principal (compat)
	updateCallbacks []func() // Multiples callbacks
}

var GlobalCustomProcessorManager = &CustomProcessorManager{
	definitions: []CustomProcessorDefinition{},
}

func (cpm *CustomProcessorManager) SetUpdateCallback(callback func()) {
	cpm.onUpdate = callback
}

// RegisterUpdateCallback ajoute un callback de mise à jour
func (cpm *CustomProcessorManager) RegisterUpdateCallback(cb func()) {
	cpm.updateCallbacks = append(cpm.updateCallbacks, cb)
}

func (cpm *CustomProcessorManager) fireUpdate() {
	if cpm.onUpdate != nil {
		cpm.onUpdate()
	}
	for _, cb := range cpm.updateCallbacks {
		if cb != nil {
			cb()
		}
	}
}

func (cpm *CustomProcessorManager) AddProcessor(name, script string) {
	cpm.definitions = append(cpm.definitions, CustomProcessorDefinition{
		Name:   name,
		Script: script,
	})
	// Sauvegarde immédiate
	if err := cpm.SaveAll(); err != nil {
		fmt.Printf("[WARN] échec de la sauvegarde des processeurs personnalisés: %v\n", err)
	}
	cpm.fireUpdate()
}

func (cpm *CustomProcessorManager) RemoveProcessor(index int) {
	if index >= 0 && index < len(cpm.definitions) {
		cpm.definitions = append(cpm.definitions[:index], cpm.definitions[index+1:]...)
		// Sauvegarde après suppression
		if err := cpm.SaveAll(); err != nil {
			fmt.Printf("[WARN] échec de la sauvegarde des processeurs personnalisés: %v\n", err)
		}
		cpm.fireUpdate()
	}
}

func (cpm *CustomProcessorManager) GetProcessors() []CustomProcessorDefinition {
	return cpm.definitions
}

func (cpm *CustomProcessorManager) CreateProcessor(name, script string) processors.Processor {
	return processors.NewCustomProcessor(name, script)
}

// CreateAddCustomProcessorDialog crée la boîte de dialogue pour ajouter un processeur personnalisé
func CreateAddCustomProcessorDialog(parent fyne.Window) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom du processeur (ex: Convertisseur Majuscules)")

	scriptEntry := widget.NewMultiLineEntry()
	scriptEntry.SetPlaceHolder(`Script JavaScript:
function process(input) {
    // Votre code ici
    return input.toUpperCase();
}

Ou simplement:
return input.toUpperCase();`)
	scriptEntry.Wrapping = fyne.TextWrapWord
	scriptEntry.Resize(fyne.NewSize(500, 200))

	// Exemples prédéfinis
	exampleSelect := widget.NewSelect([]string{
		"Convertir en majuscules",
		"Convertir en minuscules",
		"Inverser le texte",
		"Compter les mots",
		"Supprimer les espaces",
		"Ajouter des numéros de ligne",
	}, func(selected string) {
		var script string
		switch selected {
		case "Convertir en majuscules":
			nameEntry.SetText("Convertisseur Majuscules")
			script = "return input.toUpperCase();"
		case "Convertir en minuscules":
			nameEntry.SetText("Convertisseur Minuscules")
			script = "return input.toLowerCase();"
		case "Inverser le texte":
			nameEntry.SetText("Inverseur de Texte")
			script = "return input.split('').reverse().join('');"
		case "Compter les mots":
			nameEntry.SetText("Compteur de Mots")
			script = "return 'Nombre de mots: ' + input.split(/\\s+/).filter(word => word.length > 0).length;"
		case "Supprimer les espaces":
			nameEntry.SetText("Suppresseur d'Espaces")
			script = "return input.replace(/\\s+/g, '');"
		case "Ajouter des numéros de ligne":
			nameEntry.SetText("Numéroteur de Lignes")
			script = `var lines = input.split('\\n');
return lines.map(function(line, index) {
    return (index + 1) + '. ' + line;
}).join('\\n');`
		}
		scriptEntry.SetText(script)
	})

	// Zone de test
	testInput := widget.NewEntry()
	testInput.SetPlaceHolder("Texte de test...")

	testOutput := widget.NewMultiLineEntry()
	testOutput.Disable()
	testOutput.Resize(fyne.NewSize(0, 60))

	testBtn := widget.NewButton("Tester", func() {
		if nameEntry.Text == "" || scriptEntry.Text == "" {
			testOutput.SetText("Erreur: Nom et script requis")
			return
		}

		// Créer un processeur temporaire pour tester
		tempProcessor := processors.NewCustomProcessor(nameEntry.Text, scriptEntry.Text)
		vm := tempProcessor.ViewModel()

		// Charger la configuration
		config := struct{ Name, Script string }{
			Name:   nameEntry.Text,
			Script: scriptEntry.Text,
		}
		err := vm.LoadConfiguration(config)
		if err != nil {
			testOutput.SetText(fmt.Sprintf("Erreur de configuration: %v", err))
			return
		}

		// Tester avec l'entrée
		result, err := vm.Process(testInput.Text)
		if err != nil {
			testOutput.SetText(fmt.Sprintf("Erreur: %v", err))
		} else {
			testOutput.SetText(result)
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Créer un Processeur Personnalisé"),
		widget.NewSeparator(),

		widget.NewLabel("Nom du processeur:"),
		nameEntry,

		widget.NewLabel("Exemples prédéfinis:"),
		exampleSelect,

		widget.NewLabel("Script JavaScript:"),
		container.NewScroll(scriptEntry),

		widget.NewSeparator(),
		widget.NewLabel("Test du processeur:"),
		container.NewHBox(
			widget.NewLabel("Entrée:"),
			testInput,
			testBtn,
		),
		widget.NewLabel("Sortie:"),
		testOutput,
	)

	// Créer les boutons
	addBtn := widget.NewButton("Ajouter", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("le nom du processeur est requis"), parent)
			return
		}
		if scriptEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("le script est requis"), parent)
			return
		}

		// Valider le script en créant un processeur temporaire
		tempProcessor := processors.NewCustomProcessor(nameEntry.Text, scriptEntry.Text)
		vm := tempProcessor.ViewModel()
		config := struct{ Name, Script string }{
			Name:   nameEntry.Text,
			Script: scriptEntry.Text,
		}
		err := vm.LoadConfiguration(config)
		if err != nil {
			dialog.ShowError(fmt.Errorf("erreur de configuration: %v", err), parent)
			return
		}

		err = vm.Validate()
		if err != nil {
			dialog.ShowError(fmt.Errorf("validation échouée: %v", err), parent)
			return
		}

		// Ajouter le processeur
		GlobalCustomProcessorManager.AddProcessor(nameEntry.Text, scriptEntry.Text)

		// Afficher confirmation
		dialog.ShowInformation("Succès",
			fmt.Sprintf("Processeur '%s' ajouté avec succès!", nameEntry.Text),
			parent)
	})

	cancelBtn := widget.NewButton("Annuler", func() {
		// Fermer la fenêtre (sera géré par la boîte de dialogue)
	})

	buttons := container.NewHBox(addBtn, cancelBtn)
	finalContent := container.NewBorder(nil, buttons, nil, nil, content)

	// Créer et afficher la boîte de dialogue
	customDialog := dialog.NewCustom("Ajouter un Processeur Personnalisé", "Fermer", finalContent, parent)
	customDialog.Resize(fyne.NewSize(600, 700))
	customDialog.Show()
}
