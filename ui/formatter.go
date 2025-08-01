package ui

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Formatter struct {
	IndentType  string
	IndentSize  int
}

func NewFormatter(indentType string) *Formatter {
	f := &Formatter{
		IndentType: indentType,
	}
	
	switch indentType {
	case "2 espaces":
		f.IndentSize = 2
	case "4 espaces":
		f.IndentSize = 4
	default: // Tabulations
		f.IndentSize = 1
	}
	
	return f
}

func (f *Formatter) FormatJSON(input string) (string, error) {
	// Valider d'abord le JSON
	if err := ValidateJSON(input); err != nil {
		return "", err
	}
	
	// Parser le JSON
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}
	
	// Formater avec l'indentation spécifiée
	var indent string
	if f.IndentType == "Tabulations" {
		indent = "\t"
	} else {
		indent = strings.Repeat(" ", f.IndentSize)
	}
	
	var formatted bytes.Buffer
	encoder := json.NewEncoder(&formatted)
	encoder.SetIndent("", indent)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}
	
	return formatted.String(), nil
}
