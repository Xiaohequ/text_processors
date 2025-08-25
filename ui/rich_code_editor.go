package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// RichCodeEditor est un éditeur de code avec syntax highlighting basé sur RichText
type RichCodeEditor struct {
	widget.BaseWidget
	entry       *widget.Entry
	richText    *widget.RichText
	lexer       chroma.Lexer
	container   *fyne.Container
	showHighlight bool
}

// NewRichCodeEditor crée un nouvel éditeur avec RichText coloré
func NewRichCodeEditor() *RichCodeEditor {
	editor := &RichCodeEditor{
		showHighlight: true,
	}
	
	// Configuration du lexer JavaScript
	editor.lexer = lexers.Get("javascript")
	if editor.lexer == nil {
		editor.lexer = lexers.Fallback
	}
	editor.lexer = chroma.Coalesce(editor.lexer)
	
	// Créer l'Entry pour l'édition (invisible en mode highlight)
	editor.entry = widget.NewMultiLineEntry()
	editor.entry.TextStyle = fyne.TextStyle{Monospace: true}
	editor.entry.Wrapping = fyne.TextWrapOff
	
	// Créer le RichText pour l'affichage coloré
	editor.richText = widget.NewRichText()
	editor.richText.Wrapping = fyne.TextWrapOff
	
	// Synchronisation Entry -> RichText
	editor.entry.OnChanged = func(text string) {
		if editor.showHighlight {
			editor.updateSyntaxHighlighting(text)
		}
	}
	
	// Mode par défaut : affichage RichText uniquement
	editor.container = container.NewStack(editor.richText)
	
	editor.ExtendBaseWidget(editor)
	
	return editor
}

// SetText définit le texte
func (r *RichCodeEditor) SetText(text string) {
	r.entry.SetText(text)
	if r.showHighlight {
		r.updateSyntaxHighlighting(text)
	}
}

// GetText retourne le texte actuel
func (r *RichCodeEditor) GetText() string {
	return r.entry.Text
}

// SetPlaceHolder définit le placeholder
func (r *RichCodeEditor) SetPlaceHolder(text string) {
	r.entry.SetPlaceHolder(text)
	if r.GetText() == "" {
		// Afficher le placeholder en gris dans RichText
		r.richText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameDisabled,
					TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
				},
			},
		}
		r.richText.Refresh()
	}
}

// EnableEditing active le mode édition (Entry visible)
func (r *RichCodeEditor) EnableEditing() {
	r.container.Objects = []fyne.CanvasObject{r.entry}
	r.showHighlight = false
	r.container.Refresh()
}

// EnableHighlighting active le mode affichage coloré (RichText visible)
func (r *RichCodeEditor) EnableHighlighting() {
	r.showHighlight = true
	r.updateSyntaxHighlighting(r.entry.Text)
	r.container.Objects = []fyne.CanvasObject{r.richText}
	r.container.Refresh()
}

// ToggleMode bascule entre édition et affichage coloré
func (r *RichCodeEditor) ToggleMode() {
	if r.showHighlight {
		r.EnableEditing()
	} else {
		r.EnableHighlighting()
	}
}

// OnChanged définit le callback
func (r *RichCodeEditor) OnChanged(callback func(string)) {
	originalCallback := r.entry.OnChanged
	r.entry.OnChanged = func(text string) {
		if originalCallback != nil {
			originalCallback(text)
		}
		if callback != nil {
			callback(text)
		}
	}
}

// updateSyntaxHighlighting met à jour la coloration RichText
func (r *RichCodeEditor) updateSyntaxHighlighting(text string) {
	if !r.showHighlight {
		return
	}
	
	if text == "" {
		r.richText.Segments = nil
		r.richText.Refresh()
		return
	}
	
	// Tokenisation avec Chroma
	iterator, err := r.lexer.Tokenise(nil, text)
	if err != nil {
		// En cas d'erreur, afficher le texte normal
		r.richText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameForeground,
					TextStyle: fyne.TextStyle{Monospace: true},
				},
			},
		}
		r.richText.Refresh()
		return
	}
	
	// Construction des segments avec couleurs
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
		case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration,
			chroma.KeywordNamespace, chroma.KeywordPseudo, chroma.KeywordReserved, chroma.KeywordType:
			segment.Style.ColorName = theme.ColorNamePrimary
			segment.Style.TextStyle.Bold = true
		case chroma.String, chroma.StringDouble, chroma.StringSingle:
			segment.Style.ColorName = theme.ColorNameSuccess
		case chroma.Comment, chroma.CommentSingle, chroma.CommentMultiline:
			segment.Style.ColorName = theme.ColorNameDisabled
			segment.Style.TextStyle.Italic = true
		case chroma.Number, chroma.NumberInteger, chroma.NumberFloat:
			segment.Style.ColorName = theme.ColorNameWarning
		case chroma.Name, chroma.NameFunction:
			// Utiliser la couleur d'erreur pour les fonctions (souvent rouge/violet)
			segment.Style.ColorName = theme.ColorNameError
		case chroma.Operator:
			// Utiliser la couleur principale pour les opérateurs
			segment.Style.ColorName = theme.ColorNamePrimary
		default:
			segment.Style.ColorName = theme.ColorNameForeground
		}
		
		segments = append(segments, segment)
	}
	
	r.richText.Segments = segments
	r.richText.Refresh()
}

// CreateRenderer crée le renderer
func (r *RichCodeEditor) CreateRenderer() fyne.WidgetRenderer {
	return &richEditorRenderer{
		editor:    r,
		container: r.container,
		objects:   []fyne.CanvasObject{r.container},
	}
}

// richEditorRenderer est le renderer
type richEditorRenderer struct {
	editor    *RichCodeEditor
	container *fyne.Container
	objects   []fyne.CanvasObject
}

func (r *richEditorRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *richEditorRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *richEditorRenderer) Refresh() {
	r.container.Refresh()
}

func (r *richEditorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *richEditorRenderer) Destroy() {}

// Resize redimensionne l'éditeur
func (r *RichCodeEditor) Resize(size fyne.Size) {
	r.BaseWidget.Resize(size)
	r.container.Resize(size)
}

// Tapped gère les clics pour passer en mode édition
func (r *RichCodeEditor) Tapped(*fyne.PointEvent) {
	if r.showHighlight {
		r.EnableEditing()
	}
}