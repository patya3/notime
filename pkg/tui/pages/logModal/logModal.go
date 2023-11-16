package logModal

// TODO: edit logModal and change it to a form which has the followings
// StoppedAt => a field which show the elapsed time in the following format (HH:MM:SS)
//           => and calculate the new StoppedAt from that and CreatedAt
// Comment   string

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	c "github.com/patya3/notime/pkg/colors"
	logModel "github.com/patya3/notime/pkg/models/timelog"

	"github.com/patya3/notime/pkg/tui/constants"
	// "github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var logModal = tview.NewTextView()
var LogModalForm = tview.NewForm()

var modal = func(p tview.Primitive, width, height int) *tview.Grid {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitLogModal(app *tview.Application, pagePrimitive *tview.Pages) tview.Primitive {
	modalLocal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	LogModalForm.Box.
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' {
				pagePrimitive.HidePage("Log")
				logModal.SetText("")
				return nil
			}
			return event
		})

	// logModal.
	// 	SetTextColor(tcell.ColorIndianRed).
	// 	SetDynamicColors(true).
	// 	SetBackgroundColor(tcell.ColorDefault)

	// NOTE: not sure if this line is needed
	InitFormElements(pagePrimitive, nil)

	return modalLocal(LogModalForm, 50, 30)
}

func InitFormElements( /* app *tview.Application, */ pagePrimitive *tview.Pages, logID *uint) {
	LogModalForm.Clear(true)

	var extendedTimelog logModel.ExtendedLog
	if logID != nil {
		var err error
		extendedTimelog, err = constants.LogRepo.GetLogByID(*logID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO: fields needed
	// - issueKey (not editable)
	// - comment
	// - StoppedAt // merge the two with its difference
	// - StartedAt
	LogModalForm.
		AddTextView("Issue", "("+extendedTimelog.IssueKey+") "+extendedTimelog.IssueTitle, 0, 0, false, false).
		AddTextArea("Comment", extendedTimelog.Comment, 0, 0, 0, func(text string) {
			extendedTimelog.Comment = text
		}).
		AddInputField("Started At", extendedTimelog.CreatedAt.Format("2006-01-02 15:04:05"), 0, nil, func(text string) {}).
		AddInputField("Stopped At", extendedTimelog.StoppedAt.Time.Format("2006-01-02 15:04:05"), 0, nil, func(text string) {}).
		AddButton("Save", func() {
			fmt.Println("Save button pressed")
			// _, err := constants.LogRepo.CreateIssue(issue)
			// if err != nil {
			// 	notification.SetNotification(err.Error())
			// 	pagePrimitive.ShowPage("Notification")
			// }
			pagePrimitive.HidePage("AddIssue")
		}).
		AddButton("Quit", func() {
			fmt.Println("quit button pressed")
		})
}

func SetLogModalTextIssueLog(logID uint) {
	extendedTimelog, err := constants.LogRepo.GetLogByID(logID)
	if err != nil {
		log.Fatal(err)
	}
	logModal.SetText(fmt.Sprintf("Issue: %s%s\n\n%s\n[lightgreen]%s\n\n%s%s", extendedTimelog.IssueKey, extendedTimelog.IssueTitle, extendedTimelog.Title(), extendedTimelog.Description(), c.Colors["white"], extendedTimelog.Comment))
}

func SetLogModalTextForQuickLog(quickLogID uint) {
	quickLog, err := constants.LogRepo.GetQuickLogByID(quickLogID)
	if err != nil {
		log.Fatal(err)
	}
	logModal.SetText(fmt.Sprintf("%s\n[lightgreen]%s\n\n%s%s", quickLog.Title(), quickLog.Description(), c.Colors["white"], quickLog.Comment))
}
