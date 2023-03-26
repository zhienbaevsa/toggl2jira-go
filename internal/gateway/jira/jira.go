package jira

import (
	"errors"
	"time"

	jira "github.com/andygrunwald/go-jira"
)

type Client struct{}

var client *jira.Client

func New(host string, user string, pass string) (*Client, error) {
	c, err := initClient(host, user, pass)
	if err != nil {
		return &Client{}, err
	}

	client = c

	return &Client{}, nil
}

func initClient(host string, user string, pass string) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}
	client, err := jira.NewClient(tp.Client(), host)
	if err != nil {
		return &jira.Client{}, err
	}
	return client, nil
}

func (c *Client) UploadOne(issueKey string, started time.Time, timeSpentSeconds int) error {
	_, r, err := client.Issue.AddWorklogRecord(issueKey, &jira.WorklogRecord{
		Started:          (*jira.Time)(&started),
		TimeSpentSeconds: timeSpentSeconds,
	})
	if err != nil {
		return errors.New("error: " + r.Status)
	}

	return nil
}

func (c *Client) GetWorklogs(id string) ([]jira.WorklogRecord, error) {
	w, _, err := client.Issue.GetWorklogs(id)
	if err != nil {
		return []jira.WorklogRecord{}, err
	}
	return w.Worklogs, nil
}
