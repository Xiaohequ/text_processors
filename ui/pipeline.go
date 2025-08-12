package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text_processors/ui/processors"
)

// ToolType représente le type d'outil dans le pipeline
type ToolType string

const (
	JSONFormatterTool   ToolType = "json_formatter"
	TextSplitterTool    ToolType = "text_splitter"
	TextJoinerTool      ToolType = "text_joiner"
	CustomProcessorTool ToolType = "custom_processor"
)

// ToolConfig interface commune pour toutes les configurations d'outils
type ToolConfig interface {
	GetType() ToolType
	Validate() error
	GetDisplayName() string
}

// JSONFormatterConfig configuration pour le formateur JSON
type JSONFormatterConfig struct {
	IndentType string // "2 espaces", "4 espaces", "Tabulations"
}

func (c JSONFormatterConfig) GetType() ToolType {
	return JSONFormatterTool
}

func (c JSONFormatterConfig) Validate() error {
	validTypes := []string{"2 espaces", "4 espaces", "Tabulations"}
	for _, valid := range validTypes {
		if c.IndentType == valid {
			return nil
		}
	}
	return fmt.Errorf("type d'indentation invalide: %s", c.IndentType)
}

func (c JSONFormatterConfig) GetDisplayName() string {
	return fmt.Sprintf("JSON Formatter (Indentation: %s)", c.IndentType)
}

// TextSplitterConfig configuration pour le diviseur de texte
type TextSplitterConfig struct {
	Delimiter string
}

func (c TextSplitterConfig) GetType() ToolType {
	return TextSplitterTool
}

func (c TextSplitterConfig) Validate() error {
	// Le délimiteur peut être vide (utilise \n par défaut)
	return nil
}

func (c TextSplitterConfig) GetDisplayName() string {
	delimiter := c.Delimiter
	if delimiter == "" {
		delimiter = "\\n (par défaut)"
	}
	return fmt.Sprintf("Text Splitter (Délimiteur: %s)", delimiter)
}

// TextJoinerConfig configuration pour le joigneur de texte
type TextJoinerConfig struct {
	Delimiter string
}

func (c TextJoinerConfig) GetType() ToolType {
	return TextJoinerTool
}

func (c TextJoinerConfig) Validate() error {
	// Le délimiteur peut être vide
	return nil
}

func (c TextJoinerConfig) GetDisplayName() string {
	delimiter := c.Delimiter
	if delimiter == "" {
		delimiter = "(vide)"
	}
	return fmt.Sprintf("Text Joiner (Délimiteur: %s)", delimiter)
}

// CustomProcessorConfig configuration pour le processeur personnalisé
type CustomProcessorConfig struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func (c CustomProcessorConfig) GetType() ToolType {
	return CustomProcessorTool
}

func (c CustomProcessorConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("le nom du processeur ne peut pas être vide")
	}
	if c.Script == "" {
		return fmt.Errorf("le script ne peut pas être vide")
	}
	return nil
}

func (c CustomProcessorConfig) GetDisplayName() string {
	return fmt.Sprintf("Custom Processor (%s)", c.Name)
}

// PipelineStep représente une étape dans le pipeline
type PipelineStep struct {
	ID        string               `json:"id"`
	Type      ToolType             `json:"type"`
	Config    interface{}          `json:"config"`
	Name      string               `json:"name"`
	Processor processors.Processor `json:"-"`
}

// Variable globale du pipeline actuel
var CurrentPipeline = &Pipeline{
	Name:  "Mon Pipeline",
	Steps: []PipelineStep{},
}

// Callbacks pour notifier les changements du pipeline
var pipelineUpdateCallbacks []func()

// RegisterPipelineUpdateCallback enregistre un callback à appeler quand le pipeline est mis à jour
func RegisterPipelineUpdateCallback(callback func()) {
	pipelineUpdateCallbacks = append(pipelineUpdateCallbacks, callback)
}

// ClearPipelineUpdateCallbacks efface tous les callbacks enregistrés
func ClearPipelineUpdateCallbacks() {
	pipelineUpdateCallbacks = nil
}

// NotifyPipelineUpdated notifie tous les callbacks enregistrés que le pipeline a été mis à jour
func NotifyPipelineUpdated() {
	for _, callback := range pipelineUpdateCallbacks {
		callback()
	}
}

// Pipeline représente une séquence d'outils configurés
type Pipeline struct {
	Steps []PipelineStep `json:"steps"`
	Name  string         `json:"name"`
}

type pipelineStepJSON struct {
	ID     string          `json:"id"`
	Type   ToolType        `json:"type"`
	Config json.RawMessage `json:"config"`
	Name   string          `json:"name"`
}

// SaveToFile sauvegarde le pipeline dans un fichier JSON
func (p *Pipeline) SaveToFile(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("échec de la sérialisation : %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("échec de l'écriture du fichier : %w", err)
	}
	return nil
}

// LoadFromFile charge un pipeline depuis un fichier JSON
func (p *Pipeline) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("échec de la lecture du fichier : %w", err)
	}

	var temp struct {
		Steps []pipelineStepJSON `json:"steps"`
		Name  string             `json:"name"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("échec du décodage JSON : %w", err)
	}

	p.Name = temp.Name
	p.Steps = make([]PipelineStep, len(temp.Steps))

	for i, step := range temp.Steps {
		var config ToolConfig
		var processor processors.Processor

		switch step.Type {
		case JSONFormatterTool:
			config = &JSONFormatterConfig{}
			processor = processors.NewJSONFormatterUI()
		case TextSplitterTool:
			config = &TextSplitterConfig{}
			processor = processors.NewTextSplitterUI()
		case TextJoinerTool:
			config = &TextJoinerConfig{}
			processor = processors.NewTextJoinerUI()
		case CustomProcessorTool:
			config = &CustomProcessorConfig{}
			processor = processors.NewCustomProcessor("", "") // Sera configuré après
		default:
			return fmt.Errorf("type d'outil inconnu: %s", step.Type)
		}

		if err := json.Unmarshal(step.Config, config); err != nil {
			return fmt.Errorf("erreur de configuration pour l'étape %d: %w", i+1, err)
		}

		// Convertir la configuration ToolConfig vers le format attendu par le ViewModel
		var vmConfig interface{}
		switch cfg := config.(type) {
		case *JSONFormatterConfig:
			vmConfig = struct{ IndentType string }{IndentType: cfg.IndentType}
		case *TextSplitterConfig:
			vmConfig = struct{ Delimiter string }{Delimiter: cfg.Delimiter}
		case *TextJoinerConfig:
			vmConfig = struct{ Delimiter string }{Delimiter: cfg.Delimiter}
		case *CustomProcessorConfig:
			vmConfig = struct{ Name, Script string }{Name: cfg.Name, Script: cfg.Script}
		}

		if err := processor.ViewModel().LoadConfiguration(vmConfig); err != nil {
			return fmt.Errorf("chargement configuration étape %d: %w", i+1, err)
		}

		p.Steps[i] = PipelineStep{
			ID:        step.ID,
			Type:      step.Type,
			Config:    config,
			Processor: processor,
			Name:      step.Name,
		}
	}

	return nil
}

// Validate valide la configuration complète du pipeline
func (p *Pipeline) Validate() error {
	if len(p.Steps) == 0 {
		return fmt.Errorf("le pipeline doit contenir au moins une étape")
	}

	for i, step := range p.Steps {
		if err := step.Processor.ViewModel().Validate(); err != nil {
			return fmt.Errorf("erreur à l'étape %d (%s): %w", i+1, step.Name, err)
		}
	}

	return nil
}

// GetDisplaySteps retourne une représentation textuelle des étapes
func (p *Pipeline) GetDisplaySteps() []string {
	var steps []string
	for i, step := range p.Steps {
		stepName := step.Name
		if stepName == "" {
			stepName = step.Processor.Name()
		}
		steps = append(steps, fmt.Sprintf("%d. %s", i+1, stepName))
	}
	return steps
}

// PipelineExecutor exécute un pipeline sur un texte d'entrée
type PipelineExecutor struct{}

// NewPipelineExecutor crée un nouvel exécuteur de pipeline
func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{}
}

// Execute exécute le pipeline sur le texte d'entrée
func (pe *PipelineExecutor) Execute(pipeline *Pipeline, input string) (string, error) {
	if err := pipeline.Validate(); err != nil {
		return "", fmt.Errorf("pipeline invalide: %w", err)
	}

	result := input

	for i, step := range pipeline.Steps {
		vm := step.Processor.ViewModel()
		var err error
		result, err = vm.Process(result)
		if err != nil {
			stepName := step.Name
			if stepName == "" {
				stepName = step.Processor.Name()
			}
			return "", fmt.Errorf("erreur à l'étape %d (%s): %w", i+1, stepName, err)
		}
	}

	return result, nil
}

// Fonctions de traitement pour chaque outil

// ProcessJSONFormatter traite le texte avec le formateur JSON
func ProcessJSONFormatter(input string, config ToolConfig) (string, error) {
	jsonConfig, ok := config.(JSONFormatterConfig)
	if !ok {
		return "", fmt.Errorf("configuration invalide pour JSON Formatter")
	}

	formatter := processors.NewFormatter(jsonConfig.IndentType)
	return formatter.FormatJSON(input)
}

// ProcessTextSplitter traite le texte avec le diviseur
func ProcessTextSplitter(input string, config ToolConfig) (string, error) {
	splitterConfig, ok := config.(TextSplitterConfig)
	if !ok {
		return "", fmt.Errorf("configuration invalide pour Text Splitter")
	}

	delimiter := splitterConfig.Delimiter
	if delimiter == "" {
		delimiter = "\n"
	}

	parts := strings.Split(input, delimiter)
	return strings.Join(parts, "\n"), nil
}

// ProcessTextJoiner traite le texte avec le joigneur
func ProcessTextJoiner(input string, config ToolConfig) (string, error) {
	joinerConfig, ok := config.(TextJoinerConfig)
	if !ok {
		return "", fmt.Errorf("configuration invalide pour Text Joiner")
	}

	lines := strings.Split(input, "\n")
	// Filtrer les lignes vides
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, strings.TrimSpace(line))
		}
	}

	return strings.Join(nonEmptyLines, joinerConfig.Delimiter), nil
}

// GetDefaultExecutor retourne un exécuteur de pipeline
func GetDefaultExecutor() *PipelineExecutor {
	return NewPipelineExecutor()
}
