package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// EnhancedCodeEditor est un éditeur de code amélioré avec des fonctionnalités visuelles
type EnhancedCodeEditor struct {
	widget.BaseWidget
	entry        *widget.Entry
	lineNumbers  *widget.Label
	lexer        chroma.Lexer
	showLineNums bool
}

// NewEnhancedCodeEditor crée un nouvel éditeur de code amélioré
func NewEnhancedCodeEditor() *EnhancedCodeEditor {
	editor := &EnhancedCodeEditor{
		showLineNums: true,
	}
	
	// Créer l'entry principal avec police monospace
	editor.entry = widget.NewMultiLineEntry()
	editor.entry.TextStyle = fyne.TextStyle{Monospace: true}
	editor.entry.OnChanged = editor.onTextChanged
	editor.entry.Wrapping = fyne.TextWrapOff
	
	// Créer le label pour les numéros de ligne
	editor.lineNumbers = widget.NewLabel("1")
	editor.lineNumbers.TextStyle = fyne.TextStyle{Monospace: true}
	editor.lineNumbers.Alignment = fyne.TextAlignTrailing
	
	// Configuration du lexer JavaScript
	editor.lexer = lexers.Get("javascript")
	if editor.lexer == nil {
		editor.lexer = lexers.Fallback
	}
	editor.lexer = chroma.Coalesce(editor.lexer)
	
	editor.ExtendBaseWidget(editor)
	editor.updateLineNumbers()
	
	return editor
}

// SetText définit le texte de l'éditeur
func (e *EnhancedCodeEditor) SetText(text string) {
	e.entry.SetText(text)
	e.updateLineNumbers()
}

// Text retourne le texte actuel
func (e *EnhancedCodeEditor) Text() string {
	return e.entry.Text
}

// SetPlaceHolder définit le texte de placeholder
func (e *EnhancedCodeEditor) SetPlaceHolder(text string) {
	e.entry.SetPlaceHolder(text)
}

// EnableLineNumbers active/désactive les numéros de ligne
func (e *EnhancedCodeEditor) EnableLineNumbers(show bool) {
	e.showLineNums = show
	e.Refresh()
}

// onTextChanged est appelé quand le texte change
func (e *EnhancedCodeEditor) onTextChanged(text string) {
	e.updateLineNumbers()
}

// updateLineNumbers met à jour les numéros de ligne
func (e *EnhancedCodeEditor) updateLineNumbers() {
	if !e.showLineNums {
		return
	}
	
	text := e.entry.Text
	lines := strings.Split(text, "\n")
	lineCount := len(lines)
	
	// Créer les numéros de ligne
	var numbers strings.Builder
	for i := 1; i <= lineCount; i++ {
		if i > 1 {
			numbers.WriteString("\n")
		}
		// Formater le numéro avec un padding à droite
		if i < 10 {
			numbers.WriteString("  ")
		} else if i < 100 {
			numbers.WriteString(" ")
		}
		numbers.WriteString(fmt.Sprintf("%d", i))
	}
	
	e.lineNumbers.SetText(numbers.String())
}

// CreateRenderer crée le renderer pour ce widget
func (e *EnhancedCodeEditor) CreateRenderer() fyne.WidgetRenderer {
	if e.showLineNums {
		// Créer un container avec numéros de ligne à gauche
		lineNumContainer := container.NewVBox(e.lineNumbers)
		separator := widget.NewSeparator()
		
		content := container.NewBorder(nil, nil, 
			container.NewHBox(lineNumContainer, separator), 
			nil, 
			e.entry)
			
		return &enhancedCodeRenderer{
			entry:       e.entry,
			lineNumbers: e.lineNumbers,
			content:     content,
			objects:     []fyne.CanvasObject{content},
			parent:      e,
		}
	} else {
		// Mode sans numéros de ligne
		return &enhancedCodeRenderer{
			entry:   e.entry,
			content: e.entry,
			objects: []fyne.CanvasObject{e.entry},
			parent:  e,
		}
	}
}

// enhancedCodeRenderer est le renderer pour l'éditeur amélioré
type enhancedCodeRenderer struct {
	entry       *widget.Entry
	lineNumbers *widget.Label
	content     fyne.CanvasObject
	objects     []fyne.CanvasObject
	parent      *EnhancedCodeEditor
}

func (r *enhancedCodeRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *enhancedCodeRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *enhancedCodeRenderer) Refresh() {
	r.content.Refresh()
}

func (r *enhancedCodeRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *enhancedCodeRenderer) Destroy() {}

// Resize redimensionne le widget
func (e *EnhancedCodeEditor) Resize(size fyne.Size) {
	e.BaseWidget.Resize(size)
	if e.entry != nil {
		e.entry.Resize(size)
	}
}