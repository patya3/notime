package addIssueModal

import (
	"log"

	"github.com/gdamore/tcell/v2"
	issueModel "github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/patya3/notime/pkg/tui/pages/mainPage"
	"github.com/patya3/notime/pkg/tui/pages/notification"
	"github.com/rivo/tview"
)

var AddIssueForm = tview.NewForm()

func InitAddIssueForm(app *tview.Application, pagePrimitive *tview.Pages) tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	AddIssueForm.Box.
		SetTitle("Add new issue").
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
				pagePrimitive.HidePage("AddIssue")
			}
			return event
		})

	initFormElements(AddIssueForm, app, pagePrimitive)

	return modal(AddIssueForm, 50, 30)
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
	issue := issueModel.Issue{}

	AddIssueForm.
		AddDropDown("Project", projectKeys, 0, func(option string, optionIndex int) {
			if len(projects) > 0 {
				issue.ProjectID = projects[optionIndex].ID
			}
		}).
		AddInputField("Issue Key", "", 0, nil, func(text string) {
			issue.IssueKey = text
		}).
		AddTextArea("Issue Title", "", 0, 2, 0, func(text string) {
			issue.IssueTitle = text
		}).
		AddTextArea("Description", "", 0, 0, 0, func(text string) {
			issue.Desc = text
		}).
		AddButton("Save", func() {
			_, err := constants.IssueRepo.CreateIssue(issue)
			if err != nil {
				notification.SetNotification(err.Error())
				pagePrimitive.ShowPage("Notification")
			}
			mainPage.InitIssueListElements()
			pagePrimitive.HidePage("AddIssue")
		}).
		AddButton("Quit", func() {
			pagePrimitive.HidePage("AddIssue")
		}).
		AddButton("Reset", func() {
			initFormElements(issueForm, app, pagePrimitive)
			app.SetFocus(issueForm)
		})
}
