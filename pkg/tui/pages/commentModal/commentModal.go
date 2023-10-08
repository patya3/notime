package commentModal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/mainPage"
	"github.com/rivo/tview"
)

var commentModal = tview.NewTextArea()
var quickLogTextModal = tview.NewTextArea()

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

// @param type = "ISSUE_LOG" | "QUICK_LOG"
func InitCommentModal(pagePrimitive *tview.Pages, modalType string) *tview.Grid {
	var textArea *tview.TextArea
	if modalType == "ISSUE_LOG" {
		textArea = commentModal
	} else {
		textArea = quickLogTextModal
	}

	var title string
	fmt.Println("modalType", modalType)
	if modalType == "ISSUE_LOG" {
		title = "Add log comment"
	} else {
		title = "Add quick log"
	}
	textArea.Box.
		SetTitle(title).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	textArea.
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if modalType == "ISSUE_LOG" {
				return inputCaptureForLogComments(pagePrimitive, event)
			} else {
				return inputCaptureForQuickLogText(pagePrimitive, event)
			}
		})

	return modal(textArea, 40, 10)
}

// AddCommentModal shows modal for adding comment to general log
func inputCaptureForLogComments(pagePrimitive *tview.Pages, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlS {
		pagePrimitive.HidePage("AddComment")
		_, err := constants.LogRepo.StopTimerByIssueId(mainPage.SelectedIssueId, commentModal.GetText())
		if err != nil {
			log.Fatal(err)
		}
		mainPage.InitLogListElements(uint(mainPage.SelectedIssueId.Int32))
		commentModal.SetText("", true)
		return nil
	}
	return event
}

// AddCommentModal shows modal for adding text for quick logs
func inputCaptureForQuickLogText(pagePrimitive *tview.Pages, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlS {
		pagePrimitive.HidePage("AddQuickLogText")

		comment := quickLogTextModal.GetText()
		timelog, err := constants.LogRepo.StartTimerForQuickLog(comment)
		if err != nil {
			log.Fatal(err)
		}

		mainPage.QuickLogList.InsertItem(0, timelog.Title(), timelog.Comment, 0, nil)
		quickLogTextModal.SetText("", true)
		return nil
	}
	return event
}
