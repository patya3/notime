package addNoteModal

import (
	"github.com/gdamore/tcell/v2"
	noteModel "github.com/patya3/notime/pkg/models/note"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/mainPage"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var AddNoteModal = tview.NewForm()

func InitAddNoteModal(app *tview.Application, pagePrimitive *tview.Pages) tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	AddNoteModal.Box.
		SetTitle("Add new note").
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
				pagePrimitive.HidePage("AddNote")
			}
			return event
		})

	initFormElements(AddNoteModal, app, pagePrimitive)

	return modal(AddNoteModal, 50, 30)
}

func initFormElements(noteForm *tview.Form, app *tview.Application, pagePrimitive *tview.Pages) {
	noteForm.Clear(true)

	// TODO: meg nincs meg ez a szar
	note := noteModel.Note{}

	AddNoteModal.
		AddInputField("Title", "", 0, nil, func(text string) {
			note.Title = text
		}).
		// AddTextArea("Title", "", 0, 2, 0, func(text string) {
		// 	issue.IssueTitle = text
		// }).
		AddTextArea("Description", "", 0, 0, 0, func(text string) {
			note.Description = text
		}).
		AddButton("Save", func() {
			_, err := constants.NoteRepo.CreateNote(note)
			if err != nil {
				notification.SetNotification(err.Error())
				pagePrimitive.ShowPage("Notification")
			}
			mainPage.InitNoteListElements()
			pagePrimitive.HidePage("AddNote")
		}).
		AddButton("Quit", func() {
			pagePrimitive.HidePage("AddNote")
		}).
		AddButton("Reset", func() {
			initFormElements(noteForm, app, pagePrimitive)
			app.SetFocus(noteForm)
		})
}
