package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gookit/event"
	"github.com/schollz/progressbar/v3"
	"github.com/zhienbaevsa/toggle2jira-go/internal/controller"
	"github.com/zhienbaevsa/toggle2jira-go/internal/repository/toggl"
	"gopkg.in/yaml.v3"
)

const timeFormat string = "2006-01-02"
const configFile string = "config.yaml"

type config struct {
	ToggleApiKey       string            `yaml:"toggle_api_key"`
	JiraHost           string            `yaml:"jira_host"`
	JiraUser           string            `yaml:"jira_user"`
	JiraPass           string            `yaml:"jira_pass"`
	Timezone           string            `yaml:"timezone"`
	AliasToIssueKeyMap map[string]string `yaml:"alias_to_issue_key_map"`
}

func main() {
	fromDate, toDate := mustFromAndToDatesFromArgs()

	cfg, err := loadConfig()

	if err != nil {
		panic(fmt.Sprintf("error while loading config: %v", err))
	}

	wu := mustWorklogUploader(cfg)

	bar := progressbar.Default(-1, "Uploading worklogs...")
	event.Listen(controller.WorklogUploadedEvent, event.ListenerFunc(func(e event.Event) error {
		bar.Add(1)
		return nil
	}), event.Normal)

	err = wu.Start(fromDate, toDate)
	if err != nil {
		panic(fmt.Sprintf("error while loading worklogs: %v", err))
	}
}

func loadConfig() (config, error) {
	var cfg config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func mustFromAndToDatesFromArgs() (time.Time, time.Time) {
	from := flag.String("from", "", "From date")
	to := flag.String("to", "", "To date")

	flag.Parse()

	fromDate, err := time.Parse(timeFormat, *from)
	if err != nil {
		panic(fmt.Sprintf("Could not parse from date: %v", *from))
	}

	var toDate time.Time
	if *to != "" {
		toDateParsed, err := time.Parse(timeFormat, *to)
		if err != nil {
			panic(fmt.Sprintf("could not parse to date: %v", *to))
		}
		toDate = toDateParsed
	} else {
		// If to date is not specified, set it to the beginning of the current day
		now := time.Now()
		toDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	return fromDate, toDate
}

func mustWorklogUploader(cfg config) controller.WorklogUploader {
	ts := &toggl.TogglTimeEntryStorage{
		ApiKey:             cfg.ToggleApiKey,
		AliasToIssueKeyMap: cfg.AliasToIssueKeyMap,
	}
	wu, err := controller.New(ts, cfg.Timezone, controller.JiraConfig{
		Host: cfg.JiraHost,
		User: cfg.JiraUser,
		Pass: cfg.JiraPass,
	})
	if err != nil {
		panic(fmt.Sprintf("error while creating WorklogUploader: %v", err))
	}
	return *wu
}
