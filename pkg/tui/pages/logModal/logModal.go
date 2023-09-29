package logModal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	c "github.com/patya3/notime/pkg/colors"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/rivo/tview"
)

var logModal = tview.NewTextView()

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitLogModal(pagePrimitive *tview.Pages) *tview.Grid {
	logModal.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' {
				pagePrimitive.HidePage("Log")
				logModal.SetText("")
				return nil
			}
			return event
		})

	logModal.
		SetTextColor(tcell.ColorIndianRed).
		SetDynamicColors(true).
		SetBackgroundColor(tcell.ColorDefault)

	return modal(logModal, 40, 20)
}

func SetLogModalTextIssueLog(logID uint) {
	extendedTimelog, err := constants.LogRepo.GetLogByID(logID)
	if err != nil {
		log.Fatal(err)
	}
	logModal.SetText(fmt.Sprintf("Issue: %s%s\n\n%s\n[lightgreen]%s\n\n%s%s", extendedTimelog.IssueKey, extendedTimelog.IssueTitle, extendedTimelog.Title(), extendedTimelog.Description(), c.Colors["white"], extendedTimelog.Comment))
}

func SetLogModalTextForQuickLog(quickLogID uint) {
	quickLog, err := constants.LogRepo.GetQuickLogByID(quickLogID)
	if err != nil {
		log.Fatal(err)
	}
	logModal.SetText(fmt.Sprintf("%s\n[lightgreen]%s\n\n%s%s", quickLog.Title(), quickLog.Description(), c.Colors["white"], quickLog.Comment))
}
