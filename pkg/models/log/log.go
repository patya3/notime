package log

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	StartDateTime time.Time
	EndDateTime   time.Time
	IssueID       uint
}
