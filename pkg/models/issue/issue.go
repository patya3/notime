package issue

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/patya3/notime/pkg/models/log"
)

type Issue struct {
	gorm.Model
	IssueKey   string
	IssueTitle string
	Desc       string
	ProjectID  uint
	Logs       []log.Log
}

func (i Issue) Title() string { return i.IssueKey }

func (i Issue) Description() string { return fmt.Sprintf("%s", i.IssueTitle) }

func (i Issue) FilterValue() string { return "(" + i.IssueKey + ") " + i.IssueTitle }

type IssueRepo struct {
	DB *gorm.DB
}

func (g *IssueRepo) GetIssueByID(id uint) (Issue, error) {
	var issue Issue
	if err := g.DB.Where("id = ?", id).First(&issue).Error; err != nil {
		return issue, fmt.Errorf("Cannot find issue: %v", err)
	}
	return issue, nil
}

func (g *IssueRepo) GetAllIssues() ([]Issue, error) {
	var issues []Issue
	if err := g.DB.Find(&issues).Error; err != nil {
		return issues, fmt.Errorf("No issues found: %v", err)
	}
	return issues, nil
}

func (g *IssueRepo) CreateIssue(issueKey string, issueTitle string, desc string, projectID uint) (Issue, error) {
	issue := Issue{IssueKey: issueKey, IssueTitle: issueTitle, Desc: desc, ProjectID: projectID}
	if err := g.DB.Create(&issue).Error; err != nil {
		return issue, fmt.Errorf("Cannot create project: %v", err)
	}
	return issue, nil
}
