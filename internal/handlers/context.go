package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/state"
	"provbot/internal/utils"
)

// HandlerContext holds context for handlers
type HandlerContext struct {
	Bot          *tgbotapi.BotAPI
	Update       *tgbotapi.Update
	User         *models.User
	Translator   i18n.Translator
	Config       *utils.Config
	StateManager *state.StateManager
}

