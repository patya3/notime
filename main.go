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

func openSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("logs.db"), &gorm.Config{
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
