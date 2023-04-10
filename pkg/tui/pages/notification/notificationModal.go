package notification

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Notification = tview.NewTextView()

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 1).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitNotification(pagePrimitive *tview.Pages) *tview.Grid {
	Notification.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			pagePrimitive.HidePage("Notification")
			return event
		})

	Notification.
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorIndianRed).
		SetBackgroundColor(tcell.ColorDefault)

	return modal(Notification, 40, 4)
}

func SetNotification(text string) {
	Notification.
		SetText(text)
}
