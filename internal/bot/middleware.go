package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/utils"
)

// MessageLoggingMiddleware logs all incoming messages
func MessageLoggingMiddleware(logRepo *repository.LogRepository) MiddlewareFunc {
	return func(ctx *BotContext, next HandlerFunc) error {
		if ctx.Update.Message != nil {
			userID := (*int)(nil)
			if ctx.User != nil {
				userID = &ctx.User.ID
			}
			
			messageText := ctx.Update.Message.Text
			log := &models.MessageLog{
				UserID:      userID,
				TelegramID:  ctx.Update.Message.From.ID,
				Direction:   models.DirectionIncoming,
				MessageText: &messageText,
			}
			
			if ctx.Update.Message.MessageID != 0 {
				messageID := int64(ctx.Update.Message.MessageID)
				log.MessageID = &messageID
			}
			
			if err := logRepo.LogMessage(context.Background(), log); err != nil {
				utils.Logger.WithError(err).Error("Failed to log incoming message")
			}
		}
		return next(ctx)
	}
}

// UserRegistrationMiddleware ensures user is registered
func UserRegistrationMiddleware(userService *service.UserService) MiddlewareFunc {
	return func(ctx *BotContext, next HandlerFunc) error {
		if ctx.Update.Message == nil && ctx.Update.CallbackQuery == nil {
			return next(ctx)
		}

		var from *tgbotapi.User
		if ctx.Update.Message != nil {
			from = ctx.Update.Message.From
		} else if ctx.Update.CallbackQuery != nil {
			from = ctx.Update.CallbackQuery.From
		}

		if from == nil {
			return next(ctx)
		}

		// Register or update user
		var username *string
		if from.UserName != "" {
			username = &from.UserName
		}
		
		firstName := &from.FirstName
		lastName := &from.LastName
		if from.LastName == "" {
			lastName = nil
		}

		user, err := userService.RegisterOrUpdateUser(
			context.Background(),
			int64(from.ID),
			username,
			firstName,
			lastName,
		)
		if err != nil {
			utils.Logger.WithError(err).Error("Failed to register/update user")
			return fmt.Errorf("failed to register user: %w", err)
		}

		ctx.User = user
		
		// Set translator based on user language
		lang := user.Language
		if lang == "" {
			lang = i18n.DefaultLanguage
		}
		ctx.Translator = i18n.GetGlobalTranslator(lang)

		return next(ctx)
	}
}

// AdminOnlyMiddleware restricts access to admin users
func AdminOnlyMiddleware(config *utils.Config) MiddlewareFunc {
	return func(ctx *BotContext, next HandlerFunc) error {
		if ctx.User == nil {
			return fmt.Errorf("user not registered")
		}

		if !config.IsAdmin(ctx.User.TelegramID) {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_unauthorized"))
			_, _ = ctx.Bot.Send(msg)
			return fmt.Errorf("unauthorized access attempt")
		}

		return next(ctx)
	}
}

