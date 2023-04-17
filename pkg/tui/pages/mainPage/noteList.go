package mainPage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/rivo/tview"
)

func InitNotesList(app *tview.Application) {

	NoteList.
		SetBorder(true).
		SetTitle("Notes").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			NoteList.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			NoteList.SetBorderColor(tcell.ColorDefault)
		})

	NoteList.
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorLightPink.TrueColor()).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			// app.SetFocus(notesList)
		}).
		SetInputCapture(helpers.RedifineUpAndDown)

	NoteList.AddItem("Szia", "minden rendben", 0, nil)
	NoteList.AddItem("Szia", "minden rendben", 0, nil)
}
