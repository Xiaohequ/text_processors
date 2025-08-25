package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TogglableCodeEditor est un éditeur avec bouton pour basculer entre édition et highlight
type TogglableCodeEditor struct {
	widget.BaseWidget
	editor     *RichCodeEditor
	toggleBtn  *widget.Button
	container  *fyne.Container
	isEditing  bool
}

// NewTogglableCodeEditor crée un nouvel éditeur avec bouton de bascule
func NewTogglableCodeEditor() *TogglableCodeEditor {
	t := &TogglableCodeEditor{
		isEditing: true, // Commencer en mode édition
	}
	
	// Créer l'éditeur RichText
	t.editor = NewRichCodeEditor()
	t.editor.EnableEditing() // Commencer en mode édition
	
	// Créer le bouton de bascule
	t.toggleBtn = widget.NewButton("🎨 Voir Couleurs", func() {
		t.toggle()
	})
	
	// Layout avec bouton en haut
	t.container = container.NewBorder(
		t.toggleBtn, // top
		nil,         // bottom
		nil,         // left
		nil,         // right
		t.editor,    // center
	)
	
	t.ExtendBaseWidget(t)
	return t
}

// toggle bascule entre mode édition et mode highlight
func (t *TogglableCodeEditor) toggle() {
	if t.isEditing {
		// Passer en mode highlight
		t.editor.EnableHighlighting()
		t.toggleBtn.SetText("✏️ Éditer")
		t.isEditing = false
	} else {
		// Passer en mode édition
		t.editor.EnableEditing()
		t.toggleBtn.SetText("🎨 Voir Couleurs")
		t.isEditing = true
	}
}

// SetText définit le texte
func (t *TogglableCodeEditor) SetText(text string) {
	t.editor.SetText(text)
}

// GetText retourne le texte actuel
func (t *TogglableCodeEditor) GetText() string {
	return t.editor.GetText()
}

// SetPlaceHolder définit le placeholder
func (t *TogglableCodeEditor) SetPlaceHolder(text string) {
	t.editor.SetPlaceHolder(text)
}

// OnChanged définit le callback
func (t *TogglableCodeEditor) OnChanged(callback func(string)) {
	t.editor.OnChanged(callback)
}

// CreateRenderer crée le renderer
func (t *TogglableCodeEditor) CreateRenderer() fyne.WidgetRenderer {
	return &togglableEditorRenderer{
		editor:    t,
		container: t.container,
		objects:   []fyne.CanvasObject{t.container},
	}
}

// togglableEditorRenderer est le renderer
type togglableEditorRenderer struct {
	editor    *TogglableCodeEditor
	container *fyne.Container
	objects   []fyne.CanvasObject
}

func (r *togglableEditorRenderer) Layout(size fyne.Size) {
	r.container.Resize(size)
}

func (r *togglableEditorRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *togglableEditorRenderer) Refresh() {
	r.container.Refresh()
}

func (r *togglableEditorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *togglableEditorRenderer) Destroy() {}

// Resize redimensionne l'éditeur
func (t *TogglableCodeEditor) Resize(size fyne.Size) {
	t.BaseWidget.Resize(size)
	t.container.Resize(size)
}