package toggl

import (
	"fmt"
	"regexp"
	"time"

	toggl "github.com/jason0x43/go-toggl"
	"github.com/zhienbaevsa/toggle2jira-go/pkg/model"
)

var issueKeyMap map[string]string = map[string]string{
	"misc":     "COM-1536",
	"meet":     "COM-1522",
	"personal": "personal",
}

type TogglTimeEntryStorage struct {
	ApiKey string
}

func (t *TogglTimeEntryStorage) Get(from, to time.Time) ([]model.Worklog, error) {
	s := toggl.OpenSession(t.ApiKey)
	te, err := s.GetTimeEntries(from, to)

	if err != nil {
		return []model.Worklog{}, err
	}
	var res []model.Worklog

	for _, v := range te {
		issueKey, err := getIssueKey(v)
		if err != nil {
			return []model.Worklog{}, err
		}

		res = append(res, model.Worklog{
			IssueKey:         issueKey,
			Comment:          getComment(v.Description),
			StartedAt:        *v.Start,
			TimeSpentSeconds: v.Duration,
		})
	}

	return res, nil
}

func getIssueKey(t toggl.TimeEntry) (string, error) {
	r := regexp.MustCompile("^([^:]+)")
	toggleTimeEntryKey := r.FindString(t.Description)

	if toggleTimeEntryKey == "" {
		return "", fmt.Errorf("cannot define issue key from time entry description: %v, %s", t.Description, t.Start)
	}

	k := issueKeyMap[toggleTimeEntryKey]
	if k == "" {
		return toggleTimeEntryKey, nil
	}

	return k, nil
}

func getComment(description string) string {
	r := regexp.MustCompile(":[[:space:]]*(.*)")
	return r.FindStringSubmatch(description)[1]

}
