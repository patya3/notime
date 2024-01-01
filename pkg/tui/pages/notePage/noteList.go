package notePage

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/note"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/noteModal"
	"github.com/patya3/notime/pkg/utils"
	"github.com/rivo/tview"
)

var notes = make([]note.Note, 0)
var vimNotes = make([]utils.Hit, 0)

// @param type = "Notes" | "VimNotes"
func InitNotesList(app *tview.Application, pagePrimitive *tview.Pages, list *tview.List, listType string) {

	list.
		SetBorder(true).
		SetTitle(listType).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			list.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			list.SetBorderColor(tcell.ColorDefault)
		})

	list.
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorLightPink.TrueColor()).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			pagePrimitive.ShowPage("Note")
			if listType == "Notes" {
				noteModal.SetNoteModalText(notes[i].ID)
			} else {
				// log.Println(vimNotes, i)
				noteModal.SetNoteModalTextForVimNote(vimNotes[i])
			}
		}).
		SetInputCapture(helpers.RedifineUpAndDown)

	if listType == "Notes" {
		InitNoteListElements()
	} else {
		InitVimNotesListElements()
	}
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

func InitVimNotesListElements() {
	VimNoteList.Clear()
	var err error
	vimNotes, err = utils.ParseVimNotes()

	if err != nil {
		log.Fatal(err)
	}

	for _, vimNote := range vimNotes {
		VimNoteList.AddItem(vimNote.Note, "", 0, nil)
	}
}
