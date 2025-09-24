package ui

import (
	"image/color"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// SyntaxHighlightedEntry est un widget d'entrée de texte avec coloration syntaxique
type SyntaxHighlightedEntry struct {
	widget.BaseWidget
	entry            *widget.Entry
	coloredContainer *fyne.Container      // Container pour les objets Canvas.Text colorés
	lexer            chroma.Lexer
	formatter        chroma.Formatter
	style            *chroma.Style
	placeholder      string
	wrapping         fyne.TextWrap
}

// NewSyntaxHighlightedEntry crée un nouveau widget avec coloration syntaxique JavaScript
func NewSyntaxHighlightedEntry() *SyntaxHighlightedEntry {
	syntaxEntry := &SyntaxHighlightedEntry{}
	
	// Créer un Entry normal directement
	normalEntry := widget.NewMultiLineEntry()
	normalEntry.Wrapping = fyne.TextWrapOff
	normalEntry.TextStyle = fyne.TextStyle{Monospace: true}
	
	// Stocker l'Entry normal directement
	syntaxEntry.entry = normalEntry
	
	// Créer un container pour les objets Canvas.Text colorés
	syntaxEntry.coloredContainer = container.NewWithoutLayout()
	
	// Configuration pour JavaScript
	syntaxEntry.lexer = lexers.Get("javascript")
	if syntaxEntry.lexer == nil {
		syntaxEntry.lexer = lexers.Fallback
	}
	syntaxEntry.lexer = chroma.Coalesce(syntaxEntry.lexer)
	
	// Style adapté à Fyne
	syntaxEntry.style = styles.Get("github")
	if syntaxEntry.style == nil {
		syntaxEntry.style = styles.Fallback
	}
	
	// Formatter personnalisé pour Fyne
	syntaxEntry.formatter = &fyneFormatter{}
	
	// Configurer les callbacks
	normalEntry.OnChanged = syntaxEntry.onTextChanged
	
	syntaxEntry.ExtendBaseWidget(syntaxEntry)
	return syntaxEntry
}

// SetText définit le texte de l'entrée
func (s *SyntaxHighlightedEntry) SetText(text string) {
	s.entry.SetText(text)
	s.updateHighlighting()
}

// Text retourne le texte actuel
func (s *SyntaxHighlightedEntry) Text() string {
	return s.entry.Text
}

// SetPlaceHolder définit le texte de placeholder
func (s *SyntaxHighlightedEntry) SetPlaceHolder(text string) {
	s.placeholder = text
	s.entry.SetPlaceHolder(text)
}


// Resize redimensionne le widget
func (s *SyntaxHighlightedEntry) Resize(size fyne.Size) {
	s.BaseWidget.Resize(size)
	if s.entry != nil {
		s.entry.Resize(size)
	}
}

// CreateRenderer crée le renderer pour ce widget
func (s *SyntaxHighlightedEntry) CreateRenderer() fyne.WidgetRenderer {
	s.entry.TextStyle = fyne.TextStyle{Monospace: true}
	s.updateHighlighting()
	
	// Stack : coloration en dessous + Entry au-dessus
	stack := container.NewStack(s.coloredContainer, s.entry)
	
	return &transparentTextSyntaxRenderer{
		stack:  stack,
		entry:  s.entry,
		parent: s,
	}
}

// onTextChanged est appelé quand le texte change
func (s *SyntaxHighlightedEntry) onTextChanged(text string) {
	s.updateHighlighting()
}

// updateHighlighting met à jour la coloration syntaxique avec Canvas.Text
func (s *SyntaxHighlightedEntry) updateHighlighting() {
	text := s.entry.Text
	if text == "" {
		// Vider le container des objets colorés
		s.coloredContainer.Objects = nil
		s.coloredContainer.Refresh()
		return
	}
	
	// Vider les objets existants
	s.coloredContainer.Objects = nil
	
	// Tokeniser le code JavaScript avec Chroma
	iterator, err := s.lexer.Tokenise(nil, text)
	if err != nil {
		// En cas d'erreur, créer un seul objet Text normal
		plainText := canvas.NewText(text, theme.ForegroundColor())
		plainText.TextStyle = s.entry.TextStyle
		plainText.Move(fyne.NewPos(0, 0))
		s.coloredContainer.Objects = append(s.coloredContainer.Objects, plainText)
		s.coloredContainer.Refresh()
		return
	}
	
	// Créer des objets Canvas.Text colorés pour chaque token
	var xOffset float32 = 0
	var yOffset float32 = 0
	fontSize := theme.TextSize()
	lineHeight := fontSize + 2 // Espacement entre les lignes
	
	for token := iterator(); token != chroma.EOF; token = iterator() {
		value := token.Value
		tokenType := token.Type
		
		if value == "" {
			continue
		}
		
		// Déterminer la couleur selon le type de token
		var color color.Color = theme.ForegroundColor()
		textStyle := s.entry.TextStyle
		
		switch tokenType {
		case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration,
			 chroma.KeywordNamespace, chroma.KeywordPseudo, chroma.KeywordReserved, chroma.KeywordType:
			color = theme.PrimaryColor()
			textStyle.Bold = true
		case chroma.String, chroma.StringDouble, chroma.StringSingle:
			color = theme.SuccessColor()
		case chroma.Comment, chroma.CommentSingle, chroma.CommentMultiline:
			color = theme.DisabledColor()
			textStyle.Italic = true
		case chroma.Number, chroma.NumberInteger, chroma.NumberFloat:
			color = theme.WarningColor()
		case chroma.Name, chroma.NameFunction:
			color = theme.ErrorColor()
		case chroma.Operator:
			color = theme.PrimaryColor()
		}
		
		// Traiter les retours à la ligne dans le token
		lines := strings.Split(value, "\n")
		for i, line := range lines {
			if line != "" {
				// Créer un objet Canvas.Text pour cette ligne
				textObj := canvas.NewText(line, color)
				textObj.TextStyle = textStyle
				textObj.Move(fyne.NewPos(xOffset, yOffset))
				s.coloredContainer.Objects = append(s.coloredContainer.Objects, textObj)
				
				// Avancer la position X
				textWidth := fyne.MeasureText(line, fontSize, textStyle).Width
				xOffset += textWidth
			}
			
			// Passer à la ligne suivante si nécessaire
			if i < len(lines)-1 {
				yOffset += lineHeight
				xOffset = 0
			}
		}
	}
	
	s.coloredContainer.Refresh()
}

// transparentTextSyntaxRenderer gère un Entry avec texte transparent + coloration
type transparentTextSyntaxRenderer struct {
	stack  *fyne.Container
	entry  *widget.Entry
	parent *SyntaxHighlightedEntry
}

func (r *transparentTextSyntaxRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
	// Synchroniser la mise à jour de la coloration quand la taille change
	r.parent.updateHighlighting()
}

func (r *transparentTextSyntaxRenderer) MinSize() fyne.Size {
	return r.stack.MinSize()
}

func (r *transparentTextSyntaxRenderer) Refresh() {
	r.stack.Refresh()
	r.parent.updateHighlighting()
}

func (r *transparentTextSyntaxRenderer) Objects() []fyne.CanvasObject {
	// Remettre la superposition avec Stack
	todoLabel := widget.NewLabel("TODO: Syntax highlighting")
	todoLabel.TextStyle = fyne.TextStyle{Monospace: true, Italic: true}
	
	// Stack : Entry en arrière-plan + label TODO au premier plan
	stack := container.NewStack(r.entry, todoLabel)
	
	return []fyne.CanvasObject{stack}
}

func (r *transparentTextSyntaxRenderer) Destroy() {
	// Rien à nettoyer
}

// simpleTransparentEntry encapsule un Entry et rend seulement son texte transparent
type simpleTransparentEntry struct {
	*widget.Entry
}

func (s *simpleTransparentEntry) CreateRenderer() fyne.WidgetRenderer {
	baseRenderer := s.Entry.CreateRenderer()
	return &simpleTransparentEntryRenderer{
		baseRenderer: baseRenderer,
		entry:        s.Entry,
	}
}

// simpleTransparentEntryRenderer rend le texte transparent sans boucles
type simpleTransparentEntryRenderer struct {
	baseRenderer fyne.WidgetRenderer
	entry        *widget.Entry
}

func (r *simpleTransparentEntryRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *simpleTransparentEntryRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *simpleTransparentEntryRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *simpleTransparentEntryRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var result []fyne.CanvasObject
	
	for _, obj := range objects {
		if text, ok := obj.(*canvas.Text); ok {
			// Rendre le texte plus visible pour l'édition
			visibleText := canvas.NewText(text.Text, color.RGBA{128, 128, 128, 150})
			visibleText.TextStyle = text.TextStyle
			visibleText.Alignment = text.Alignment
			visibleText.Move(text.Position())
			visibleText.Resize(text.Size())
			result = append(result, visibleText)
		} else {
			// Garder tous les autres éléments (curseur, sélection, bordures, etc.)
			result = append(result, obj)
		}
	}
	
	return result
}

func (r *simpleTransparentEntryRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// NewCodeEditor crée un éditeur de code avec coloration syntaxique
func NewCodeEditor() *SyntaxHighlightedEntry {
	editor := NewSyntaxHighlightedEntry()
	editor.entry.TextStyle = fyne.TextStyle{Monospace: true}
	
	return editor
}

// fyneFormatter est un formatter personnalisé pour Fyne
type fyneFormatter struct{}

func (f *fyneFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	// Implémentation simplifiée - pour le moment on ne fait rien
	return nil
}