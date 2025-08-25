package ui

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// SyntaxHighlightedEntry est un widget d'entrée de texte avec coloration syntaxique
type SyntaxHighlightedEntry struct {
	widget.BaseWidget
	entry        *widget.Entry
	richText     *widget.RichText
	lexer        chroma.Lexer
	formatter    chroma.Formatter
	style        *chroma.Style
	placeholder  string
	wrapping     fyne.TextWrap
	showHighlights bool
}

// NewSyntaxHighlightedEntry crée un nouveau widget avec coloration syntaxique JavaScript
func NewSyntaxHighlightedEntry() *SyntaxHighlightedEntry {
	entry := &SyntaxHighlightedEntry{}
	entry.entry = widget.NewMultiLineEntry()
	entry.richText = widget.NewRichText()
	entry.showHighlights = false // Désactivé par défaut pour compatibilité
	
	// Configuration pour JavaScript
	entry.lexer = lexers.Get("javascript")
	if entry.lexer == nil {
		entry.lexer = lexers.Fallback
	}
	entry.lexer = chroma.Coalesce(entry.lexer)
	
	// Style adapté à Fyne
	entry.style = styles.Get("github")
	if entry.style == nil {
		entry.style = styles.Fallback
	}
	
	// Formatter personnalisé pour Fyne
	entry.formatter = &fyneFormatter{}
	
	entry.entry.OnChanged = entry.onTextChanged
	entry.entry.Wrapping = fyne.TextWrapWord
	
	entry.ExtendBaseWidget(entry)
	return entry
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

// EnableSyntaxHighlighting active la coloration syntaxique
func (s *SyntaxHighlightedEntry) EnableSyntaxHighlighting() {
	s.showHighlights = true
	s.updateHighlighting()
	s.Refresh()
}

// DisableSyntaxHighlighting désactive la coloration syntaxique  
func (s *SyntaxHighlightedEntry) DisableSyntaxHighlighting() {
	s.showHighlights = false
	s.Refresh()
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
	// Utiliser un container Stack avec RichText en arrière-plan et Entry transparent par-dessus
	if s.showHighlights {
		// Rendre l'entry transparent pour voir le RichText coloré en dessous
		s.entry.TextStyle = fyne.TextStyle{Monospace: true}
		container := container.NewStack(s.richText, s.entry)
		return &syntaxHighlightedRenderer{
			entry:     s.entry,
			richText:  s.richText,
			container: container,
			objects:   []fyne.CanvasObject{container},
			parent:    s,
		}
	} else {
		// Mode normal sans coloration
		return &syntaxHighlightedRenderer{
			entry:     s.entry,
			richText:  s.richText,
			container: nil,
			objects:   []fyne.CanvasObject{s.entry},
			parent:    s,
		}
	}
}

// onTextChanged est appelé quand le texte change
func (s *SyntaxHighlightedEntry) onTextChanged(text string) {
	s.updateHighlighting()
}

// updateHighlighting met à jour la coloration syntaxique
func (s *SyntaxHighlightedEntry) updateHighlighting() {
	if !s.showHighlights {
		return
	}
	
	text := s.entry.Text
	if text == "" {
		s.richText.Segments = nil
		s.richText.Refresh()
		return
	}
	
	// Tokenisation du code JavaScript
	iterator, err := s.lexer.Tokenise(nil, text)
	if err != nil {
		return
	}
	
	// Construction des segments avec couleurs réelles
	var segments []widget.RichTextSegment
	
	for token := iterator(); token != chroma.EOF; token = iterator() {
		value := token.Value
		tokenType := token.Type
		
		// Créer un segment avec style approprié
		segment := &widget.TextSegment{
			Text: value,
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		}
		
		// Appliquer des couleurs selon le type de token
		switch tokenType {
		case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration, chroma.KeywordNamespace, chroma.KeywordPseudo, chroma.KeywordReserved, chroma.KeywordType:
			segment.Style.ColorName = "primary" // Bleu pour les mots-clés
			segment.Style.TextStyle.Bold = true
		case chroma.String, chroma.StringDouble, chroma.StringSingle:
			segment.Style.ColorName = "success" // Vert pour les chaînes
		case chroma.Comment, chroma.CommentSingle, chroma.CommentMultiline:
			segment.Style.ColorName = "disabled" // Gris pour les commentaires
			segment.Style.TextStyle.Italic = true
		case chroma.Number, chroma.NumberInteger, chroma.NumberFloat:
			segment.Style.ColorName = "warning" // Orange pour les nombres
		case chroma.Punctuation:
			segment.Style.ColorName = "foreground" // Couleur normale pour la ponctuation
		default:
			segment.Style.ColorName = "foreground" // Couleur normale
		}
		
		segments = append(segments, segment)
	}
	
	s.richText.Segments = segments
	s.richText.Refresh()
}

// syntaxHighlightedRenderer est le renderer pour le widget
type syntaxHighlightedRenderer struct {
	entry     *widget.Entry
	richText  *widget.RichText
	container *fyne.Container
	objects   []fyne.CanvasObject
	parent    *SyntaxHighlightedEntry
}

func (r *syntaxHighlightedRenderer) Layout(size fyne.Size) {
	if r.container != nil {
		r.container.Resize(size)
	} else {
		r.entry.Resize(size)
	}
}

func (r *syntaxHighlightedRenderer) MinSize() fyne.Size {
	return r.entry.MinSize()
}

func (r *syntaxHighlightedRenderer) Refresh() {
	r.entry.Refresh()
}

func (r *syntaxHighlightedRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *syntaxHighlightedRenderer) Destroy() {}

// fyneFormatter est un formatter personnalisé pour Fyne
type fyneFormatter struct{}

func (f *fyneFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	// Implémentation simplifiée - pour le moment on ne fait rien
	return nil
}

// NewCodeEditor crée un éditeur de code avec police monospace
func NewCodeEditor() *SyntaxHighlightedEntry {
	editor := NewSyntaxHighlightedEntry()
	editor.entry.TextStyle = fyne.TextStyle{Monospace: true}
	// Activer la coloration syntaxique
	editor.EnableSyntaxHighlighting() 
	return editor
}