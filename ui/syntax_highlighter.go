package ui

import (
	"fmt"
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
	entry        *transparentMultiLineEntry
	richText     *widget.RichText
	coloredContainer *fyne.Container  // Container pour les objets Canvas.Text colorés
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
	
	// Créer un Entry normal
	normalEntry := widget.NewMultiLineEntry()
	normalEntry.Wrapping = fyne.TextWrapOff
	normalEntry.TextStyle = fyne.TextStyle{Monospace: true}
	
	// Créer un thème transparent modifié
	hiddenTheme := &transparentTheme{}
	
	// Utiliser un wrapper qui appliquera le thème transparent
	entry.entry = &transparentMultiLineEntry{
		Entry:            normalEntry,
		transparentTheme: hiddenTheme,
	}
	
	// Créer un container pour les objets Canvas.Text colorés
	entry.coloredContainer = container.NewWithoutLayout()
	entry.richText = widget.NewRichText() // Garder pour compatibilité avec l'ancien renderer
	
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
	
	// Configurer les callbacks
	normalEntry.OnChanged = entry.onTextChanged
	
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
	return s.entry.Entry.Text
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
		s.entry.Entry.Resize(size)
	}
}

// CreateRenderer crée le renderer pour ce widget
func (s *SyntaxHighlightedEntry) CreateRenderer() fyne.WidgetRenderer {
	if s.showHighlights {
		s.entry.Entry.TextStyle = fyne.TextStyle{Monospace: true}
		s.updateHighlighting()
		
		// Utiliser l'Entry avec texte caché
		// Le texte coloré sera visible par-dessus
		stack := container.NewStack(s.entry.Entry, s.coloredContainer)
		
		return &simpleSyntaxRenderer{
			stack:  stack,
			entry:  s.entry.Entry,
			parent: s,
		}
	} else {
		// Mode normal sans coloration
		return &syntaxHighlightedRenderer{
			entry:     s.entry,
			richText:  s.richText,
			container: nil,
			objects:   []fyne.CanvasObject{s.entry.Entry},
			parent:    s,
		}
	}
}

// onTextChanged est appelé quand le texte change
func (s *SyntaxHighlightedEntry) onTextChanged(text string) {
	s.updateHighlighting()
}

// updateHighlighting met à jour la coloration syntaxique avec Canvas.Text
func (s *SyntaxHighlightedEntry) updateHighlighting() {
	if !s.showHighlights {
		return
	}
	
	text := s.entry.Entry.Text
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
		plainText.TextStyle = s.entry.Entry.TextStyle
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
		textStyle := s.entry.Entry.TextStyle
		
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
	fmt.Printf("DEBUG: Created %d Canvas.Text objects for syntax highlighting\n", len(s.coloredContainer.Objects))
}

// syntaxHighlightedRenderer est le renderer pour le widget
type syntaxHighlightedRenderer struct {
	entry     *transparentMultiLineEntry
	richText  *widget.RichText
	container *fyne.Container
	objects   []fyne.CanvasObject
	parent    *SyntaxHighlightedEntry
}

func (r *syntaxHighlightedRenderer) Layout(size fyne.Size) {
	if r.container != nil {
		r.container.Resize(size)
		// S'assurer que le RichText a exactement la même position que l'Entry
		entryPos := r.entry.Position()
		entrySize := r.entry.Size()
		r.richText.Move(entryPos)
		r.richText.Resize(entrySize)
	} else {
		r.entry.Entry.Resize(size)
	}
}

func (r *syntaxHighlightedRenderer) MinSize() fyne.Size {
	if r.container != nil {
		return r.container.MinSize()
	}
	return r.entry.Entry.MinSize()
}

func (r *syntaxHighlightedRenderer) Refresh() {
	if r.container != nil {
		r.container.Refresh()
	} else {
		r.entry.Entry.Refresh()
	}
}

func (r *syntaxHighlightedRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *syntaxHighlightedRenderer) Destroy() {}

// simpleSyntaxRenderer utilise un Stack simple avec Entry transparent
type simpleSyntaxRenderer struct {
	stack  *fyne.Container
	entry  *widget.Entry
	parent *SyntaxHighlightedEntry
}

func (r *simpleSyntaxRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *simpleSyntaxRenderer) MinSize() fyne.Size {
	return r.stack.MinSize()
}

func (r *simpleSyntaxRenderer) Refresh() {
	r.stack.Refresh()
}

func (r *simpleSyntaxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *simpleSyntaxRenderer) Destroy() {
	// Rien à nettoyer
}

// fyneFormatter est un formatter personnalisé pour Fyne
type fyneFormatter struct{}

func (f *fyneFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	// Implémentation simplifiée - pour le moment on ne fait rien
	return nil
}

// NewCodeEditor crée un éditeur de code avec police monospace
func NewCodeEditor() *SyntaxHighlightedEntry {
	editor := NewSyntaxHighlightedEntry()
	editor.entry.Entry.TextStyle = fyne.TextStyle{Monospace: true}
	// Activer la coloration syntaxique immédiatement
	editor.showHighlights = true
	editor.EnableSyntaxHighlighting()
	
	return editor
}


// transparentMultiLineEntry est un Entry multiligne avec texte transparent
type transparentMultiLineEntry struct {
	*widget.Entry
	transparentTheme *transparentTheme
}

func newTransparentMultiLineEntry() *transparentMultiLineEntry {
	entry := &transparentMultiLineEntry{
		Entry:            widget.NewMultiLineEntry(),
		transparentTheme: &transparentTheme{},
	}
	return entry
}

func (t *transparentMultiLineEntry) CreateRenderer() fyne.WidgetRenderer {
	// Appliquer temporairement le thème transparent
	originalTheme := fyne.CurrentApp().Settings().Theme()
	fyne.CurrentApp().Settings().SetTheme(t.transparentTheme)
	
	// Créer le renderer
	renderer := t.Entry.CreateRenderer()
	
	// Restaurer le thème original dans une goroutine pour éviter les boucles
	go func() {
		fyne.CurrentApp().Settings().SetTheme(originalTheme)
	}()
	
	return renderer
}

// simpleTransparentRenderer rend seulement les rectangles transparents
type simpleTransparentRenderer struct {
	baseRenderer fyne.WidgetRenderer
	entry        *widget.Entry
}

func (r *simpleTransparentRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *simpleTransparentRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *simpleTransparentRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *simpleTransparentRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var transparentObjects []fyne.CanvasObject
	
	for i, obj := range objects {
		fmt.Printf("DEBUG: Simple Entry object %d: %T\n", i, obj)
		if text, ok := obj.(*canvas.Text); ok {
			fmt.Printf("DEBUG: Simple - Remplacement texte '%s' par transparent\n", text.Text)
			// Utiliser une couleur vraiment transparente pour le TEXTE SEULEMENT
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			transparentObjects = append(transparentObjects, transparentText)
		} else {
			// Garder TOUS les autres objets intacts pour préserver l'interaction
			// (rectangles, scrolls, curseurs, etc.)
			transparentObjects = append(transparentObjects, obj)
		}
	}
	
	return transparentObjects
}

func (r *simpleTransparentRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// transparentEntryWrapperRenderer encapsule le renderer de base
type transparentEntryWrapperRenderer struct {
	baseRenderer fyne.WidgetRenderer
	entry        *widget.Entry
}

func (r *transparentEntryWrapperRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *transparentEntryWrapperRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *transparentEntryWrapperRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *transparentEntryWrapperRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var transparentObjects []fyne.CanvasObject
	
	for i, obj := range objects {
		fmt.Printf("DEBUG: Entry object %d: %T\n", i, obj)
		if text, ok := obj.(*canvas.Text); ok {
			fmt.Printf("DEBUG: Remplacement texte '%s' par transparent\n", text.Text)
			// Remplacer par un texte transparent
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			transparentObjects = append(transparentObjects, transparentText)
		} else if rect, ok := obj.(*canvas.Rectangle); ok {
			fmt.Printf("DEBUG: Remplacement Rectangle background par transparent\n")
			// Rendre le background transparent
			transparentRect := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
			transparentRect.Move(rect.Position())
			transparentRect.Resize(rect.Size())
			transparentObjects = append(transparentObjects, transparentRect)
		} else if scroll, ok := obj.(*container.Scroll); ok {
			fmt.Printf("DEBUG: Trouvé Scroll, création d'un wrapper transparent\n")
			// Créer un wrapper pour le Scroll qui rend son contenu transparent
			transparentScroll := &transparentScrollWrapper{Scroll: scroll}
			transparentObjects = append(transparentObjects, transparentScroll)
		} else {
			transparentObjects = append(transparentObjects, obj)
		}
	}
	
	return transparentObjects
}

func (r *transparentEntryWrapperRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// transparentScrollWrapper encapsule un Scroll et rend son contenu transparent
type transparentScrollWrapper struct {
	*container.Scroll
}

func (t *transparentScrollWrapper) CreateRenderer() fyne.WidgetRenderer {
	baseRenderer := t.Scroll.CreateRenderer()
	return &transparentScrollRenderer{
		baseRenderer: baseRenderer,
		scroll:       t.Scroll,
	}
}

// transparentScrollRenderer rend le Scroll avec contenu transparent
type transparentScrollRenderer struct {
	baseRenderer fyne.WidgetRenderer
	scroll       *container.Scroll
}

func (r *transparentScrollRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *transparentScrollRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *transparentScrollRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *transparentScrollRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var transparentObjects []fyne.CanvasObject
	
	for i, obj := range objects {
		fmt.Printf("DEBUG: Scroll object %d: %T\n", i, obj)
		if text, ok := obj.(*canvas.Text); ok {
			fmt.Printf("DEBUG: Scroll - Remplacement texte '%s' par transparent\n", text.Text)
			// Utiliser une couleur vraiment transparente
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			transparentObjects = append(transparentObjects, transparentText)
		} else if objType := fmt.Sprintf("%T", obj); objType == "*widget.entryContent" {
			fmt.Printf("DEBUG: Trouvé entryContent, création d'un wrapper transparent\n")
			// Créer un wrapper pour l'entryContent qui rend le texte transparent
			transparentEntry := &transparentEntryContentWrapper{content: obj}
			transparentObjects = append(transparentObjects, transparentEntry)
		} else {
			// Garder les autres objets (scrollbars, shadows, etc.)
			transparentObjects = append(transparentObjects, obj)
		}
	}
	
	return transparentObjects
}

func (r *transparentScrollRenderer) makeObjectsTransparent(obj fyne.CanvasObject, result *[]fyne.CanvasObject) {
	if text, ok := obj.(*canvas.Text); ok {
		fmt.Printf("DEBUG: Recursive - Remplacement texte '%s' par transparent\n", text.Text)
		transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
		transparentText.TextStyle = text.TextStyle
		transparentText.Alignment = text.Alignment
		transparentText.Move(text.Position())
		transparentText.Resize(text.Size())
		*result = append(*result, transparentText)
	} else if container, ok := obj.(*fyne.Container); ok {
		// Explorer les objets du container
		newContainer := container
		for _, childObj := range container.Objects {
			r.makeObjectsTransparent(childObj, result)
		}
		*result = append(*result, newContainer)
	} else {
		*result = append(*result, obj)
	}
}

func (r *transparentScrollRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// transparentEntryContentWrapper encapsule un entryContent et rend son texte transparent
type transparentEntryContentWrapper struct {
	content fyne.CanvasObject
}

func (t *transparentEntryContentWrapper) Position() fyne.Position {
	return t.content.Position()
}

func (t *transparentEntryContentWrapper) Size() fyne.Size {
	return t.content.Size()
}

func (t *transparentEntryContentWrapper) MinSize() fyne.Size {
	return t.content.MinSize()
}

func (t *transparentEntryContentWrapper) Move(pos fyne.Position) {
	t.content.Move(pos)
}

func (t *transparentEntryContentWrapper) Resize(size fyne.Size) {
	t.content.Resize(size)
}

func (t *transparentEntryContentWrapper) Hide() {
	t.content.Hide()
}

func (t *transparentEntryContentWrapper) Show() {
	t.content.Show()
}

func (t *transparentEntryContentWrapper) Visible() bool {
	return t.content.Visible()
}

func (t *transparentEntryContentWrapper) Refresh() {
	t.content.Refresh()
}

// Implémente fyne.Widget pour pouvoir créer un renderer
func (t *transparentEntryContentWrapper) CreateRenderer() fyne.WidgetRenderer {
	// Utiliser la réflexion pour obtenir le renderer de l'entryContent original
	if widget, ok := t.content.(fyne.Widget); ok {
		baseRenderer := widget.CreateRenderer()
		return &transparentEntryContentRenderer{
			baseRenderer: baseRenderer,
			content:      t.content,
		}
	}
	// Fallback si ce n'est pas un widget
	return &transparentEntryContentRenderer{
		baseRenderer: nil,
		content:      t.content,
	}
}

// transparentEntryContentRenderer rend l'entryContent avec texte transparent
type transparentEntryContentRenderer struct {
	baseRenderer fyne.WidgetRenderer
	content      fyne.CanvasObject
}

func (r *transparentEntryContentRenderer) Layout(size fyne.Size) {
	if r.baseRenderer != nil {
		r.baseRenderer.Layout(size)
	}
}

func (r *transparentEntryContentRenderer) MinSize() fyne.Size {
	if r.baseRenderer != nil {
		return r.baseRenderer.MinSize()
	}
	return r.content.MinSize()
}

func (r *transparentEntryContentRenderer) Refresh() {
	if r.baseRenderer != nil {
		r.baseRenderer.Refresh()
	}
}

func (r *transparentEntryContentRenderer) Objects() []fyne.CanvasObject {
	if r.baseRenderer == nil {
		return []fyne.CanvasObject{r.content}
	}
	
	objects := r.baseRenderer.Objects()
	var transparentObjects []fyne.CanvasObject
	
	for i, obj := range objects {
		fmt.Printf("DEBUG: EntryContent object %d: %T\n", i, obj)
		if text, ok := obj.(*canvas.Text); ok {
			fmt.Printf("DEBUG: EntryContent - Remplacement canvas.Text '%s' par transparent\n", text.Text)
			// Utiliser une couleur vraiment transparente
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			transparentObjects = append(transparentObjects, transparentText)
		} else if rect, ok := obj.(*canvas.Rectangle); ok {
			fmt.Printf("DEBUG: EntryContent - Remplacement Rectangle background par transparent\n")
			// Rendre le background transparent
			transparentRect := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
			transparentRect.Move(rect.Position())
			transparentRect.Resize(rect.Size())
			transparentObjects = append(transparentObjects, transparentRect)
		} else if richText, ok := obj.(*widget.RichText); ok {
			fmt.Printf("DEBUG: EntryContent - Trouvé RichText, création d'un wrapper transparent\n")
			// Créer un wrapper pour le RichText qui rend son contenu transparent
			transparentRichText := &transparentRichTextWrapper{RichText: richText}
			transparentObjects = append(transparentObjects, transparentRichText)
		} else {
			transparentObjects = append(transparentObjects, obj)
		}
	}
	
	return transparentObjects
}

func (r *transparentEntryContentRenderer) Destroy() {
	if r.baseRenderer != nil {
		r.baseRenderer.Destroy()
	}
}

// transparentRichTextWrapper encapsule un RichText et rend son texte transparent
type transparentRichTextWrapper struct {
	*widget.RichText
}

func (t *transparentRichTextWrapper) CreateRenderer() fyne.WidgetRenderer {
	baseRenderer := t.RichText.CreateRenderer()
	return &transparentRichTextRenderer{
		baseRenderer: baseRenderer,
		richText:     t.RichText,
	}
}

// transparentRichTextRenderer rend le RichText avec texte transparent
type transparentRichTextRenderer struct {
	baseRenderer fyne.WidgetRenderer
	richText     *widget.RichText
}

func (r *transparentRichTextRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *transparentRichTextRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *transparentRichTextRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *transparentRichTextRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var transparentObjects []fyne.CanvasObject
	
	for i, obj := range objects {
		fmt.Printf("DEBUG: RichText object %d: %T\n", i, obj)
		if text, ok := obj.(*canvas.Text); ok {
			fmt.Printf("DEBUG: RichText - Remplacement texte '%s' par transparent\n", text.Text)
			// Utiliser une couleur vraiment transparente
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			transparentObjects = append(transparentObjects, transparentText)
		} else {
			transparentObjects = append(transparentObjects, obj)
		}
	}
	
	return transparentObjects
}

func (r *transparentRichTextRenderer) Destroy() {
	r.baseRenderer.Destroy()
}


// hiddenTextEntry est un Entry dont le texte a la même couleur que le background
type hiddenTextEntry struct {
	*widget.Entry
}

func (h *hiddenTextEntry) CreateRenderer() fyne.WidgetRenderer {
	// Créer un thème temporaire qui rend le texte de la même couleur que le background
	hiddenTheme := &hiddenTextTheme{}
	
	// Appliquer temporairement ce thème
	originalTheme := fyne.CurrentApp().Settings().Theme()
	fyne.CurrentApp().Settings().SetTheme(hiddenTheme)
	
	// Créer le renderer
	renderer := h.Entry.CreateRenderer()
	
	// Restaurer le thème original dans une goroutine pour éviter les boucles
	go func() {
		fyne.CurrentApp().Settings().SetTheme(originalTheme)
	}()
	
	return renderer
}

// hiddenTextTheme rend le texte de la même couleur que le background
type hiddenTextTheme struct{}

func (h *hiddenTextTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		// Retourner la couleur de background pour "cacher" le texte
		return theme.DefaultTheme().Color(theme.ColorNameBackground, variant)
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (h *hiddenTextTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (h *hiddenTextTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (h *hiddenTextTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// transparentTextEntry encapsule un Entry et force la transparence du texte
type transparentTextEntry struct {
	widget.BaseWidget
	Entry *widget.Entry
}

func (t *transparentTextEntry) CreateRenderer() fyne.WidgetRenderer {
	renderer := t.Entry.CreateRenderer()
	return &transparentTextRenderer{
		baseRenderer: renderer,
		entry:        t.Entry,
	}
}

func (t *transparentTextEntry) Resize(size fyne.Size) {
	t.BaseWidget.Resize(size)
	if t.Entry != nil {
		t.Entry.Resize(size)
	}
}

func (t *transparentTextEntry) Move(pos fyne.Position) {
	t.BaseWidget.Move(pos)
	if t.Entry != nil {
		t.Entry.Move(pos)
	}
}

func (t *transparentTextEntry) MinSize() fyne.Size {
	if t.Entry != nil {
		return t.Entry.MinSize()
	}
	return fyne.NewSize(0, 0)
}

// transparentTextRenderer rend le texte transparent
type transparentTextRenderer struct {
	baseRenderer fyne.WidgetRenderer
	entry        *widget.Entry
}

func (r *transparentTextRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *transparentTextRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *transparentTextRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *transparentTextRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var result []fyne.CanvasObject
	
	for _, obj := range objects {
		if text, ok := obj.(*canvas.Text); ok {
			// Créer un texte transparent
			transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
			transparentText.TextStyle = text.TextStyle
			transparentText.Alignment = text.Alignment
			transparentText.Move(text.Position())
			transparentText.Resize(text.Size())
			result = append(result, transparentText)
			fmt.Printf("DEBUG: TransparentText - Rendu texte transparent: '%s'\n", text.Text)
		} else {
			result = append(result, obj)
		}
	}
	
	return result
}

func (r *transparentTextRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// forceTransparentEntry est un Entry qui force la transparence du texte
type forceTransparentEntry struct {
	widget.Entry
}

func (f *forceTransparentEntry) CreateRenderer() fyne.WidgetRenderer {
	renderer := f.Entry.CreateRenderer()
	return &forceTransparentRenderer{
		baseRenderer: renderer,
		entry:        &f.Entry,
	}
}

// forceTransparentRenderer force tous les textes à être transparents
type forceTransparentRenderer struct {
	baseRenderer fyne.WidgetRenderer
	entry        *widget.Entry
}

func (r *forceTransparentRenderer) Layout(size fyne.Size) {
	r.baseRenderer.Layout(size)
}

func (r *forceTransparentRenderer) MinSize() fyne.Size {
	return r.baseRenderer.MinSize()
}

func (r *forceTransparentRenderer) Refresh() {
	r.baseRenderer.Refresh()
}

func (r *forceTransparentRenderer) Objects() []fyne.CanvasObject {
	objects := r.baseRenderer.Objects()
	var result []fyne.CanvasObject
	
	for _, obj := range objects {
		r.makeTransparent(obj, &result)
	}
	
	return result
}

func (r *forceTransparentRenderer) makeTransparent(obj fyne.CanvasObject, result *[]fyne.CanvasObject) {
	if text, ok := obj.(*canvas.Text); ok {
		// Forcer la couleur à être vraiment transparente
		transparentText := canvas.NewText(text.Text, color.RGBA{0, 0, 0, 0})
		transparentText.TextStyle = text.TextStyle
		transparentText.Alignment = text.Alignment
		transparentText.Move(text.Position())
		transparentText.Resize(text.Size())
		*result = append(*result, transparentText)
		fmt.Printf("DEBUG: ForceTransparent - Texte transparent: '%s'\n", text.Text)
	} else if rect, ok := obj.(*canvas.Rectangle); ok {
		// Garder les rectangles transparents
		transparentRect := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
		transparentRect.Move(rect.Position())
		transparentRect.Resize(rect.Size())
		*result = append(*result, transparentRect)
	} else if container, ok := obj.(*fyne.Container); ok {
		// Traiter récursivement les containers
		newContainer := container
		newContainer.Objects = nil
		for _, childObj := range container.Objects {
			r.makeTransparent(childObj, &newContainer.Objects)
		}
		*result = append(*result, newContainer)
	} else {
		*result = append(*result, obj)
	}
}

func (r *forceTransparentRenderer) Destroy() {
	r.baseRenderer.Destroy()
}

// transparentTheme est un thème qui rend le texte transparent
type transparentTheme struct {
	fyne.Theme
}

func (t *transparentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Rendre le texte de la même couleur que le background pour le "cacher"
	if name == theme.ColorNameForeground {
		return theme.DefaultTheme().Color(theme.ColorNameBackground, variant)
	}
	
	// Pour tous les autres éléments, utiliser le thème par défaut
	return theme.DefaultTheme().Color(name, variant)
}

func (t *transparentTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *transparentTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *transparentTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}