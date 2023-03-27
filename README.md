# Toggl2Jira

Tool for uploading work logs to Jira issues from some source of time entries. Currently, only Toggl is supported.

## Requirements

Go 1.20

## Usage

1. Clone repository
2. Create a .env file in a root directory by copying .env.template file. Fill config.
3. Run an app providing "from" and "to" arguments with a range of dates you want to upload in a "2006-01-02T15:04:05" date time format. Example:

    ```go run main.go --from=2023-03-01 --to=2023-04-01```

    If `to` date is not specified, the beginning of the current day will be used