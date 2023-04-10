package commentModal

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/mainPage"
	"github.com/rivo/tview"
)

var commentModal = tview.NewTextArea()

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitCommentModal(pagePrimitive *tview.Pages) *tview.Grid {
	commentModal.Box.
		SetTitle("Add log commment").
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	commentModal.
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlS {
				pagePrimitive.HidePage("AddComment")
				_, err := constants.LogRepo.StopTimerByIssueId(mainPage.SelectedIssueId, commentModal.GetText())
				if err != nil {
					log.Fatal(err)
				}
				mainPage.InitLogListElements(mainPage.SelectedIssueId)
				commentModal.SetText("", true)
				return nil
			}
			return event
		})

	return modal(commentModal, 40, 10)
}
