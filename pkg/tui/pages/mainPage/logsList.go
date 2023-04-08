package mainPage

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/rivo/tview"
)

var logsList = tview.NewList()

func InitLogsList() *tview.List {
	logsList.
		SetBorder(true).
		SetTitle("Logs").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			logsList.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			logsList.SetBorderColor(tcell.ColorDefault)
		})

	logsList.
		SetSelectedBackgroundColor(tcell.ColorLightPink.TrueColor()).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'j':
				return tcell.NewEventKey(tcell.KeyDown, rune(tcell.KeyDown), tcell.ModNone)
			case 'k':
				return tcell.NewEventKey(tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone)
			}
			return event
		})

	return logsList
}

func initLogListElements(issueID uint) {
	logs, err := constants.LogRepo.GetAllLogsForAnIssue(issueID)
	if err != nil {
		log.Fatal(err)
	}
	for _, log := range logs {
		logsList.AddItem(log.Title(), log.Description(), 0, nil)
	}
}
