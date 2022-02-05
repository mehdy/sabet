package telegram

import (
	"github.com/mehdy/sabet/pkg/jobs/meta"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Job struct {
	meta.TypeMeta `json:",inline"`
	meta.JobMeta  `json:"metadata,omitempty"`

	Spec Spec `json:"spec,omitempty"`

	bot *tgbotapi.BotAPI
}

type Spec struct {
	TokenEnv string   `json:"tokenEnv,omitempty"`
	Channel  string   `json:"channel,omitempty"`
	Template Template `json:"template,omitempty"`
}

type Template struct {
	ParseMode string `json:"parseMode,omitempty"`
	Text      string `json:"text,omitempty"`
}
