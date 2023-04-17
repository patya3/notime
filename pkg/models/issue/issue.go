package issue

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/patya3/notime/pkg/models/timelog"
)

type Issue struct {
	gorm.Model
	IssueKey   string `gorm:"unique"`
	IssueTitle string
	Desc       string
	ProjectID  uint
	Logs       []timelog.Log
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
	if err := g.DB.
		Preload("Logs", func(db *gorm.DB) *gorm.DB { return db.Order("logs.stopped_at").Limit(1) }).
		Order("created_at desc").
		Find(&issues).Error; err != nil {

		return issues, fmt.Errorf("No issues found: %v", err)
	}
	// str, _ := json.MarshalIndent(issues, "", "\t")
	// fmt.Println(string(str))
	// os.Exit(0)
	return issues, nil
}

func (g *IssueRepo) CreateIssue(issue Issue) (Issue, error) {
	if err := g.DB.Create(&issue).Error; err != nil {
		return issue, fmt.Errorf("Cannot create project: %v", err)
	}
	return issue, nil
}
