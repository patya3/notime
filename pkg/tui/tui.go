package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/models/project"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui/constants"
)

func StartTea(projectRepo project.ProjectRepo, issueRepo issue.IssueRepo, logRepo timelog.LogRepo) {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	constants.ProjectRepo = &projectRepo
	constants.IssueRepo = &issueRepo
	constants.LogRepo = &logRepo

	m := InitIssue()
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if err := constants.P.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
