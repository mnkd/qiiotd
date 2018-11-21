package main

import (
	"fmt"
	"os"
	"time"

	slackposter "github.com/mnkd/slackposter"
)

type App struct {
	Config   Config
	QiitaAPI *QiitaAPI
	Slack    slackposter.SlackPoster
	Days     int
	YearsAgo int
}

func prepareDate(yearsAgo int, days int) (string, string) {
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

func (app *App) fetchItems(yearsAgo int, days int) ([]*QiitaItem, error) {
	t1, t2 := prepareDate(yearsAgo, days)
	items, err := app.QiitaAPI.Items(t1, t2)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *App) Run() int {
	fmt.Fprintf(os.Stdout, "fetch %v year ago...\n", app.YearsAgo)
	items, err := app.fetchItems(app.YearsAgo, app.Days)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	fmt.Printf("len: %v\n", len(items))
	for i, item := range items {
		fmt.Printf("%d: %v\n", i, item)
	}

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
		attachment := builder.BuildAttachment(*item)
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

func NewApp(config Config, yearsAgo int, days int) *App {
	return &App{
		Config:   config,
		QiitaAPI: NewQiitaAPI(config),
		Slack:    slackposter.NewSlackPoster(config.Slack),
		Days:     days,
		YearsAgo: yearsAgo,
	}
}
