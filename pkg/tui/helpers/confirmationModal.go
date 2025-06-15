package helpers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ConfirmationModal = tview.NewModal()

func SetConfirmationModal(text string, app *tview.Application, pagePrimitive *tview.Pages, elementToFocus tview.Primitive, callback func()) {
	ConfirmationModal.ClearButtons()

	ConfirmationModal.
		SetText(text).
		AddButtons([]string{"Yes", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				pagePrimitive.HidePage("ConfirmationModal")
			} else if buttonLabel == "Yes" {
				callback()
			}
			app.SetFocus(elementToFocus)
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
				pagePrimitive.HidePage("ConfirmationModal")
				app.SetFocus(elementToFocus)
			}
			return event
		})

}
