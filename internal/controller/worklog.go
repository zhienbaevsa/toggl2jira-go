package controller

import (
	"fmt"
	"strings"
	"time"

	event "github.com/gookit/event"
	"github.com/zhienbaevsa/toggle2jira-go/internal/gateway/jira"
	"github.com/zhienbaevsa/toggle2jira-go/pkg/model"
)

type timeEntryStorage interface {
	Get(from, to time.Time) ([]model.Worklog, error)
}

type WorklogUploader struct {
	timeEntryStorage
	jira jira.Client
}

type JiraConfig struct {
	Host string
	User string
	Pass string
}

var issuesWorklogsMap map[string]bool

var loc *time.Location

const WorklogUploadedEvent = "worklog.uploaded"

const timeFormat = "2006-01-02T15:04:05"

func New(ts timeEntryStorage, tz string, jc JiraConfig) (*WorklogUploader, error) {
	l, err := time.LoadLocation(tz)
	if err != nil {
		return &WorklogUploader{}, err
	}
	loc = l

	j, err := jira.New(jc.Host, jc.User, jc.Pass)
	if err != nil {
		return &WorklogUploader{}, err
	}

	return &WorklogUploader{ts, *j}, nil
}

func (u *WorklogUploader) Start(from, to time.Time) error {
	ww, err := u.timeEntryStorage.Get(from, to)
	if err != nil {
		return err
	}

	ik := []string{}
	for _, v := range ww {
		ik = append(ik, v.IssueKey)
	}
	err = u.loadIssuesWorklogsMap(ik)
	if err != nil {
		return err
	}

	for _, v := range ww {
		err = u.UploadOne(v)
		if err != nil {
			return err
		}
		event.MustFire(WorklogUploadedEvent, nil)
	}
	return nil
}

func (u *WorklogUploader) loadIssuesWorklogsMap(ids []string) error {
	res := make(map[string]bool)
	for _, v := range ids {
		ww, err := u.jira.GetWorklogs(v)
		if err != nil {
			return err
		}

		for _, w := range ww {
			k := getIssueWorklogMapKey(v, time.Time(*w.Started))
			res[k] = true
		}
	}
	issuesWorklogsMap = res
	return nil
}

func (u *WorklogUploader) UploadOne(w model.Worklog) error {
	k := getIssueWorklogMapKey(w.IssueKey, w.StartedAt)

	if _, exists := issuesWorklogsMap[k]; exists {
		return nil
	}

	s := w.StartedAt.In(loc)
	ts := int(w.TimeSpentSeconds)
	if ts < 60 {
		ts = 60
	}
	err := u.jira.UploadOne(w.IssueKey, s, ts)
	if err != nil {
		return fmt.Errorf("cannot upload worklog %+v: %v", w, err)
	}

	return nil
}

func getIssueWorklogMapKey(issueKey string, time time.Time) string {
	return fmt.Sprintf("%v-%s", strings.ToLower(issueKey), time.In(loc).Format(timeFormat))
}
