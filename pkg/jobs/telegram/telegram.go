package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mehdy/sabet/pkg/jobs/meta"
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

func (j Job) Execute(payload io.Reader) (io.Reader, error) {
	var input map[string]interface{}

	if err := json.NewDecoder(payload).Decode(&input); err != nil {
		return nil, err
	}

	if _, ok := input["items"]; !ok {
		return nil, nil
	}

	for _, item := range input["items"].([]interface{}) {
		buf := new(bytes.Buffer)

		templ := template.Must(template.New("").Parse(j.Spec.Template.Text))

		if err := templ.Execute(buf, item); err != nil {
			return nil, err
		}

		msg := tgbotapi.NewMessageToChannel(j.Spec.Channel, buf.String())
		msg.ParseMode = j.Spec.Template.ParseMode

		if _, err := j.bot.Send(msg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
