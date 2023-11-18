package notePage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var NotePageContainer = tview.NewFlex()
var NoteList = tview.NewList()
var VimNoteList = tview.NewList()

func InitNotePage(app *tview.Application, pagePrimitive *tview.Pages) *tview.Flex {
	NotePageContainer.SetBackgroundColor(tcell.ColorDefault)

	InitNotesList(app, pagePrimitive, NoteList, "Notes")
	InitNotesList(app, pagePrimitive, VimNoteList, "VimNotes")
	// InitVimNoteListElements()

	NotePageContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '1':
			app.SetFocus(NoteList)
			break
		case '2':
			app.SetFocus(VimNoteList)
			break
		case 'q':
			pagePrimitive.HidePage("NotePage")
			break
		}

		return event
	})

	NotePageContainer.AddItem(
		tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NoteList, 0, 5, true).
			AddItem(VimNoteList, 0, 5, false),
		0, 1, true)

	return NotePageContainer
}

// case '3':
// 	app.SetFocus(NoteList)
// 	break
