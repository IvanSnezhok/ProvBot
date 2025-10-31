package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"provbot/internal/handlers"
	"provbot/internal/handlers/admin"
	"provbot/internal/handlers/support"
	"provbot/internal/handlers/user"
	"provbot/internal/repository"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"
)

// BotHandler handles bot updates
type BotHandler struct {
	bot              *tgbotapi.BotAPI
	config           *utils.Config
	userService      *service.UserService
	billingService   *service.BillingService
	supportService   *service.SupportService
	adminService     *service.AdminService
	notificationService *service.NotificationService
	stateManager     *state.StateManager
	userRepo         *repository.UserRepository
	billingRepo      *repository.BillingRepository
	logRepo          *repository.LogRepository
	
	// User handlers
	startHandler     *user.StartHandler
	payBillHandler   *user.PayBillHandler
	timePayHandler   *user.TimePayHandler
	supportChatHandler *support.SupportChatHandler
	
	// Admin handlers
	adminPanelHandler *admin.PanelHandler
	
	handlers         map[string]HandlerFunc
	middlewares      []MiddlewareFunc
}

// NewBotHandler creates a new bot handler
func NewBotHandler(
	bot *tgbotapi.BotAPI,
	config *utils.Config,
	userService *service.UserService,
	billingService *service.BillingService,
	supportService *service.SupportService,
	adminService *service.AdminService,
	notificationService *service.NotificationService,
	stateManager *state.StateManager,
	userRepo *repository.UserRepository,
	billingRepo *repository.BillingRepository,
	logRepo *repository.LogRepository,
) *BotHandler {
	h := &BotHandler{
		bot:                bot,
		config:             config,
		userService:        userService,
		billingService:     billingService,
		supportService:     supportService,
		adminService:       adminService,
		notificationService: notificationService,
		stateManager:       stateManager,
		userRepo:           userRepo,
		billingRepo:        billingRepo,
		logRepo:            logRepo,
		handlers:           make(map[string]HandlerFunc),
		middlewares:        []MiddlewareFunc{},
	}

	// Initialize user handlers
	h.startHandler = user.NewStartHandler(userService, billingRepo, userRepo, stateManager, config)
	h.payBillHandler = user.NewPayBillHandler(billingService, billingRepo, userRepo, stateManager, config)
	h.timePayHandler = user.NewTimePayHandler(billingRepo, userRepo, config)
	h.supportChatHandler = support.NewSupportChatHandler(logRepo, stateManager, config, bot)
	
	// Initialize admin handlers
	h.adminPanelHandler = admin.NewPanelHandler(adminService, billingRepo, userRepo, stateManager, config)

	h.registerHandlers()
	return h
}

// Use adds middleware to the handler chain
func (h *BotHandler) Use(middleware MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middleware)
}

// registerHandlers registers all command handlers
func (h *BotHandler) registerHandlers() {
	h.handlers["/start"] = h.handleStart
	h.handlers["/help"] = h.handleHelp
	h.handlers["/profile"] = h.handleProfile
	h.handlers["/balance"] = h.handleBalance
	h.handlers["/services"] = h.handleServices
	h.handlers["/support"] = h.handleSupport
	h.handlers["/language"] = h.handleLanguage
	
	// Admin commands
	h.handlers["/admin"] = h.handleAdmin
	h.handlers["/broadcast"] = h.handleBroadcast
	h.handlers["/outage"] = h.handleOutage
	h.handlers["/users"] = h.handleUsers
	h.handlers["/connect"] = h.handleConnect
	h.handlers["/end_chat"] = h.handleEndChat
	h.handlers["/stats"] = h.handleStats
}

// HandleUpdate processes an update
func (h *BotHandler) HandleUpdate(update tgbotapi.Update) {
	ctx := &BotContext{
		Bot:          h.bot,
		Update:       &update,
		Config:       h.config,
		UserService:  h.userService,
		StateManager: h.stateManager,
	}

	// Apply middlewares
	handler := h.getHandler(update)
	if handler == nil {
		return
	}

	next := func(ctx *BotContext) error {
		return handler(ctx)
	}

	// Apply middlewares in reverse order
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		middleware := h.middlewares[i]
		prevNext := next
		next = func(ctx *BotContext) error {
			return middleware(ctx, prevNext)
		}
	}

	if err := next(ctx); err != nil {
		utils.Logger.WithError(err).Error("Handler error")
		if ctx.User != nil && ctx.Update.Message != nil {
			msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
			_, _ = h.bot.Send(msg)
		}
	}
}

// getHandler returns appropriate handler for update
func (h *BotHandler) getHandler(update tgbotapi.Update) HandlerFunc {
	if update.Message != nil {
		// Handle commands
		if update.Message.IsCommand() {
			command := update.Message.Command()
			if handler, ok := h.handlers["/"+command]; ok {
				return handler
			}
		}
		
		// Handle contact (phone number)
		if update.Message.Contact != nil {
			return h.handleContact
		}
		
		// Handle text messages based on state
		return h.handleTextMessage
	}
	if update.CallbackQuery != nil {
		return h.handleCallbackQuery
	}
	if update.PreCheckoutQuery != nil {
		return h.handlePreCheckout
	}
	if update.Message != nil && update.Message.SuccessfulPayment != nil {
		return h.handleSuccessfulPayment
	}
	return nil
}

// handleStart handles /start command
func (h *BotHandler) handleStart(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	
	// If user is admin, show admin panel instead of regular start
	if ctx.Update.Message != nil && h.config.IsAdmin(ctx.Update.Message.From.ID) {
		// Clear any existing state
		ctx.StateManager.ClearState(int64(ctx.Update.Message.From.ID))
		return h.adminPanelHandler.ShowAdminPanel(handlerCtx)
	}
	
	return h.startHandler.HandleStart(handlerCtx)
}

// handleContact handles phone number contact
func (h *BotHandler) handleContact(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.startHandler.HandleContact(handlerCtx)
}

// handleHelp handles /help command
func (h *BotHandler) handleHelp(ctx *BotContext) error {
	text := ctx.Translator.Get("help_title") + "\n\n" +
		ctx.Translator.Get("help_start") + "\n" +
		ctx.Translator.Get("help_profile") + "\n" +
		ctx.Translator.Get("help_balance") + "\n" +
		ctx.Translator.Get("help_services") + "\n" +
		ctx.Translator.Get("help_support") + "\n" +
		ctx.Translator.Get("help_language") + "\n" +
		ctx.Translator.Get("help_help")
	
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err := ctx.Bot.Send(msg)
	return err
}

// handleProfile handles /profile command
func (h *BotHandler) handleProfile(ctx *BotContext) error {
	user := ctx.User
	if user == nil {
		return fmt.Errorf("user not found")
	}

	var text string
	if user.Username != nil {
		text = ctx.Translator.Getf("profile_username", *user.Username)
	} else {
		text = ctx.Translator.Getf("profile_id", user.ID)
	}
	
	if user.FirstName != nil {
		lastName := ""
		if user.LastName != nil {
			lastName = *user.LastName
		}
		text += "\n" + ctx.Translator.Getf("profile_name", *user.FirstName, lastName)
	}
	
	text += "\n" + ctx.Translator.Getf("profile_language", user.Language)

	// Get balance from billing
	billingUser, err := h.billingService.GetBillingUser(context.Background(), int64(user.ID))
	if err == nil && billingUser != nil {
		text += "\n" + ctx.Translator.Getf("profile_balance", billingUser.Balance)
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err = ctx.Bot.Send(msg)
	return err
}

// handleBalance handles /balance command
func (h *BotHandler) handleBalance(ctx *BotContext) error {
	user := ctx.User
	if user == nil {
		return fmt.Errorf("user not found")
	}

	billingUser, err := h.billingService.GetBillingUser(context.Background(), int64(user.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	text := ctx.Translator.Getf("balance_amount", billingUser.Balance)
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	
	// Add top-up button
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("balance_topup"), "topup_balance"),
		),
	)
	msg.ReplyMarkup = keyboard
	
	_, err = ctx.Bot.Send(msg)
	return err
}

// handleServices handles /services command
func (h *BotHandler) handleServices(ctx *BotContext) error {
	user := ctx.User
	if user == nil {
		return fmt.Errorf("user not found")
	}

	services, err := h.billingService.GetServices(context.Background(), int64(user.ID))
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("error"))
		_, _ = ctx.Bot.Send(msg)
		return err
	}

	text := ctx.Translator.Get("services_title") + "\n\n"
	if len(services) == 0 {
		text += ctx.Translator.Get("services_none")
	} else {
		text += ctx.Translator.Get("services_list") + "\n"
		for _, service := range services {
			status := ctx.Translator.Get("service_" + service.Status)
			text += fmt.Sprintf("• %s - %s\n", service.Name, status)
		}
	}

	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err = ctx.Bot.Send(msg)
	return err
}

// handleSupport handles /support command
func (h *BotHandler) handleSupport(ctx *BotContext) error {
	text := ctx.Translator.Get("support_message")
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, text)
	_, err := ctx.Bot.Send(msg)
	return err
}

// handleLanguage handles /language command
func (h *BotHandler) handleLanguage(ctx *BotContext) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("language_ua"), "lang_ua"),
			tgbotapi.NewInlineKeyboardButtonData(ctx.Translator.Get("language_en"), "lang_en"),
		),
	)
	
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("language_title"))
	msg.ReplyMarkup = keyboard
	_, err := ctx.Bot.Send(msg)
	return err
}

// handleTextMessage handles regular text messages
func (h *BotHandler) handleTextMessage(ctx *BotContext) error {
	if ctx.Update.Message == nil || ctx.Update.Message.Text == "" {
		return nil
	}

	text := ctx.Update.Message.Text
	userState, _, exists := h.stateManager.GetState(int64(ctx.Update.Message.From.ID))

	// Create handler context
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}

	// Check if admin sending message
	if h.config.IsAdmin(ctx.Update.Message.From.ID) {
		// Handle "Назад" button for admins - show admin panel
		if text == ctx.Translator.Get("menu_back") {
			ctx.StateManager.ClearState(int64(ctx.Update.Message.From.ID))
			return h.adminPanelHandler.ShowAdminPanel(handlerCtx)
		}
		
		// Check if admin is connected to an active support chat
		// This check should happen before admin panel handling
		if h.supportChatHandler.IsAdminInChat(int64(ctx.Update.Message.From.ID)) {
			// Admin is in active support chat - handle as support message
			supportCtx := &support.SupportContext{
				Update:       ctx.Update,
				User:         ctx.User,
				Translator:   ctx.Translator,
				Config:       ctx.Config,
				StateManager: ctx.StateManager,
			}
			return h.supportChatHandler.HandleAdminMessage(supportCtx)
		}
		
		// Handle admin panel text messages if not in support chat
		return h.adminPanelHandler.HandleTextMessage(handlerCtx)
	}

	// Handle states first (before menu buttons) to ensure support chat messages are processed
	if exists {
		switch userState {
		case state.StateWaitingInvoicePayload:
			return h.payBillHandler.HandleAmountInput(handlerCtx)
		case state.StateWaitingInvoiceContract, state.StateInvalidPayload:
			return h.payBillHandler.HandleContractInput(handlerCtx)
		case state.StateWaitingForSupport, state.StateChatting:
			// Handle support chat messages - check if user has active chat first
			supportCtx := &support.SupportContext{
				Update:       ctx.Update,
				User:         ctx.User,
				Translator:   ctx.Translator,
				Config:       ctx.Config,
				StateManager: ctx.StateManager,
			}
			// Always try to handle as support message if in support state
			return h.supportChatHandler.HandleSupportMessage(supportCtx)
		}
	}

	// Handle menu buttons (only if not in a state)
	if text == ctx.Translator.Get("menu_topup") {
		return h.payBillHandler.HandleTopUp(handlerCtx)
	}
	if text == ctx.Translator.Get("menu_support") {
		supportCtx := &support.SupportContext{
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return h.supportChatHandler.HandleStartSupport(supportCtx)
	}
	if text == ctx.Translator.Get("menu_time_pay") {
		return h.timePayHandler.HandleTimePay(handlerCtx)
	}
	if text == ctx.Translator.Get("support_end_chat_button") {
		// User wants to end support chat
		supportCtx := &support.SupportContext{
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return h.supportChatHandler.HandleEndChat(supportCtx)
	}
	if text == ctx.Translator.Get("menu_back") {
		return h.startHandler.ShowMainMenu(handlerCtx, ctx.User != nil && ctx.User.Contract != nil)
	}

	// Default: create support ticket
	userID := (*int)(nil)
	if ctx.User != nil {
		userID = &ctx.User.ID
	}
	
	if err := h.supportService.CreateTicket(
		context.Background(),
		userID,
		ctx.Update.Message.From.ID,
		text,
	); err != nil {
		utils.Logger.WithError(err).Error("Failed to create support ticket")
	} else {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("support_sent"))
		_, _ = ctx.Bot.Send(msg)
	}
	
	return nil
}

// handleCallbackQuery handles callback queries
func (h *BotHandler) handleCallbackQuery(ctx *BotContext) error {
	callback := ctx.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	data := callback.Data

	// Handle language selection callbacks
	if strings.HasPrefix(data, "start_lang_") {
		handlerCtx := &handlers.HandlerContext{
			Bot:          ctx.Bot,
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return h.startHandler.HandleLanguageCallback(handlerCtx)
	}

	switch data {
	case "lang_ua":
		if err := h.userService.UpdateLanguage(context.Background(), int64(callback.From.ID), "ua"); err == nil {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, ctx.Translator.Get("language_changed"))
			_, _ = ctx.Bot.Send(msg)
		}
	case "lang_en":
		if err := h.userService.UpdateLanguage(context.Background(), int64(callback.From.ID), "en"); err == nil {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, ctx.Translator.Get("language_changed"))
			_, _ = ctx.Bot.Send(msg)
		}
	case "lang_ru":
		if err := h.userService.UpdateLanguage(context.Background(), int64(callback.From.ID), "ru"); err == nil {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, ctx.Translator.Get("language_changed"))
			_, _ = ctx.Bot.Send(msg)
		}
	}

	// Handle support callbacks (check before admin callbacks since support can be for admins)
	if strings.HasPrefix(data, "support_connect_") {
		if !h.config.IsAdmin(callback.From.ID) {
			callbackConfig := tgbotapi.NewCallback(callback.ID, "У вас немає прав доступу")
			_, _ = ctx.Bot.Request(callbackConfig)
			return nil
		}
		
		userIDStr := strings.TrimPrefix(data, "support_connect_")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Помилка: невірний ID користувача")
			_, _ = ctx.Bot.Request(callbackConfig)
			return nil
		}
		
		supportCtx := &support.SupportContext{
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return h.supportChatHandler.HandleAdminConnect(supportCtx, userID)
	}
	
	if data == "support_end_chat" {
		if !h.config.IsAdmin(callback.From.ID) {
			callbackConfig := tgbotapi.NewCallback(callback.ID, "У вас немає прав доступу")
			_, _ = ctx.Bot.Request(callbackConfig)
			return nil
		}
		
		supportCtx := &support.SupportContext{
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return h.supportChatHandler.HandleEndChat(supportCtx)
	}

	// Handle admin callbacks
	if h.config.IsAdmin(callback.From.ID) {
		handlerCtx := &handlers.HandlerContext{
			Bot:          ctx.Bot,
			Update:       ctx.Update,
			User:         ctx.User,
			Translator:   ctx.Translator,
			Config:       ctx.Config,
			StateManager: ctx.StateManager,
		}
		return admin.HandleCallbackQuery(handlerCtx, h.adminPanelHandler)
	}

	// Answer callback query
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	_, _ = ctx.Bot.Request(callbackConfig)
	return nil
}

// handlePreCheckout handles pre-checkout query for payments
func (h *BotHandler) handlePreCheckout(ctx *BotContext) error {
	query := ctx.Update.PreCheckoutQuery
	if query == nil {
		return nil
	}

	// Always approve for now (should validate in production)
	answer := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: query.ID,
		OK:                  true,
	}
	_, err := ctx.Bot.Request(answer)
	return err
}

// handleSuccessfulPayment handles successful payment
func (h *BotHandler) handleSuccessfulPayment(ctx *BotContext) error {
	if ctx.Update.Message == nil || ctx.Update.Message.SuccessfulPayment == nil {
		return nil
	}
	payment := ctx.Update.Message.SuccessfulPayment

	amount := float64(payment.TotalAmount) / 100.0
	payload := payment.InvoicePayload

	// Parse amount from payload
	amountKopecks, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to parse payment amount")
		return err
	}

	// Extract contract from payment description or state
	// This is simplified - in production should get from payment metadata
	contract := ctx.Update.Message.SuccessfulPayment.InvoicePayload

	// Update balance in billing
	// Note: This requires getting user ID from contract - simplified here
	utils.Logger.Infof("Payment successful: %.2f UAH, contract: %s", amount, contract)

	// Send confirmation to user
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, 
		ctx.Translator.Getf("balance_topup_success", float64(amountKopecks)/100.0))
	_, _ = ctx.Bot.Send(msg)

	return nil
}

// handleConnect handles /connect command for admins
func (h *BotHandler) handleConnect(ctx *BotContext) error {
	if !h.config.IsAdmin(ctx.Update.Message.From.ID) {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_unauthorized"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	// Parse user ID from command arguments
	args := strings.Fields(ctx.Update.Message.Text)
	if len(args) < 2 {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Використання: /connect <user_id>")
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Невірний формат user_id")
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	supportCtx := &support.SupportContext{
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.supportChatHandler.HandleAdminConnect(supportCtx, userID)
}

// handleEndChat handles /end_chat command for admins
func (h *BotHandler) handleEndChat(ctx *BotContext) error {
	if !h.config.IsAdmin(ctx.Update.Message.From.ID) {
		msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, ctx.Translator.Get("admin_unauthorized"))
		_, _ = ctx.Bot.Send(msg)
		return nil
	}

	supportCtx := &support.SupportContext{
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.supportChatHandler.HandleEndChat(supportCtx)
}

// handleStats handles /stats command for admins
func (h *BotHandler) handleStats(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.adminPanelHandler.HandleStats(handlerCtx)
}

// Admin handlers
func (h *BotHandler) handleAdmin(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.adminPanelHandler.HandleAdmin(handlerCtx)
}

func (h *BotHandler) handleBroadcast(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.adminPanelHandler.BroadcastHandler.HandleSendMessage(handlerCtx)
}

func (h *BotHandler) handleOutage(ctx *BotContext) error {
	// Outage management not implemented yet
	return nil
}

func (h *BotHandler) handleUsers(ctx *BotContext) error {
	handlerCtx := &handlers.HandlerContext{
		Bot:          ctx.Bot,
		Update:       ctx.Update,
		User:         ctx.User,
		Translator:   ctx.Translator,
		Config:       ctx.Config,
		StateManager: ctx.StateManager,
	}
	return h.adminPanelHandler.UsersHandler.HandleAccountMenu(handlerCtx)
}

