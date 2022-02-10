package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/sirupsen/logrus"
)

func (j Job) GetStoreType() string {
	return ""
}

func (j Job) SetStore(_ meta.Store) {
}

func (j *Job) Init() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(j.Spec.TokenEnv))
	if err != nil {
		return err
	}

	j.bot = bot

	return nil
}

func (j *Job) log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"type": j.GetType(),
		"name": j.GetName(),
	})
}

func (j Job) Execute(payload io.Reader) (io.Reader, error) {
	var input map[string]interface{}

	if err := json.NewDecoder(payload).Decode(&input); err != nil {
		return nil, err
	}

	if _, ok := input["items"]; !ok {
		return nil, fmt.Errorf("items not found in input")
	}

	items, ok := input["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("items is not an array")
	}

	itemsCount := len(items)
	sentCounter := 0

	for _, item := range items {
		buf := new(bytes.Buffer)

		templ := template.Must(template.New("").Parse(j.Spec.Template.Text))

		if err := templ.Execute(buf, item); err != nil {
			j.log().WithError(err).Error("Failed to execute template")
			continue
		}

		msg := tgbotapi.NewMessageToChannel(j.Spec.Channel, buf.String())
		msg.ParseMode = j.Spec.Template.ParseMode

		if _, err := j.bot.Send(msg); err != nil {
			j.log().WithError(err).Error("Failed to send message")
			continue
		}

		sentCounter++
	}

	j.log().WithField("items", itemsCount).WithField("sent", sentCounter).Info("Message sent")

	return nil, nil
}
