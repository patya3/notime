package mainPage

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/helpers"
	"github.com/rivo/tview"
)

var issueList = tview.NewList()
var issues = make([]issue.Issue, 0)
var selectedIssueId uint

func InitIssuesList(app *tview.Application, pagePrimitive *tview.Pages) *tview.List {

	issueList.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitle("Issues").
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			issueList.SetBorderColor(tcell.ColorRed)
		}).
		SetBlurFunc(func() {
			issueList.SetBorderColor(tcell.ColorDefault)
		})

	issueList.
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorLightPink.TrueColor()).
		SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			if selectedIssueId != issues[i].ID {
				initLogListElements(issues[i].ID)
				selectedIssueId = issues[i].ID
			}
			app.SetFocus(logsList)
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'c':
				pagePrimitive.ShowPage("AddIssue")
			}
			return helpers.RedifineUpAndDown(event)
		})

	initIssueListElements()

	return issueList
}

func initIssueListElements() {
	var err error
	issues, err = constants.IssueRepo.GetAllIssues()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(issues))
	for _, issue := range issues {
		issueList.AddItem(issue.FilterValue(), "", 0, nil)
	}
}
