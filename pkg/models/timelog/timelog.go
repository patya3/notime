package timelog

import (
	"database/sql"
	"fmt"
	"time"

	c "github.com/patya3/notime/pkg/colors"
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	StoppedAt sql.NullTime
	IssueID   uint
	Comment   string
	Logged    bool
}

func (l *Log) Title(colors bool) string {
	// return fmt.Sprint(diff.Format("15:04:05"))
	if l.StoppedAt.Valid {
		diff := time.Time{}.Add(l.StoppedAt.Time.Sub(l.CreatedAt))
		if colors {
			return "[red](" + c.Colors["lightorange"] + fmt.Sprintf("%d", l.ID) + "[red])[pink] - Duration: " + c.Colors["lightpurple"] + diff.Format("15:04:05")
		}
		return "(" + fmt.Sprintf("%d", l.ID) + ") - Duration: " + diff.Format("15:04:05")
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

func (g *LogRepo) StopTimerByLogId(logID uint) (Log, error) {
	var log Log
	err := g.DB.Model(&log).Where("id = ?", logID).Update("StoppedAt", time.Now()).Error
	if err != nil {
		return log, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) StopTimerByIssueId(issueID uint, comment string) (Log, error) {
	var log Log
	err := g.DB.Model(&log).
		Where("issue_id = ? AND stopped_at IS NULL", issueID).
		Updates(map[string]interface{}{"stopped_at": time.Now(), "comment": comment}).
		Error
	if err != nil {
		return log, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) HasRunningLog(issueID uint) (bool, error) {
	var count int64
	var log Log
	err := g.DB.Model(&log).Where("issue_id = ? AND stopped_at IS NULL", issueID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return count > 0, nil
}

func (g *LogRepo) GetAllLogsForAnIssue(issueID uint) ([]Log, error) {
	var logs []Log
	if err := g.DB.Where("issue_id = ?", issueID).Order("created_at desc").Find(&logs).Error; err != nil {
		return logs, fmt.Errorf("No logs found: %v", err)
	}
	return logs, nil
}

func (g *LogRepo) GetLogByID(id uint) (Log, error) {
	var log Log
	if err := g.DB.Where("id = ?", id).First(&log).Error; err != nil {
		return log, fmt.Errorf("No log found with the given ID: %v", err)
	}
	return log, nil
}
