package mainPage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var issues = make([]issue.Issue, 0)
var SelectedIssueId sql.NullInt32

func InitIssueList(app *tview.Application, pagePrimitive *tview.Pages) {
	IssueList.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitle("Issues").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			IssueList.SetBorderColor(tcell.ColorRed)
			for i, issue := range issues {
				if len(issue.Logs) > 0 && !issue.Logs[0].StoppedAt.Valid {
					IssueList.SetCurrentItem(i)
					break
				}
			}
		}).
		SetBlurFunc(func() {
			IssueList.SetBorderColor(tcell.ColorDefault)
		})

	IssueList.
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorDarkSlateGray).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			app.SetFocus(LogList)
		}).
		SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
			if uint(SelectedIssueId.Int32) != issues[i].ID {
				LogList.Clear()
				InitLogListElements(issues[i].ID)
				SelectedIssueId.Scan(issues[i].ID)
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			switch event.Rune() {
			case 'c':
				pagePrimitive.ShowPage("AddIssue")
			case 'a':
				if len(issues) == 0 {
					notification.SetNotification("Create an issue first with the letter 'c'.")
					pagePrimitive.ShowPage("Notification")
					break
				}
				hasRunningLog, err := constants.LogRepo.HasRunningLog(SelectedIssueId)
				if err != nil {
					log.Fatal(err)
				}
				if hasRunningLog {
					notification.SetNotification("This Issue already has a running log.")
					pagePrimitive.ShowPage("Notification")
				} else {
					timelog, err := constants.LogRepo.StartTimer(SelectedIssueId)
					if err != nil {
						log.Fatal(err)
					}
					LogList.InsertItem(0, timelog.Title(), timelog.Comment, 0, nil)
				}
			case 's':
				if len(issues) == 0 {
					break
				}
				currentIssueId := issues[IssueList.GetCurrentItem()].ID
				hasRunningLog, err := constants.LogRepo.HasRunningLog(sql.NullInt32{Int32: int32(currentIssueId), Valid: true})
				if err != nil {
					log.Fatal(err)
				}
				if hasRunningLog {
					fmt.Println("hello", currentIssueId)
					pagePrimitive.ShowPage("AddComment")
				} else {
					notification.SetNotification("This Issue hasn't got a running log.")
					pagePrimitive.ShowPage("Notification")
				}
			}
			return helpers.RedifineUpAndDown(event)
		})

	InitIssueListElements()
}

func InitIssueListElements() {
	var err error
	IssueList.Clear()
	issues, err = constants.IssueRepo.GetAllIssues()
	if err != nil {
		log.Fatal(err)
	}
	for _, issue := range issues {
		IssueList.AddItem("[lightgreen]("+issue.IssueKey+") [pink]"+issue.IssueTitle, "", 0, nil)
	}
}
