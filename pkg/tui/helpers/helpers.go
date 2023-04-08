package helpers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// var FocusedList *tview.List

func SetFocusedList(app *tview.Application, newFocusedList *tview.List, oldFocusedList *tview.List) {
	if oldFocusedList != nil {
		oldFocusedList.SetBorderColor(tcell.ColorDefault)
	}
	if newFocusedList != nil {
		newFocusedList.SetBorderColor(tcell.ColorRed)
		app.SetFocus(newFocusedList)
	}
}
