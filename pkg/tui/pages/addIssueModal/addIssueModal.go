package addIssueModal

import (
	// "github.com/gdamore/tcell/v2"
	"log"

	"github.com/gdamore/tcell/v2"
	issueModel "github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/tui/constants"
	"github.com/rivo/tview"
)

var addIssueForm = tview.NewForm()

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

	addIssueForm.Box.
		SetTitle("Add new issue").
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEscape:
				pagePrimitive.HidePage("AddIssue")
			}
			return event
		})

		// type Issue struct {
		// 	gorm.Model
		// 	IssueKey   string
		// 	IssueTitle string
		// 	Desc       string
		// 	ProjectID  uint
		// 	Logs       []timelog.Log
		// }
	initFormElements(addIssueForm, app, pagePrimitive)

	return modal(addIssueForm, 50, 30)
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

	addIssueForm.
		AddInputField("Issue Key", "", 0, nil, func(text string) {
			issue.IssueKey = text
		}).
		AddTextArea("Issue Title", "", 0, 2, 0, func(text string) {
			issue.IssueTitle = text
		}).
		AddTextArea("Description", "", 0, 0, 0, func(text string) {
			issue.Desc = text
		}).
		AddDropDown("Project", projectKeys, 0, func(option string, optionIndex int) {
			issue.ProjectID = projects[optionIndex].ID
		}).
		AddButton("Save", func() {
			_, err := constants.IssueRepo.CreateIssue(issue)
			if err != nil {
				log.Fatal(err)
			}
		}).
		AddButton("Quit", func() {
			pagePrimitive.HidePage("AddIssue")
		}).
		AddButton("Reset", func() {
			initFormElements(issueForm, app, pagePrimitive)
			app.SetFocus(issueForm)
		})
}
