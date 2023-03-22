package controller

import (
	"fmt"
	"time"

	"github.com/zhienbaevsa/toggle2jira-go/pkg/model"
)

type timeEntryStorage interface {
	Get(from, to time.Time) ([]model.Worklog, error)
}

type worklogUploader interface {
	Upload([]model.Worklog) error
}

type WorklogService struct {
	tes timeEntryStorage
	wlu worklogUploader
}

func NewWorklogService(tes timeEntryStorage, wlu worklogUploader) *WorklogService {
	return &WorklogService{tes, wlu}
}

func (u *WorklogService) Start(from, to time.Time) error {
	ts, err := u.tes.Get(from, to)

	if err != nil {
		return err
	}
	for _, v := range ts {
		fmt.Printf("%+v\n", v)
	}
	return nil
}
