package mainPage

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var MainPageContainer = tview.NewFlex()

func InitMainPage(app *tview.Application, pagePrimitive *tview.Pages) *tview.Flex {
	MainPageContainer.SetBackgroundColor(tcell.ColorDefault)

	issueList := InitIssuesList(app, pagePrimitive)
	notesList := InitNotesList(app)
	logsList := InitLogsList(pagePrimitive)

	MainPageContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 49:
			app.SetFocus(issueList)
			break
		case 50:
			app.SetFocus(logsList)
			break
		case 51:
			app.SetFocus(notesList)
			break
		}
		return event
	})

	MainPageContainer.
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(issueList, 0, 7, true).
			AddItem(notesList, 0, 5, true), 0, 8, true).
		AddItem(logsList, 0, 4, false)

	return MainPageContainer
}
