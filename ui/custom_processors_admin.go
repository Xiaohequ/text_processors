package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// CreateManageCustomProcessorsDialog affiche une fenêtre de gestion des processeurs personnalisés
func CreateManageCustomProcessorsDialog(parent fyne.Window) {
	// Widgets principaux
	list := widget.NewList(
		func() int { return len(GlobalCustomProcessorManager.GetProcessors()) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			procs := GlobalCustomProcessorManager.GetProcessors()
			if i >= 0 && i < len(procs) {
				o.(*widget.Label).SetText(procs[i].Name)
			}
		},
	)

	var selectedIndex int = -1
	list.OnSelected = func(id widget.ListItemID) { selectedIndex = int(id) }

	// Boutons d'action
	renameBtn := widget.NewButton("Renommer", func() {
		if selectedIndex < 0 {
			return
		}
		procs := GlobalCustomProcessorManager.GetProcessors()
		if selectedIndex >= len(procs) {
			return
		}
		current := procs[selectedIndex]
		nameEntry := widget.NewEntry()
		nameEntry.SetText(current.Name)
		content := container.NewVBox(
			widget.NewLabel("Nouveau nom:"),
			nameEntry,
		)
		d := dialog.NewCustomConfirm("Renommer", "OK", "Annuler", content, func(ok bool) {
			if !ok {
				return
			}
			newName := nameEntry.Text
			if newName == "" {
				dialog.ShowError(fmt.Errorf("nom vide"), parent)
				return
			}
			// Mise à jour en mémoire
			procs := GlobalCustomProcessorManager.GetProcessors()
			procs[selectedIndex].Name = newName
			// Appliquer la modification et sauvegarder
			// Remplacer la liste interne par la version modifiée
			GlobalCustomProcessorManager.definitions = procs
			if err := GlobalCustomProcessorManager.SaveAll(); err != nil {
				dialog.ShowError(err, parent)
			}
			if GlobalCustomProcessorManager.onUpdate != nil {
				GlobalCustomProcessorManager.onUpdate()
			}
			list.Refresh()
		}, parent)
		d.Show()
	})

	deleteBtn := widget.NewButton("Supprimer", func() {
		if selectedIndex < 0 {
			return
		}
		confirm := dialog.NewConfirm("Supprimer", "Supprimer le processeur sélectionné ?", func(ok bool) {
			if !ok {
				return
			}
			GlobalCustomProcessorManager.RemoveProcessor(selectedIndex)
			selectedIndex = -1
			list.Refresh()
		}, parent)
		confirm.Show()
	})

	// Bouton éditer (nom + script)
	editBtn := widget.NewButton("Éditer", func() {
		if selectedIndex < 0 {
			return
		}
		procs := GlobalCustomProcessorManager.GetProcessors()
		if selectedIndex >= len(procs) {
			return
		}
		current := procs[selectedIndex]
		nameEntry := widget.NewEntry()
		nameEntry.SetText(current.Name)
		scriptEntry := widget.NewMultiLineEntry()
		scriptEntry.SetPlaceHolder("Script JavaScript (fonction process(input) { return input; })")
		scriptEntry.SetText(current.Script)
		// UI du dialogue
		content := container.NewVBox(
			widget.NewLabel("Nom:"),
			nameEntry,
			widget.NewLabel("Script JavaScript:"),
			container.NewVScroll(scriptEntry),
		)
		d := dialog.NewCustomConfirm("Éditer le processeur", "Enregistrer", "Annuler", content, func(ok bool) {
			if !ok {
				return
			}
			newName := nameEntry.Text
			newScript := scriptEntry.Text
			if newName == "" {
				dialog.ShowError(fmt.Errorf("le nom est requis"), parent)
				return
			}
			if newScript == "" {
				dialog.ShowError(fmt.Errorf("le script est requis"), parent)
				return
			}
			// Appliquer
			procs := GlobalCustomProcessorManager.GetProcessors()
			procs[selectedIndex].Name = newName
			procs[selectedIndex].Script = newScript
			GlobalCustomProcessorManager.definitions = procs
			if err := GlobalCustomProcessorManager.SaveAll(); err != nil {
				dialog.ShowError(err, parent)
			}
			if GlobalCustomProcessorManager.onUpdate != nil {
				GlobalCustomProcessorManager.onUpdate()
			}
			list.Refresh()
		}, parent)
		d.Resize(fyne.NewSize(600, 500))
		d.Show()
	})

	// Layout
	actions := container.NewHBox(renameBtn, editBtn, deleteBtn)
	content := container.NewBorder(nil, actions, nil, nil, list)

	dlg := dialog.NewCustom("Gérer les processeurs personnalisés", "Fermer", content, parent)
	dlg.Resize(fyne.NewSize(500, 400))
	dlg.Show()
}
