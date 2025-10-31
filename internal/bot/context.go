package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"
)

// BotContext holds context for bot handlers
type BotContext struct {
	Bot           *tgbotapi.BotAPI
	Update        *tgbotapi.Update
	User          *models.User
	Translator    i18n.Translator
	UserService   *service.UserService
	Config        *utils.Config
	StateManager  *state.StateManager
}

// HandlerFunc is a handler function type
type HandlerFunc func(*BotContext) error

// MiddlewareFunc is a middleware function type
type MiddlewareFunc func(*BotContext, HandlerFunc) error

