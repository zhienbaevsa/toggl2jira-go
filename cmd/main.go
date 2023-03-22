package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zhienbaevsa/toggle2jira-go/internal/controller"
	"github.com/zhienbaevsa/toggle2jira-go/internal/repository/jira"
	"github.com/zhienbaevsa/toggle2jira-go/internal/repository/toggl"
)

const timeFormat string = "2006-01-02"

type config struct {
	Toggl struct {
		host   string
		apiKey string
	}
	Jira struct {
		host string
		user string
		pass string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.Toggl.host = os.Getenv("TOGGL_HOST")
	cfg.Toggl.apiKey = os.Getenv("TOGGL_API_KEY")

	cfg.Jira.host = os.Getenv("JIRA_HOST")
	cfg.Jira.user = os.Getenv("JIRA_USER")
	cfg.Jira.pass = os.Getenv("JIRA_PASS")

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	tes := &toggl.TogglTimeEntryStorage{
		ApiKey: cfg.Toggl.apiKey,
	}
	wlu := &jira.JiraWorklogUploader{
		Host:     cfg.Jira.host,
		User:     cfg.Jira.user,
		Password: cfg.Jira.pass,
	}
	wls := controller.NewWorklogService(tes, wlu)

	from := flag.String("from", "", "From date")
	to := flag.String("to", "", "To date")

	flag.Parse()

	fromDate, err := time.Parse(timeFormat, *from)
	if err != nil {
		panic(fmt.Sprintf("Could not parse from time: %v", *from))
	}

	var toDate time.Time
	if *to != "" {
		toDateParsed, err := time.Parse(timeFormat, *to)
		if err != nil {
			panic(fmt.Sprintf("Could not parse to time: %v", *to))
		}
		toDate = toDateParsed
	} else {
		now := time.Now()
		toDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	err = wls.Start(fromDate, toDate)
	if err != nil {
		panic(err) // TODO Proper error handling
	}

	// worklogs := 10
	// bar := progressbar.Default(int64(worklogs))

	// for i := 0; i < worklogs; i++ {
	// 	bar.Add(1)
	// }
}
