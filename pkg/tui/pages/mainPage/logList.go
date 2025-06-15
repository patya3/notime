package mainPage

import (
	"log"

	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/logModal"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/patya3/notime/pkg/utils"

	"github.com/gdamore/tcell/v2"
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
				logModal.InitFormElements(app, pagePrimitive, logType, &issueLogs[i].ID)
				app.SetFocus(logModal.LogModalForm)
			} else if logType == "QUICK_LOG" {
				// logModal.SetLogModalTextForQuickLog(quickLogs[i].ID)
				logModal.InitFormElements(app, pagePrimitive, logType, &quickLogs[i].ID)
				app.SetFocus(logModal.LogModalForm)
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			// NOTE: copy not working correctly at the moment:
			// dont take notes after copy and display comment on running time
			case 'c':
				currentIssueId := issues[IssueList.GetCurrentItem()].ID
				currentLogId := issueLogs[list.GetCurrentItem()].ID

				copiedTimeLog, err := constants.LogRepo.CopyTimerByLogAndIssueId(currentLogId, currentIssueId)
				if err != nil {
					log.Fatal(err)
				}
				LogList.InsertItem(0, copiedTimeLog.Title(), copiedTimeLog.Comment, 0, nil)
				break
			case 'D':
				// DELETE
				helpers.SetConfirmationModal("Are you sure you want to delete this log?", app, pagePrimitive, LogList, func() {
					currentLogId := issueLogs[list.GetCurrentItem()].ID
					err := constants.LogRepo.DeleteLogByID(currentLogId)
					if err != nil {
						log.Fatal(err)
					}
					LogList.RemoveItem(list.GetCurrentItem())
					pagePrimitive.HidePage("ConfirmationModal")
				})
				pagePrimitive.ShowPage("ConfirmationModal")
				break
			case 'L':
				// LOG to jira
				helpers.SetConfirmationModal("Are you sure you want to LOG this entry?", app, pagePrimitive, LogList, func() {
					// timelog timelog.ExtendedLog
					currentLogId := issueLogs[list.GetCurrentItem()].ID
					extendedTimelog, err := constants.LogRepo.GetLogByID(currentLogId)

					if err != nil {
						notification.SetNotification(err.Error())
						pagePrimitive.HidePage("ConfirmationModal")
						pagePrimitive.ShowPage("Notification")
						return
					}

					_, err = utils.CreateJiraWorkLog(extendedTimelog)
					if err != nil {
						notification.SetNotification(err.Error())
						pagePrimitive.HidePage("ConfirmationModal")
						pagePrimitive.ShowPage("Notification")
						return
					}
					LogList.RemoveItem(list.GetCurrentItem())
					pagePrimitive.HidePage("ConfirmationModal")
				})
				pagePrimitive.ShowPage("ConfirmationModal")
				break
			case 'A': // TODO:
				// assign to quicklog list to issue
				if logType == "ISSUE_LOG" {
					// display a modal where you can select from already synced issues
				}
				break

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
