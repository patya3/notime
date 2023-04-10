package mainPage

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var IssueList = tview.NewList()
var issues = make([]issue.Issue, 0)
var SelectedIssueId uint

func InitIssuesList(app *tview.Application, pagePrimitive *tview.Pages) *tview.List {

	IssueList.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitle("Issues").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			IssueList.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			IssueList.SetBorderColor(tcell.ColorDefault)
		})

	IssueList.
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorDarkSlateGray).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			app.SetFocus(logsList)
		}).
		SetChangedFunc(func(i int, mainText, secondaryText string, shortcut rune) {
			if SelectedIssueId != issues[i].ID {
				logsList.Clear()
				InitLogListElements(issues[i].ID)
				SelectedIssueId = issues[i].ID
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			switch event.Rune() {
			case 'c':
				pagePrimitive.ShowPage("AddIssue")
			case 'a':
				if len(issues) == 0 {
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
					logsList.AddItem(timelog.Title(true), timelog.Description(), 0, nil)
				}
			case 's':
				if len(issues) == 0 {
					break
				}
				currentIssueId := issues[IssueList.GetCurrentItem()].ID
				hasRunningLog, err := constants.LogRepo.HasRunningLog(currentIssueId)
				if err != nil {
					log.Fatal(err)
				}
				if !hasRunningLog {
					notification.SetNotification("This Issue hasn't got a running log.")
					pagePrimitive.ShowPage("Notification")
				} else {
					pagePrimitive.ShowPage("AddComment")
				}
			}
			return helpers.RedifineUpAndDown(event)
		})

	InitIssueListElements()

	return IssueList
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
