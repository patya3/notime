package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
	"github.com/patya3/notime/pkg/models/issue"
	"github.com/patya3/notime/pkg/models/note"
	"github.com/patya3/notime/pkg/models/project"
	"github.com/patya3/notime/pkg/models/timelog"
	"github.com/patya3/notime/pkg/tui"
	"github.com/patya3/notime/pkg/tui/constants"
)

func createFileLogger(homeDir string) (logger.Interface, error) {

	file, err := os.Create(homeDir + "/notime/runtimeLogs.txt")
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	return logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer

		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	), nil

}

func getDatabaseFilePath(homeDir string) (string, error) {

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbFilePath, err := getDatabaseFilePath(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	fileLogger, err := createFileLogger(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: fileLogger,
	})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	err = db.AutoMigrate(&project.Project{}, &issue.Issue{}, &timelog.Log{}, &note.Note{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(filepath.Join(homeDir, "notime", ".env"))

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := openSqlite()
	projectRepo := project.ProjectRepo{DB: db}
	issueRepo := issue.IssueRepo{DB: db}
	logRepo := timelog.LogRepo{DB: db}
	noteRepo := note.NoteRepo{DB: db}

	constants.IssueRepo = &issueRepo
	constants.LogRepo = &logRepo
	constants.ProjectRepo = &projectRepo
	constants.NoteRepo = &noteRepo
	projectRepo.CreateProject("VM", "Virtual Microscope")
	projectRepo.CreateProject("IDV", "IDV")

	tui.StartTui()
}
