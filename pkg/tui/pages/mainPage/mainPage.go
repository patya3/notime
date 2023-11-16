package mainPage

import (
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
	InitNotesList(app, pagePrimitive)
	InitLogList(LogList, "ISSUE_LOG", app, pagePrimitive)
	InitLogList(QuickLogList, "QUICK_LOG", app, pagePrimitive)
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
			hasRunningQuickLog, err := constants.LogRepo.HasRunningLogQuickLog()
			if err != nil {
				log.Fatal(err)
			}
			if hasRunningQuickLog {
				notification.SetNotification("There is a running quick log")
				pagePrimitive.ShowPage("Notification")
			} else {
				pagePrimitive.ShowPage("AddQuickLogText")
			}
			break
		case 'S':
			hasRunningQuickLog, err := constants.LogRepo.HasRunningLogQuickLog()
			if err != nil {
				log.Fatal(err)
			}
			if !hasRunningQuickLog {
				notification.SetNotification("There are no running quick logs")
				pagePrimitive.ShowPage("Notification")
			} else {
				// TODO: maybe need a confirmation and a comment rewrite if needed (load old comment from database)
				_, err := constants.LogRepo.StopTimerForQuickLog()
				if err != nil {
					log.Fatal(err)
				}
				InitQuickLogListElements()
			}
			break
		}
		return event
	})

	MainPageContainer.
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(IssueList, 0, 7, true).
			AddItem(NoteList, 0, 5, true), 0, 7, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(LogList, 0, 7, false).
			AddItem(QuickLogList, 0, 5, false), 0, 5, true)

	return MainPageContainer
}
