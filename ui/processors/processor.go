package processors

import (
	"fyne.io/fyne/v2"
)

// Processor définit l'interface commune pour tous les processeurs
type Processor interface {
	Name() string
	Description() string
	CreateConfigurationUI() fyne.CanvasObject
	ViewModel() ViewModel
}

// ViewModel interface pour la logique métier des processeurs
type ViewModel interface {
	Process(input string) (string, error)
	GetConfiguration() interface{}
	LoadConfiguration(config interface{}) error
	Validate() error
}
