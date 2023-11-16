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
	IssueID   sql.NullInt32 `sql:"DEFAULT:NULL"`
	Comment   string
	Logged    bool
}

type ExtendedLog struct {
	Log
	IssueTitle string
	IssueKey   string
}

func (l *Log) Title() string {
	// return fmt.Sprint(diff.Format("15:04:05"))
	if l.StoppedAt.Valid {
		diff := time.Time{}.Add(l.StoppedAt.Time.Sub(l.CreatedAt))
		return fmt.Sprintf("[red](%s%s[red])[pink] - Duration: %s%s",
			c.Colors["lightorange"],
			fmt.Sprintf("%d", l.ID),
			c.Colors["lightpurple"],
			diff.Format("15:04:05"),
		)
	}
	currentTime := time.Time{}.Add(time.Now().Sub(l.CreatedAt))
	return fmt.Sprintf("%s(%s%s%s)[pink] - Running: %s%s",
		c.Colors["darkgreen"],
		c.Colors["green"],
		fmt.Sprintf("%d", l.ID),
		c.Colors["darkgreen"],
		c.Colors["lightpurple"],
		currentTime.Format("15:04:05"),
	)
}

func (l *Log) Description() string {
	if l.StoppedAt.Valid {
		return l.CreatedAt.Format("[lightgreen]2006. 01. 02 --> 15:04") + " - " + l.StoppedAt.Time.Format("15:04")
	}
	return l.CreatedAt.Format("[lightgreen]2006-01-02 15:04")
}

func (l *Log) FilterValue() string { return l.Description() }

type LogRepo struct {
	DB *gorm.DB
}

// create log row
func (g *LogRepo) StartTimer(issueID sql.NullInt32) (Log, error) {
	log := Log{IssueID: issueID}
	if err := g.DB.Create(&log).Error; err != nil {
		return log, fmt.Errorf("Cannot start timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) StartTimerForQuickLog(comment string) (Log, error) {
	var log Log
	log = Log{Comment: comment}
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

func (g *LogRepo) StopTimerForQuickLog() (Log, error) {
	var log Log
	if err := g.DB.Model(&log).Where("issue_id IS NULL AND stopped_at IS NULL").Update("StoppedAt", time.Now()).Error; err != nil {
		return log, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) CopyTimerByLogAndIssueId(logID uint, issueID uint) (Log, error) {
	var oldLog Log
	err := g.DB.Model(&oldLog).Where("id = ?", logID).Error
	if err != nil {
		return oldLog, fmt.Errorf("Cannot continue timer (log not found): %v", err)
	}

	newLog := Log{IssueID: sql.NullInt32{
		Int32: int32(issueID), Valid: true},
		Comment: oldLog.Comment,
		Logged:  oldLog.Logged,
	}
	if err := g.DB.Create(&newLog).Error; err != nil {
		return newLog, fmt.Errorf("Cannot create copy log: %v", err)
	}

	return newLog, nil
}

func (g *LogRepo) StopTimerByIssueId(issueID sql.NullInt32, comment string) (Log, error) {
	var log Log
	query := g.DB.Model(&log).Where("issue_id = ? AND stopped_at IS NULL", issueID)
	if !issueID.Valid {
		query = g.DB.Model(&log).Where("issue_id IS NULL AND stopped_at IS NULL")
	}
	err := query.Updates(map[string]interface{}{"stopped_at": time.Now(), "comment": comment}).Error
	if err != nil {
		return log, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return log, nil
}

func (g *LogRepo) HasRunningLog(issueID sql.NullInt32) (bool, error) {
	var count int64
	var log *Log
	query := g.DB.Model(&log).Where("issue_id = ? AND stopped_at IS NULL", issueID).Count(&count)
	if !issueID.Valid {
		query = g.DB.Model(&log).Where("issue_id IS NULL AND stopped_at IS NULL").Count(&count)
	}
	err := query.Error
	if err != nil {
		return false, fmt.Errorf("Cannot stop timer: %v", err)
	}
	return count > 0, nil
}

func (g *LogRepo) HasRunningLogQuickLog() (bool, error) {
	var count int64
	var log Log

	query := g.DB.Model(&log).Where("issue_id IS NULL AND stopped_at IS NULL").Count(&count)
	if query.Error != nil {
		return false, fmt.Errorf("Query failed: %v", query.Error)
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

func (g *LogRepo) GetAllQuickLogs() ([]Log, error) {
	var logs []Log
	if err := g.DB.Where("issue_id IS NULL").Order("created_at desc").Find(&logs).Error; err != nil {
		return logs, fmt.Errorf("No logs found: %v", err)
	}
	return logs, nil
}

func (g *LogRepo) GetQuickLogByID(id uint) (Log, error) {
	var log Log
	if err := g.DB.Model(&log).Where("id = ?", id).First(&log).Error; err != nil {
		// NOTE: maybe return with null pointer would be better
		return log, fmt.Errorf("No log found with the given ID: %v", err)
	}
	return log, nil
}

func (g *LogRepo) GetLogByID(id uint) (ExtendedLog, error) {
	var log ExtendedLog
	var l Log
	if err := g.DB.Model(&l).
		Select("logs.*, issues.issue_title, issues.issue_key").
		Joins("INNER JOIN issues ON logs.issue_id = issues.id").
		Where("logs.id = ?", id).
		First(&log).
		Error; err != nil {

		return log, fmt.Errorf("No log found with the given ID: %v", err)
	}
	return log, nil
}
