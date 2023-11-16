package mainPage

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/logModal"
	"github.com/patya3/notime/pkg/tui/pages/notification"

	// "github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var issueLogs = make([]timelog.Log, 0)
var quickLogs = make([]timelog.Log, 0)

// @param logType = "ISSUE_LOG" | "QUICK_LOG"
func InitLogList(list *tview.List, logType string, app *tview.Application, pagePrimitive *tview.Pages) {
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
				// logModal.SetLogModalTextIssueLog(issueLogs[i].ID)
				logModal.InitFormElements(pagePrimitive, &issueLogs[i].ID)
				app.SetFocus(logModal.LogModalForm)
			} else if logType == "QUICK_LOG" {
				logModal.SetLogModalTextForQuickLog(quickLogs[i].ID)
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			// NOTE: not working correctly at the moment:
			// dont take notes after copy and display comment on running time
			case 'c':
				currentIssueId := issues[IssueList.GetCurrentItem()].ID
				currentLogId := issueLogs[list.GetCurrentItem()].ID

				copiedTimeLog, err := constants.LogRepo.CopyTimerByLogAndIssueId(currentLogId, currentIssueId)
				if err != nil {
					log.Fatal(err)
					// notification.SetNotification("Something went wrong cannot copy Log.")
					// pagePrimitive.ShowPage("Notification")
				}
				LogList.InsertItem(0, copiedTimeLog.Title(), copiedTimeLog.Comment, 0, nil)
				break
			case 'u':
				// NOTE: UPDATE
				currentLogId := issueLogs[list.GetCurrentItem()].ID
				log, err := constants.LogRepo.GetLogByID(currentLogId)
				if err != nil {
					pagePrimitive.ShowPage("Notification")
					notification.SetNotification("Log not found")
				}
				fmt.Println(log)

			}
			return helpers.RedifineUpAndDown(event)
		})
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
