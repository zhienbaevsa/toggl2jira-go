package model

import "time"

type Worklog struct {
	IssueKey         string
	Comment          string
	StartedAt        time.Time
	TimeSpentSeconds int64
}
