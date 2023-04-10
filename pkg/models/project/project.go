package project

import (
	"fmt"

	"github.com/patya3/notime/pkg/models/issue"
	"gorm.io/gorm"
)

const (
	format string = "%d : %s\n"
)

type Project struct {
	gorm.Model
	ProjectKey string `gorm:"unique"`
	Name       string
	Issues     []issue.Issue
}

// Title the project title to display in a list
func (p Project) Title() string { return p.ProjectKey }

// Description the project description to display in a list
func (p Project) Description() string { return fmt.Sprintf("%s", p.Name) }

// FilterValue choose what field to use for filtering in a Bubbletea list component
func (p Project) FilterValue() string { return "(" + p.ProjectKey + ") " + p.Name }

type ProjectRepo struct {
	DB *gorm.DB
}

func (g *ProjectRepo) GetProjectByID(id uint) (Project, error) {
	var project Project
	if err := g.DB.Where("id = ?", id).First(&project).Error; err != nil {
		return project, fmt.Errorf("Cannot find project: %v", err)
	}
	return project, nil
}

func (g *ProjectRepo) GetAllProjects() ([]Project, error) {
	var projects []Project
	if err := g.DB.Preload("Issues.Logs").Find(&projects).Error; err != nil {
		return projects, fmt.Errorf("No projects found: %v", err)
	}
	return projects, nil
}

func (g *ProjectRepo) CreateProject(projectKey string, name string) (Project, error) {
	project := Project{ProjectKey: projectKey, Name: name}
	if err := g.DB.Create(&project).Error; err != nil {
		return project, fmt.Errorf("Cannot create project: %v", err)
	}
	return project, nil
}

func (g *ProjectRepo) RenameProject(id uint, name string) (Project, error) {
	var project Project

	if err := g.DB.Where("id = ?", id).First(&project).Error; err != nil {
		return project, fmt.Errorf("cannot update project name %v", err)
	}
	project.Name = name

	if err := g.DB.Where("id = ?", id).First(&project).Error; err != nil {
		return project, fmt.Errorf("cannot update project name %v", err)
	}

	return project, nil
}
