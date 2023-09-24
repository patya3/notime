package mainPage

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/note"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/rivo/tview"
)

var notes = make([]note.Note, 0)

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

	InitNoteListElements()
}

func InitNoteListElements() {
	var err error
	NoteList.Clear()
	notes, err = constants.NoteRepo.GetAllNotes()
	if err != nil {
		log.Fatal(err)
	}
	for _, note := range notes {
		NoteList.AddItem(note.Title, "", 0, nil)
	}
}
