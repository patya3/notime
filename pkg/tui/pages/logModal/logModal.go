package logModal

// TODO: edit logModal and change it to a form which has the followings
// StoppedAt => a field which show the elapsed time in the following format (HH:MM:SS)
//           => and calculate the new StoppedAt from that and CreatedAt
// Comment   string

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	c "github.com/patya3/notime/pkg/colors"
	logModel "github.com/patya3/notime/pkg/models/timelog"

	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/notification"

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
	// InitFormElements(app, pagePrimitive,  nil)

	return modalLocal(LogModalForm, 50, 30)
}

// @param logType = "ISSUE_LOG" | "QUICK_LOG"
func InitFormElements(app *tview.Application, pagePrimitive *tview.Pages, logType string, logID *uint) {
	LogModalForm.Clear(true)

	var extendedTimelog logModel.ExtendedLog
	if logType == "ISSUE_LOG" {
		if logID != nil {
			var err error
			extendedTimelog, err = constants.LogRepo.GetLogByID(*logID)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("no logID provided")
		}
		LogModalForm.
			AddTextView("Issue", "("+extendedTimelog.IssueKey+") "+extendedTimelog.IssueTitle, 0, 0, false, false)
	} else if logType == "QUICK_LOG" {
		quickLog, err := constants.LogRepo.GetQuickLogByID(*logID)
		if err != nil {
			log.Fatal(err)
		}
		extendedTimelog = logModel.ExtendedLog{Log: quickLog, IssueKey: "", IssueTitle: ""}

	}

	LogModalForm.
		AddTextArea("Comment", extendedTimelog.Comment, 0, 0, 0, func(text string) {
			extendedTimelog.Comment = text
		}).
		AddInputField("Started At", extendedTimelog.CreatedAt.Format("2006-01-02 15:04:05"), 0, nil, func(text string) {
			date, error := time.Parse("2006-01-02 15:04:05", text)
			if error != nil {
				LogModalForm.GetButton(0).SetDisabled(true)
			} else {
				LogModalForm.GetButton(0).SetDisabled(false)
			}
			extendedTimelog.CreatedAt = date
		}).
		AddInputField("Stopped At", extendedTimelog.StoppedAt.Time.Format("2006-01-02 15:04:05"), 0, nil, func(text string) {
			dateTime, error := time.Parse("2006-01-02 15:04:05", text)
			if error != nil {
				LogModalForm.GetButton(0).SetDisabled(true)
			} else {
				LogModalForm.GetButton(0).SetDisabled(false)
			}
			extendedTimelog.StoppedAt = sql.NullTime{Time: dateTime, Valid: true}
		}).
		AddButton("Save", func() {
			_, err := constants.LogRepo.EditLog(extendedTimelog.Log)
			if err != nil {
				notification.SetNotification(err.Error())
				pagePrimitive.ShowPage("Notification")
			}

			pagePrimitive.HidePage("Log")
		}).
		AddButton("Reset", func() {
			InitFormElements(app, pagePrimitive, logType, logID)
			app.SetFocus(LogModalForm)
		}).
		AddButton("Quit", func() {
			pagePrimitive.HidePage("Log")
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
