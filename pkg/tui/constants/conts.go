package constants

import (
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/models/note"
	"github.com/patya3/notime/pkg/models/project"
	"github.com/patya3/notime/pkg/models/timelog"
)

var (
	IssueRepo   *issue.IssueRepo
	LogRepo     *timelog.LogRepo
	ProjectRepo *project.ProjectRepo
	NoteRepo    *note.NoteRepo
)
