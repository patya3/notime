package mainPage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/rivo/tview"
)

var notesList = tview.NewList()

func InitNotesList(app *tview.Application) *tview.List {

	notesList.
		SetBorder(true).
		SetTitle("Notes").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			notesList.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			notesList.SetBorderColor(tcell.ColorDefault)
		})

	notesList.
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorLightPink.TrueColor()).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			// app.SetFocus(notesList)
		}).
		SetInputCapture(helpers.RedifineUpAndDown)

	notesList.AddItem("Szia", "minden rendben", 0, nil)
	notesList.AddItem("Szia", "minden rendben", 0, nil)

	return notesList
}
