package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/models/project"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui"
	"github.com/patya3/notime/pkg/tui/constants"
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second,   // Slow SQL threshold
		LogLevel:                  logger.Silent, // Log level
		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,          // Disable color
	},
)

func getDatabaseFilePath() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(homeDir + "/notime"); os.IsNotExist(err) {
		if err := os.MkdirAll(homeDir+"/notime", os.ModePerm); err != nil {
			return "", err
		}
	}

	file := homeDir + "/notime/logs.db"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if _, err := os.Create(file); err != nil {
			return "", err
		}
	}

	return file, nil
}

func openSqlite() *gorm.DB {
	dbFilePath, err := getDatabaseFilePath()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	err = db.AutoMigrate(&project.Project{}, &issue.Issue{}, &timelog.Log{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := openSqlite()
	projectRepo := project.ProjectRepo{DB: db}
	issueRepo := issue.IssueRepo{DB: db}
	logRepo := timelog.LogRepo{DB: db}

	constants.IssueRepo = &issueRepo
	constants.LogRepo = &logRepo
	constants.ProjectRepo = &projectRepo
	projectRepo.CreateProject("VM", "Virtual Microscope")
	projectRepo.CreateProject("IDV", "IDV")
	tui.StartTui()
}
