package golog

import (
	"fmt"
	"log"
	"os"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type SlackWebhook struct {
	URL      string
	Username string
	Channel  string
	Icon     string
}

var Slack SlackWebhook

func New() {
	Slack = SlackWebhook{
		URL:      GetEnv("SLACK_WEBHOOK_URL"),
		Username: GetEnv("SLACK_USERNAME"),
		Channel:  GetEnv("SLACK_CHANNEL"),
	}
}

func GetEnv(key string) string {
	env := os.Getenv(key)

	if env == "" || len(env) < 1 {
		log.Fatalf("Error : %s variable not found on your system, please add to environtment variable.", key)
	}

	return env
}

func (s *SlackWebhook) sendToSlack(payload slack.Payload) {
	err := slack.Send(s.URL, "", payload)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SlackWebhook) compose(message string, messageType string, color string, emoji string, errors error) {
	attachment := slack.Attachment{
		Color: &color,
	}
	attachment.AddField(slack.Field{
		Title: "Message",
		Value: message,
	}).AddField(slack.Field{
		Title: "Level",
		Value: messageType,
	})

	if errors != nil {
		attachment.AddField(slack.Field{
			Title: "Exception",
			Value: fmt.Sprintf("``` %s ```", errors.Error()),
		})
	}

	payload := slack.Payload{
		Username:    s.Username,
		Channel:     s.Channel,
		IconEmoji:   emoji,
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
	s.sendToSlack(payload)
}

func (s *SlackWebhook) composeWithData(message string, messageType string, color string, emoji string, data []byte, e error) {
	attachment := slack.Attachment{
		Color: &color,
	}
	attachment.AddField(slack.Field{
		Title: "Message",
		Value: message,
	}).AddField(slack.Field{
		Title: "Level",
		Value: messageType,
	}).AddField(slack.Field{
		Title: "Data",
		Value: fmt.Sprintf("``` %s ```", string(data)),
	})

	if e != nil {
		attachment.AddField(slack.Field{
			Title: "Exception",
			Value: fmt.Sprintf("``` %s ```", e.Error()),
		})
	}

	payload := slack.Payload{
		Username:    s.Username,
		Channel:     s.Channel,
		IconEmoji:   emoji,
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
	s.sendToSlack(payload)
}

func (s *SlackWebhook) Info(message string) {
	s.compose(message, "INFO", "#2eb886", ":ok_hand:", nil)
}

func (s *SlackWebhook) InfoWidthData(message string, data []byte) {
	s.composeWithData(message, "INFO", "#2eb886", ":ok_hand:", data, nil)
}

func (s *SlackWebhook) Error(message string, e error) {
	s.compose(message, "ERROR", "#a30200", ":bomb:", e)
}

func (s *SlackWebhook) ErrorWithData(message string, data []byte, e error) {

	s.composeWithData(message, "ERROR", "#a30200", ":bomb:", data, e)
}

func (s *SlackWebhook) Warning(message string, e error) {
	s.compose(message, "WARNING", "#ffc107", ":warning:", e)
}

func (s *SlackWebhook) WarningWithData(message string, data []byte, e error) {
	s.composeWithData(message, "WARNING", "#ffc107", ":warning:", data, e)
}
