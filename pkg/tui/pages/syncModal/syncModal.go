package syncModal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/patya3/notime/pkg/utils"
	"github.com/rivo/tview"
)

var SyncForm = tview.NewForm()

var modal = func(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func InitSyncModal(app *tview.Application, pagePrimitive *tview.Pages) tview.Primitive {

	SyncForm.Box.
		SetTitle("Add new issue").
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
				pagePrimitive.HidePage("SyncModal")
			}
			return event
		})

	initFormElements(SyncForm, app, pagePrimitive)

	return modal(SyncForm, 30, 10)
}

func initFormElements(issueForm *tview.Form, app *tview.Application, pagePrimitive *tview.Pages) {
	issueForm.Clear(true)

	projects, err := constants.ProjectRepo.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	projectKeys := make([]string, 0, len(projects))
	for _, project := range projects {
		projectKeys = append(projectKeys, project.ProjectKey)
	}

	var projectKey string

	SyncForm.
		AddDropDown("Project", projectKeys, 0, func(option string, optionIndex int) {
			if len(projects) > 0 {
				projectKey = projects[optionIndex].ProjectKey
			}
		}).
		AddButton("Sync", func() {
			pagePrimitive.HidePage("SyncModal")
			err := utils.SyncIssuesFromJira(projectKey)
			if err != nil {
				fmt.Println(err)
				notification.SetNotification("Something went wrong while syncing issues")
			} else {
				notification.SetNotification("Issues synced")
			}
			pagePrimitive.ShowPage("Notification")
		}).
		AddButton("Cancel", func() {
			pagePrimitive.HidePage("SyncModal")
		})
}
