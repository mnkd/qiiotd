package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mnkd/slackposter"
)

type App struct {
	Config   Config
	QiitaAPI QiitaAPI
	Slack    slackposter.Slack
	Days     int
	YearsAgo int
}

func (app App) prepareDate(yearsAgo int, days int) (string, string) {
	JST, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Fprintln(os.Stderr, "App: <error>: %v\n", err)
		return "", ""
	}

	now := time.Now()
	jst := now.In(JST)

	// e.g. now: 2018-11-21, t1: 2017-11-20, t2: 2017-11-22
	t1 := jst.AddDate(-yearsAgo, 0, -1)
	t2 := jst.AddDate(-yearsAgo, 0, 1)

	return t1.Format("2006-01-02"), t2.Format("2006-01-02")
}

func (app App) fetchItems(yearsAgo int, days int) ([]QiitaItem, error) {
	t1, t2 := app.prepareDate(yearsAgo, days)
	items, err := app.QiitaAPI.Items(t1, t2)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app App) Run() int {
	fmt.Fprintf(os.Stdout, "fetch %v year ago...\n", app.YearsAgo)
	items, err := app.fetchItems(app.YearsAgo, app.Days)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	fmt.Fprintf(os.Stdout, "len: %v\n", len(items))

	var payload slackposter.Payload
	payload.Channel = app.Slack.Channel
	payload.Username = app.Slack.Username
	payload.IconEmoji = app.Slack.IconEmoji
	payload.LinkNames = true
	payload.Mrkdwn = true

	fmt.Fprintf(os.Stdout, "Build message...\n")
	builder := NewMessageBuilder(app.Config.Qiita.Domain, app.YearsAgo)

	// Prepare summary
	summary := builder.BudildSummary(len(items))
	payload.Text = summary
	fmt.Fprintf(os.Stdout, "\n%v\n", summary)

	var attachments []slackposter.Attachment
	for _, item := range items {
		fmt.Fprintln(os.Stdout, item.Title)
		attachment := builder.BuildAttachment(item)
		attachments = append(attachments, attachment)
	}
	payload.Attachments = attachments

	// Post payload
	fmt.Fprintf(os.Stdout, "Post payload...\n")
	err = app.Slack.PostPayload(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "App: <error> send a payload to slack: %v\n", err)
		return ExitCodeError
	}

	fmt.Fprintf(os.Stdout, "Post payload... Done\n")
	return ExitCodeOK
}

func NewApp(config Config, yearsAgo int, days int) App {
	var app = App{}
	app.Config = config
	app.QiitaAPI = NewQiitaAPI(config)
	app.Slack = slackposter.NewSlack(config.SlackWebhooks[0])
	app.Days = days
	app.YearsAgo = yearsAgo
	return app
}
