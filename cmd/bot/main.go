package main

import (
	"os"
	"os/signal"
	"syscall"

	"provbot/internal/bot"
	"provbot/internal/database"
	"provbot/internal/repository"
	"provbot/internal/scheduler"
	"provbot/internal/service"
	"provbot/internal/state"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load configuration first to get log level
	config, err := utils.LoadConfig()
	if err != nil {
		// Use basic logger if config fails
		utils.InitLogger("info")
		utils.Logger.WithError(err).Fatal("Failed to load configuration")
		return
	}

	// Initialize file logger with rotation
	logLevel := config.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}
	if err := utils.InitFileLogger(logLevel, nil); err != nil {
		// Fallback to basic logger if file logger fails
		utils.InitLogger(logLevel)
		utils.Logger.WithError(err).Warn("Failed to initialize file logger, using stdout only")
	} else {
		utils.Logger.Info("File logger initialized successfully")
	}

	utils.Logger.Info("Starting ProvBot...")

	// Initialize databases
	if err := database.InitPostgres(config); err != nil {
		utils.Logger.WithError(err).Fatal("Failed to initialize PostgreSQL")
	}
	defer database.ClosePostgres()

	if err := database.InitMySQL(config); err != nil {
		utils.Logger.WithError(err).Warn("Failed to initialize MySQL (billing), continuing without billing features")
	} else {
		defer database.CloseMySQL()
	}

	// Initialize Telegram bot
	telegramBot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		utils.Logger.WithError(err).Fatal("Failed to initialize Telegram bot")
	}

	utils.Logger.Infof("Authorized on account %s", telegramBot.Self.UserName)

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	logRepo := repository.NewLogRepository()
	outageRepo := repository.NewOutageRepository()
	billingRepo := repository.NewBillingRepository()

	// Initialize services
	userService := service.NewUserService(userRepo)
	billingService := service.NewBillingService(billingRepo)
	supportService := service.NewSupportService(logRepo)
	adminService := service.NewAdminService(userRepo, outageRepo, billingRepo)
	notificationService := service.NewNotificationService(telegramBot, logRepo, userRepo)

	// Initialize scheduler for balance notifications
	balanceScheduler := scheduler.NewBalanceNotificationScheduler(
		telegramBot,
		userRepo,
		billingRepo,
		config,
	)
	balanceScheduler.Start()
	defer balanceScheduler.Stop()

	// Initialize state manager
	stateManager := state.GetStateManagerInstance()

	// Initialize bot handler
	botHandler := bot.NewBotHandler(
		telegramBot,
		config,
		userService,
		billingService,
		supportService,
		adminService,
		notificationService,
		stateManager,
		userRepo,
		billingRepo,
		logRepo,
	)

	// Add middlewares (order matters - last added executes first)
	botHandler.Use(bot.UserRegistrationMiddleware(userService))
	botHandler.Use(bot.MessageLoggingMiddleware(logRepo))
	botHandler.Use(bot.HandlerLoggingMiddleware()) // Log handler execution

	// Setup update channel
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := telegramBot.GetUpdatesChan(updateConfig)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	utils.Logger.Info("Bot is running. Press Ctrl+C to stop.")

	// Process updates
	for {
		select {
		case update := <-updates:
			go botHandler.HandleUpdate(update)
		case <-sigChan:
			utils.Logger.Info("Shutting down...")
			telegramBot.StopReceivingUpdates()
			return
		}
	}
}
