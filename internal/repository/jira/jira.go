package jira

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/zhienbaevsa/toggle2jira-go/pkg/model"
)

type JiraWorklogUploader struct {
	Host     string
	User     string
	Password string
}

func (j *JiraWorklogUploader) Upload(wl []model.Worklog) error {
	tp := jira.BasicAuthTransport{
		Username: j.User,
		Password: j.Password,
	}
	client, err := jira.NewClient(tp.Client(), j.Host)
	if err != nil {
		return err
	}

	for _, v := range wl {
		client.Issue.AddWorklogRecord(v.IssueKey, &jira.WorklogRecord{
			TimeSpentSeconds: int(v.TimeSpentSeconds),
		})
	}

	return nil
}
