package helpers

import (
	// "github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// var FocusedList *tview.List

func SetFocusedList(app *tview.Application, newFocusedList *tview.List, oldFocusedList *tview.List) {
	if oldFocusedList != nil {
		// oldFocusedList.SetBorderColor(tcell.ColorDefault)
	}
	if newFocusedList != nil {
		// newFocusedList.SetBorderColor(tcell.ColorRed)
	}
}

func RedifineUpAndDown(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, rune(tcell.KeyDown), tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone)
	}
	return event
}

// func FilterList(list *tview.List) {
//     list.FindItems()
// }
