package ui

import (
	"fmt"
	"strings"
)

// ToolType représente le type d'outil dans le pipeline
type ToolType string

const (
	JSONFormatterTool ToolType = "json_formatter"
	TextSplitterTool  ToolType = "text_splitter"
	TextJoinerTool    ToolType = "text_joiner"
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

// PipelineStep représente une étape dans le pipeline
type PipelineStep struct {
	ID     string     // Identifiant unique de l'étape
	Config ToolConfig // Configuration de l'outil
	Name   string     // Nom personnalisé de l'étape
}

// Pipeline représente une séquence d'outils configurés
type Pipeline struct {
	Steps []PipelineStep
	Name  string
}

// Validate valide la configuration complète du pipeline
func (p *Pipeline) Validate() error {
	if len(p.Steps) == 0 {
		return fmt.Errorf("le pipeline doit contenir au moins une étape")
	}
	
	for i, step := range p.Steps {
		if err := step.Config.Validate(); err != nil {
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
			stepName = step.Config.GetDisplayName()
		}
		steps = append(steps, fmt.Sprintf("%d. %s", i+1, stepName))
	}
	return steps
}

// ProcessorFunc type de fonction pour traiter le texte
type ProcessorFunc func(input string, config ToolConfig) (string, error)

// PipelineExecutor exécute un pipeline sur un texte d'entrée
type PipelineExecutor struct {
	processors map[ToolType]ProcessorFunc
}

// NewPipelineExecutor crée un nouvel exécuteur de pipeline
func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{
		processors: make(map[ToolType]ProcessorFunc),
	}
}

// RegisterProcessor enregistre un processeur pour un type d'outil
func (pe *PipelineExecutor) RegisterProcessor(toolType ToolType, processor ProcessorFunc) {
	pe.processors[toolType] = processor
}

// Execute exécute le pipeline sur le texte d'entrée
func (pe *PipelineExecutor) Execute(pipeline *Pipeline, input string) (string, error) {
	if err := pipeline.Validate(); err != nil {
		return "", fmt.Errorf("pipeline invalide: %w", err)
	}
	
	result := input
	
	for i, step := range pipeline.Steps {
		processor, exists := pe.processors[step.Config.GetType()]
		if !exists {
			return "", fmt.Errorf("processeur non trouvé pour l'outil %s à l'étape %d", step.Config.GetType(), i+1)
		}
		
		var err error
		result, err = processor(result, step.Config)
		if err != nil {
			stepName := step.Name
			if stepName == "" {
				stepName = step.Config.GetDisplayName()
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
	
	formatter := NewFormatter(jsonConfig.IndentType)
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

// GetDefaultExecutor retourne un exécuteur avec tous les processeurs enregistrés
func GetDefaultExecutor() *PipelineExecutor {
	executor := NewPipelineExecutor()
	executor.RegisterProcessor(JSONFormatterTool, ProcessJSONFormatter)
	executor.RegisterProcessor(TextSplitterTool, ProcessTextSplitter)
	executor.RegisterProcessor(TextJoinerTool, ProcessTextJoiner)
	return executor
}
