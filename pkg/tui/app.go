package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/pages/addIssueModal"
	"github.com/patya3/notime/pkg/tui/pages/commentModal"
	"github.com/patya3/notime/pkg/tui/pages/logModal"
	"github.com/patya3/notime/pkg/tui/pages/mainPage"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pagePrimitive = tview.NewPages()

func StartTui() {

	pagePrimitive.Box.
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' && mainPage.MainPageContainer.HasFocus() {
				app.Stop()
			}
			return event
		})

	pagePrimitive.
		AddPage("Main", mainPage.InitMainPage(app, pagePrimitive), true, true).
		AddPage("AddIssue", addIssueModal.InitAddIssueForm(app, pagePrimitive), true, false).
		AddPage("Notification", notification.InitNotification(pagePrimitive), true, false).
		AddPage("AddComment", commentModal.InitCommentModal(pagePrimitive, "ISSUE_LOG"), true, false).
		AddPage("AddQuickLogText", commentModal.InitCommentModal(pagePrimitive, "QUICK_LOG"), true, false).
		AddPage("Log", logModal.InitLogModal(pagePrimitive), true, false)

	if err := app.SetRoot(pagePrimitive, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
