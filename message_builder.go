package main

import (
	"fmt"

	"github.com/mnkd/slackposter"
)

type MessageBuilder struct {
	QiitaDomain string
	YearsAgo    int
}

func (builder MessageBuilder) BudildSummary(itemsCount int) string {
	repo := "Qiita:Team"
	url := "https://" + builder.QiitaDomain
	link := fmt.Sprintf("<%s|%s>", url, repo)

	var summary string
	switch itemsCount {
	case 0:
		summary = fmt.Sprintf("%d 年前の今日の投稿はありません。%s", builder.YearsAgo, link)
	default:
		summary = fmt.Sprintf("%d 年前の今日の投稿が %d 件みつかりました %s", builder.YearsAgo, itemsCount, link)
	}
	return summary
}

func (builder MessageBuilder) fallbackString(item QiitaItem) string {
	return fmt.Sprintf("<%s|%s>\nby %s", item.URL, item.Title, item.User.ID)
}

func (builder MessageBuilder) BuildAttachment(item QiitaItem) slackposter.Attachment {
	color := "#3287C8"
	message := builder.fallbackString(item)

	var timestamp int64
	time, err := item.Time_CreatedAt()
	if err != nil {
		// Do nothing
	}
	timestamp = time.Unix()

	var attachment slackposter.Attachment
	attachment = slackposter.Attachment{
		Fallback:  message,
		Color:     color,
		Title:     item.Title,
		TitleLink: item.URL,
		Footer:    item.User.ID,
		Ts:        timestamp,
		ThumbUrl:  item.User.ProfileImageURL,
		MrkdwnIn:  []string{"text", "fallback"},
	}
	return attachment
}

func NewMessageBuilder(domain string, yearsAgo int) MessageBuilder {
	return MessageBuilder{
		QiitaDomain: domain,
		YearsAgo:    yearsAgo,
	}
}
