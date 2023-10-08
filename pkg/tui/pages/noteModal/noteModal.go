package noteModal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/rivo/tview"
)

var noteModal = tview.NewTextView()
var currentNoteId uint

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitNoteModal(app *tview.Application, pagePrimitive *tview.Pages, noteList *tview.List) *tview.Grid {
	noteModal.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' {
				pagePrimitive.HidePage("Note")
				noteModal.SetText("")
				app.SetFocus(noteList)
				return nil
			} else if event.Rune() == 'n' {
				setNextOrPreviousNoteModalText(true)
				return nil
			} else if event.Rune() == 'N' {
				setNextOrPreviousNoteModalText(false)
				return nil
			}
			return event
		})

	noteModal.
		SetTextColor(tcell.ColorIndianRed).
		SetDynamicColors(true).
		SetBackgroundColor(tcell.ColorDefault)

	return modal(noteModal, 80, 30)

}

func SetNoteModalText(noteID uint) {
	currentNoteId = noteID
	note, err := constants.NoteRepo.GetNoteByID(noteID)
	if err != nil {
		log.Fatal(err)
	}
	noteModal.SetText(fmt.Sprintf("%s\n\n[lightgreen]%s", note.Title, note.Description))
}

func setNextOrPreviousNoteModalText(isNext bool) {

	// NOTE: not the most performant (+ database query)
	noteIds, err := constants.NoteRepo.GetAllNoteIds()
	if err != nil {
		log.Fatal(err)
	}

	var currentNoteIdIndex int
	for index, noteId := range noteIds {
		if noteId == int(currentNoteId) {
			currentNoteIdIndex = index
			break
		}
	}

	var nextIndex int
	if isNext {
		nextIndex = currentNoteIdIndex + 1
	} else {
		nextIndex = currentNoteIdIndex - 1
	}
	if len(noteIds) > nextIndex && nextIndex >= 0 {
		nextNoteId := noteIds[nextIndex]
		SetNoteModalText(uint(nextNoteId))
		return
	}

	return
}
