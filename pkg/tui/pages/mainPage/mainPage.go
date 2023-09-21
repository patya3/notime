package mainPage

import (
	"database/sql"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var MainPageContainer = tview.NewFlex()
var IssueList = tview.NewList()
var LogList = tview.NewList()
var NoteList = tview.NewList()
var QuickLogList = tview.NewList()

func InitMainPage(app *tview.Application, pagePrimitive *tview.Pages) *tview.Flex {
	MainPageContainer.SetBackgroundColor(tcell.ColorDefault)

	InitIssueList(app, pagePrimitive)
	InitNotesList(app)
	InitLogList(LogList, "ISSUE_LOG", pagePrimitive)
	InitLogList(QuickLogList, "QUICK_LOG", pagePrimitive)
	InitQuickLogListElements()

	MainPageContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '1':
			app.SetFocus(IssueList)
			break
		case '2':
			app.SetFocus(LogList)
			break
		case '3':
			app.SetFocus(NoteList)
			break
		case '4':
			app.SetFocus(QuickLogList)
			break
		case 'A':
			timelog, err := constants.LogRepo.StartTimer(sql.NullInt32{})
			if err != nil {
				log.Fatal(err)
			}
			QuickLogList.InsertItem(0, timelog.Title(), timelog.Comment, 0, nil)
			break
		case 'S':
			hasRunningLog, err := constants.LogRepo.HasRunningLog(sql.NullInt32{})
			if err != nil {
				log.Fatal(err)
			}
			if !hasRunningLog {
				notification.SetNotification("There are no running quick logs")
				pagePrimitive.ShowPage("Notification")
			} else {
				SelectedIssueId.Scan(nil)
				pagePrimitive.ShowPage("AddComment")
			}
			break
		}
		return event
	})

	MainPageContainer.
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(IssueList, 0, 7, true).
			AddItem(NoteList, 0, 5, true), 0, 8, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(LogList, 0, 7, false).
			AddItem(QuickLogList, 0, 5, false), 0, 4, true)

	return MainPageContainer
}
