package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"provbot/internal/i18n"
	"provbot/internal/models"
	"provbot/internal/repository"
	"provbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BalanceNotificationScheduler handles scheduled balance notifications
type BalanceNotificationScheduler struct {
	bot         *tgbotapi.BotAPI
	userRepo    *repository.UserRepository
	billingRepo *repository.BillingRepository
	config      *utils.Config
	workerCount int
	running     bool
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// NewBalanceNotificationScheduler creates a new scheduler instance
func NewBalanceNotificationScheduler(
	bot *tgbotapi.BotAPI,
	userRepo *repository.UserRepository,
	billingRepo *repository.BillingRepository,
	config *utils.Config,
) *BalanceNotificationScheduler {
	// Use worker count from config or default to 10 for parallel processing
	workerCount := 10
	if config != nil && config.SchedulerWorkers > 0 {
		workerCount = config.SchedulerWorkers
	}

	return &BalanceNotificationScheduler{
		bot:         bot,
		userRepo:    userRepo,
		billingRepo: billingRepo,
		config:      config,
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *BalanceNotificationScheduler) Start() {
	if s.running {
		utils.Logger.Warn("Scheduler is already running")
		return
	}

	s.running = true
	s.wg.Add(1)
	go s.run()
	utils.Logger.Info("Balance notification scheduler started")
}

// Stop stops the scheduler gracefully
func (s *BalanceNotificationScheduler) Stop() {
	if !s.running {
		return
	}

	close(s.stopChan)
	s.wg.Wait()
	s.running = false
	utils.Logger.Info("Balance notification scheduler stopped")
}

// run is the main scheduler loop
func (s *BalanceNotificationScheduler) run() {
	defer s.wg.Done()

	// Calculate next run time (9th of month at 23:45)
	now := time.Now()
	nextRun := s.calculateNextRun(now)

	utils.Logger.Infof("Next balance notification scheduled for: %s", nextRun.Format(time.RFC3339))

	for {
		now = time.Now()
		waitDuration := nextRun.Sub(now)

		if waitDuration <= 0 {
			// Time to run
			s.processNotifications()
			nextRun = s.calculateNextRun(now)
			utils.Logger.Infof("Next balance notification scheduled for: %s", nextRun.Format(time.RFC3339))
			continue
		}

		// Wait until next run or stop signal
		select {
		case <-time.After(waitDuration):
			// Time to run
			s.processNotifications()
			nextRun = s.calculateNextRun(time.Now())
			utils.Logger.Infof("Next balance notification scheduled for: %s", nextRun.Format(time.RFC3339))
		case <-s.stopChan:
			return
		}
	}
}

// calculateNextRun calculates the next run time (9th of month at 23:45)
func (s *BalanceNotificationScheduler) calculateNextRun(now time.Time) time.Time {
	// Target: 9th of current month at 23:45
	targetDay := 9
	targetHour := 23
	targetMinute := 45

	// If we're past the 9th at 23:45, schedule for next month
	nextRun := time.Date(now.Year(), now.Month(), targetDay, targetHour, targetMinute, 0, 0, now.Location())
	if now.After(nextRun) || now.Equal(nextRun) {
		// Move to next month
		nextRun = nextRun.AddDate(0, 1, 0)
	}

	return nextRun
}

// processNotifications processes balance notifications for all users
func (s *BalanceNotificationScheduler) processNotifications() {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	utils.Logger.WithFields(map[string]interface{}{
		"component": "scheduler",
		"action":    "process_notifications",
		"workers":   s.workerCount,
	}).Info("Starting balance notification process")

	// Get all users with contracts
	users, err := s.userRepo.GetUsersWithContracts(ctx)
	if err != nil {
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "get_users",
			"error":     err.Error(),
		}).Error("Failed to get users with contracts")
		return
	}

	if len(users) == 0 {
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "process_notifications",
		}).Info("No users with contracts found")
		return
	}

	utils.Logger.WithFields(map[string]interface{}{
		"component":   "scheduler",
		"action":      "process_notifications",
		"users_count": len(users),
		"workers":     s.workerCount,
	}).Info("Processing users for balance notifications")

	// Create worker pool for parallel processing
	userChan := make(chan models.User, len(users))
	var wg sync.WaitGroup

	// Statistics
	var stats struct {
		processed int
		notified  int
		skipped   int
		errors    int
		sync.Mutex
	}

	// Start workers
	for i := 0; i < s.workerCount; i++ {
		wg.Add(1)
		workerID := i
		go s.worker(ctx, userChan, &wg, workerID, &stats)
	}

	// Send users to workers
	for _, user := range users {
		userChan <- user
	}
	close(userChan)

	// Wait for all workers to finish
	wg.Wait()

	duration := time.Since(startTime)
	utils.Logger.WithFields(map[string]interface{}{
		"component":   "scheduler",
		"action":      "process_notifications",
		"duration_ms": duration.Milliseconds(),
		"processed":   stats.processed,
		"notified":    stats.notified,
		"skipped":     stats.skipped,
		"errors":      stats.errors,
		"users_total": len(users),
	}).Info("Balance notification process completed")
}

// worker processes users from the channel
func (s *BalanceNotificationScheduler) worker(ctx context.Context, userChan <-chan models.User, wg *sync.WaitGroup, workerID int, stats *struct {
	processed int
	notified  int
	skipped   int
	errors    int
	sync.Mutex
}) {
	defer wg.Done()

	workerStartTime := time.Now()
	processed := 0

	for user := range userChan {
		processed++
		s.processUser(ctx, user, workerID, stats)
	}

	workerDuration := time.Since(workerStartTime)
	utils.Logger.WithFields(map[string]interface{}{
		"component":   "scheduler",
		"action":      "worker_completed",
		"worker_id":   workerID,
		"processed":   processed,
		"duration_ms": workerDuration.Milliseconds(),
	}).Debug("Worker completed processing")
}

// processUser processes a single user's balance notification
func (s *BalanceNotificationScheduler) processUser(ctx context.Context, user models.User, workerID int, stats *struct {
	processed int
	notified  int
	skipped   int
	errors    int
	sync.Mutex
}) {
	processedUserStartTime := time.Now()

	if user.Contract == nil || *user.Contract == "" {
		stats.Lock()
		stats.errors++
		stats.Unlock()
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "process_user",
			"user_id":   user.TelegramID,
			"error":     "no_contract",
		}).Warn("User has no contract, skipping")
		return
	}

	// Get balance from billing system
	billingUser, err := s.billingRepo.GetUserByContract(ctx, *user.Contract)
	if err != nil {
		stats.Lock()
		stats.errors++
		stats.Unlock()
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "get_balance",
			"user_id":   user.TelegramID,
			"contract":  *user.Contract,
			"error":     err.Error(),
		}).Error("Failed to get balance for contract")
		return
	}

	if billingUser == nil {
		stats.Lock()
		stats.errors++
		stats.Unlock()
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "get_balance",
			"user_id":   user.TelegramID,
			"contract":  *user.Contract,
			"error":     "user not found in billing",
		}).Error("User not found in billing")
		return
	}

	balance := billingUser.Balance

	stats.Lock()
	stats.processed++
	stats.Unlock()

	// Check if balance is low (less than threshold, e.g., 0 or negative)
	// This logic can be adjusted based on business requirements
	if balance > 0 {
		// Balance is positive, skip notification
		stats.Lock()
		stats.skipped++
		stats.Unlock()
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "process_user",
			"user_id":   user.TelegramID,
			"contract":  *user.Contract,
			"balance":   balance,
			"reason":    "balance_positive",
		}).Debug("User balance is positive, skipping notification")
		return
	}

	// Get translator for user's language
	translator := i18n.GetGlobalTranslator(user.Language)
	message := translator.Getf("balance_notification_message", balance)
	if message == "balance_notification_message" {
		// Fallback if translation not found
		message = fmt.Sprintf("⚠️ Увага! Ваш баланс низький (%.2f грн). Можливе блокування послуг 12 числа. Рекомендуємо поповнити рахунок.", balance)
	}

	// Send notification
	msg := tgbotapi.NewMessage(user.TelegramID, message)
	msg.ParseMode = tgbotapi.ModeHTML

	// Use context with timeout for API call
	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Send message (non-blocking with timeout)
	done := make(chan error, 1)
	go func() {
		_, err := s.bot.Send(msg)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			stats.Lock()
			stats.errors++
			stats.Unlock()
			utils.Logger.WithFields(map[string]interface{}{
				"component": "scheduler",
				"action":    "send_notification",
				"user_id":   user.TelegramID,
				"contract":  *user.Contract,
				"balance":   balance,
				"error":     err.Error(),
			}).Error("Failed to send balance notification")
		} else {
			stats.Lock()
			stats.notified++
			stats.Unlock()
			duration := time.Since(processedUserStartTime)
			utils.Logger.WithFields(map[string]interface{}{
				"component":   "scheduler",
				"action":      "send_notification",
				"user_id":     user.TelegramID,
				"contract":    *user.Contract,
				"balance":     balance,
				"duration_ms": duration.Milliseconds(),
			}).Info("Balance notification sent successfully")
		}
	case <-sendCtx.Done():
		stats.Lock()
		stats.errors++
		stats.Unlock()
		utils.Logger.WithFields(map[string]interface{}{
			"component": "scheduler",
			"action":    "send_notification",
			"user_id":   user.TelegramID,
			"contract":  *user.Contract,
			"balance":   balance,
			"error":     "timeout",
		}).Warn("Timeout sending balance notification")
	}

	// Small delay to avoid rate limiting
	time.Sleep(50 * time.Millisecond)
}
