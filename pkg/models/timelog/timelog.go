package timelog

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	StoppedAt sql.NullTime
	IssueID   uint
}

func (l *Log) Title() string {
	// return fmt.Sprint(diff.Format("15:04:05"))
	if l.StoppedAt.Valid {
		diff := time.Time{}.Add(l.StoppedAt.Time.Sub(l.CreatedAt))
		return fmt.Sprintf("%d", l.ID) + " - Duration: " + diff.Format("15:04:05")
	}
	return fmt.Sprintf("%d", l.ID) + " - Running"
}

func (l *Log) Description() string {
	if l.StoppedAt.Valid {
		return l.CreatedAt.Format("2006. 01. 02 --> 15:04") + " - " + l.StoppedAt.Time.Format("15:04")
	}
	return l.CreatedAt.Format("2006-01-02 15:04")
}

func (l *Log) FilterValue() string { return l.Description() }

type LogRepo struct {
	DB *gorm.DB
}

// create log row
func (g *LogRepo) StartTimer(issueID uint) (Log, error) {
	log := Log{IssueID: issueID}
	if err := g.DB.Create(&log).Error; err != nil {
		return log, fmt.Errorf("Cannot start timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) StopTimer(logID uint) (Log, error) {
	var log Log
	err := g.DB.Model(&log).Where("id = ?", logID).Update("StoppedAt", time.Now()).Error
	if err != nil {
		return log, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) GetAllLogsForAnIssue(issueID uint) ([]Log, error) {
	var logs []Log
	if err := g.DB.Where("issue_id = ?", issueID).Find(&logs).Error; err != nil {
		return logs, fmt.Errorf("No logs found: %v", err)
	}
	return logs, nil
}
