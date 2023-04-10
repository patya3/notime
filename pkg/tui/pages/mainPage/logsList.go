package mainPage

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/logModal"
	"github.com/rivo/tview"
)

var logsList = tview.NewList()
var logs = make([]timelog.Log, 0)

func InitLogsList(pagePrimitive *tview.Pages) *tview.List {
	logsList.Box.
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
		SetHighlightFullLine(true).
		SetSecondaryTextColor(tcell.ColorLightGreen).
		SetSelectedBackgroundColor(tcell.ColorDarkSlateGray).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			pagePrimitive.ShowPage("Log")
			logModal.SetLogModalText(logs[i].ID)
		}).
		SetInputCapture(helpers.RedifineUpAndDown)

	return logsList
}

func InitLogListElements(issueID uint) {
	var err error
	logsList.Clear()
	logs, err = constants.LogRepo.GetAllLogsForAnIssue(issueID)
	if err != nil {
		log.Fatal(err)
	}
	for _, log := range logs {
		logsList.AddItem(log.Title(true), log.Comment, 0, nil)
	}
}
