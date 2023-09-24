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

var issueLogs = make([]timelog.Log, 0)
var quickLogs = make([]timelog.Log, 0)

// @param logType = "ISSUE_LOG" | "QUICK_LOG"
func InitLogList(list *tview.List, logType string, pagePrimitive *tview.Pages) {
	title := "Logs"
	if logType == "QUICK_LOG" {
		title = "Quick Logs"
	}

	list.Box.
		SetBorder(true).
		SetTitle(title).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			list.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			list.SetBorderColor(tcell.ColorDefault)
		})

	list.
		SetHighlightFullLine(true).
		SetSecondaryTextColor(tcell.ColorLightGreen).
		SetSelectedBackgroundColor(tcell.ColorDarkSlateGray).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			pagePrimitive.ShowPage("Log")
			if logType == "ISSUE_LOG" {
				logModal.SetLogModalText(issueLogs[i].ID)
			} else if logType == "QUICK_LOG" {
				logModal.SetLogModalText(quickLogs[i].ID)
			}
		}).
		SetInputCapture(helpers.RedifineUpAndDown)
}

func InitLogListElements(issueID uint) {
	var err error
	LogList.Clear()
	issueLogs, err = constants.LogRepo.GetAllLogsForAnIssue(issueID)
	if err != nil {
		log.Fatal(err)
	}
	for _, log := range issueLogs {
		LogList.AddItem(log.Title(), log.Comment, 0, nil)
	}
}

func InitQuickLogListElements() {
	var err error
	QuickLogList.Clear()
	quickLogs, err = constants.LogRepo.GetAllQuickLogs()
	if err != nil {
		log.Fatal(err)
	}
	for _, log := range quickLogs {
		QuickLogList.AddItem(log.Title(), log.Comment, 0, nil)
	}
}
