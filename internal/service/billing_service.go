package service

import (
	"context"
	"fmt"
	"time"

	"provbot/internal/models"
	"provbot/internal/repository"
)

type BillingService struct {
	billingRepo *repository.BillingRepository
}

func NewBillingService(billingRepo *repository.BillingRepository) *BillingService {
	return &BillingService{
		billingRepo: billingRepo,
	}
}

// GetBalance retrieves user balance
func (s *BillingService) GetBalance(ctx context.Context, contract string) (float64, error) {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}
	return user.Balance, nil
}

// CheckBalance checks if user has enough balance for their plan
// Returns inequality (shortage) if balance is low, or nil if enough
func (s *BillingService) CheckBalance(ctx context.Context, contract string) (*float64, error) {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.PlanID == nil {
		return nil, fmt.Errorf("user has no plan")
	}

	plan, err := s.billingRepo.GetPlanByID(ctx, *user.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	shortage := plan.Price - user.Balance
	if shortage > 0 {
		return &shortage, nil
	}
	return nil, nil
}

// PayBalance processes a payment
func (s *BillingService) PayBalance(ctx context.Context, contract string, amount float64, admin string) error {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Update balance
	newBalance := user.Balance + amount
	if err := s.billingRepo.UpdateBalance(ctx, contract, newBalance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Add payment record
	pay := &models.BillingPay{
		UserID:  user.ID,
		Amount:  amount,
		Time:    float64(time.Now().Unix()),
		Admin:   admin,
		Reason:  fmt.Sprintf("%f", float64(time.Now().Unix())), // Python uses timestamp as reason?
		Comment: "Popolnenie via BOT",
		Bonus:   "",
		Flag:    "",
	}
	if err := s.billingRepo.AddPayment(ctx, pay); err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}

	// Check if balance is enough to unlock
	if user.PlanID != nil {
		plan, err := s.billingRepo.GetPlanByID(ctx, *user.PlanID)
		if err == nil && plan != nil {
			if newBalance >= plan.Price {
				if err := s.billingRepo.SetUserStatus(ctx, contract, "on"); err != nil {
					return fmt.Errorf("failed to unlock user: %w", err)
				}
			}
		}
	}

	return nil
}

// TemporaryPay enables temporary payment (credit)
func (s *BillingService) TemporaryPay(ctx context.Context, contract string) (bool, error) {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	if user.TimePay != 0 {
		return false, nil // Already used
	}

	if user.PlanID == nil {
		return false, fmt.Errorf("user has no plan")
	}

	plan, err := s.billingRepo.GetPlanByID(ctx, *user.PlanID)
	if err != nil {
		return false, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan == nil {
		return false, fmt.Errorf("plan not found")
	}

	// Calculate new balance logic from Python:
	// if old_balance > 0: balance = old_balance + price
	// else: balance = -old_balance + price (Wait, Python says: balance = -old_balance + price? No, let's re-read)
	// Python:
	// if old_balance > 0: balance = old_balance + price
	// else: balance = -old_balance + price  <-- This looks weird in Python code.
	// Let's assume it means we give enough credit to cover the debt + price?
	// Or maybe it's a typo in Python code? "balance = -old_balance + price" would mean if I have -100, I get 100 + price?
	// Let's stick to a safer logic: Add price to balance.
	// Actually, let's look at Python code again.
	// if old_balance > 0: balance = old_balance + price
	// else: balance = -old_balance + price
	// If balance is -50, and price is 150. New balance = -(-50) + 150 = 50 + 150 = 200.
	// This seems to set balance to POSITIVE price + abs(debt).
	// This effectively gives them `price + 2*debt` credit? No.
	// If I have -50. I need 150.
	// If I set balance to 200. I effectively gave 250.
	// Maybe the intention was `balance = price`?
	// Let's implement as "Add Price to Balance" for now, but maybe add a comment.
	// Actually, let's replicate Python logic exactly to be safe, or maybe "Set balance to Price" if negative?
	// Let's look at Python again:
	// `balance = -old_balance + price`
	// If old = -10. Price = 100. New = 10 + 100 = 110.
	// If I had -10, and now 110. I gained 120.
	// If I had 10. New = 10 + 100 = 110. I gained 100.
	// It seems it tries to compensate for negative balance?

	var newBalance float64
	if user.Balance > 0 {
		newBalance = user.Balance + plan.Price
	} else {
		newBalance = -user.Balance + plan.Price
	}

	// Enable temp pay
	if err := s.billingRepo.EnableTemporaryPayment(ctx, contract, newBalance); err != nil {
		return false, fmt.Errorf("failed to enable temp pay: %w", err)
	}

	// Add payment record
	pay := &models.BillingPay{
		UserID:  user.ID,
		Amount:  plan.Price,
		Time:    float64(time.Now().Unix() + 86400), // Next time? Python: next_t
		Admin:   "BOT",
		Reason:  fmt.Sprintf("Platej sozdan %f", float64(time.Now().Unix())),
		Comment: "Razblokirovan na 24 chasa",
		Bonus:   "y",
		Flag:    "t",
	}
	if err := s.billingRepo.AddPayment(ctx, pay); err != nil {
		// Log error but don't fail as DB update succeeded
		fmt.Printf("failed to add payment record for temp pay: %v\n", err)
	}

	return true, nil
}

// SearchUser searches for a user by various criteria
func (s *BillingService) SearchUser(ctx context.Context, query string) (*models.BillingUser, error) {
	// Try by contract first
	user, err := s.billingRepo.GetUserByContract(ctx, query)
	if err == nil && user != nil {
		return user, nil
	}

	// Try by phone
	user, err = s.billingRepo.SearchByPhone(ctx, query)
	if err == nil && user != nil {
		return user, nil
	}

	return nil, nil
}

// GetTariff retrieves tariff info (name, price) for a contract
func (s *BillingService) GetTariff(ctx context.Context, contract string) (string, float64, error) {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", 0, fmt.Errorf("user not found")
	}

	if user.PlanID == nil {
		return "No Plan", 0, nil
	}

	plan, err := s.billingRepo.GetPlanByID(ctx, *user.PlanID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan == nil {
		return "Unknown Plan", 0, nil
	}

	return plan.Name, plan.Price, nil
}

// CheckContract checks if a contract exists
func (s *BillingService) CheckContract(ctx context.Context, contract string) (bool, error) {
	user, err := s.billingRepo.GetUserByContract(ctx, contract)
	if err != nil {
		return false, fmt.Errorf("failed to check contract: %w", err)
	}
	return user != nil, nil
}

// GetUserByID retrieves billing user by ID
func (s *BillingService) GetUserByID(ctx context.Context, id int64) (*models.BillingUser, error) {
	user, err := s.billingRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

// UpdateBalance updates user balance directly (for admin usage)
func (s *BillingService) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	// We need contract to update balance in repo, but we have userID
	user, err := s.billingRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Update balance
	// Note: Repo UpdateBalance takes contract.
	// If we want to set specific balance (not add), we need to know if 'amount' is new balance or delta.
	// Admin handler says "UpdateBalance", usually means "Set Balance" or "Add to Balance"?
	// Python `balance_change` sets balance: `UPDATE users set balance = {cash}`.
	// My `BillingRepository.UpdateBalance` does `UPDATE users SET balance = ?`.
	// So it sets the balance.

	return s.billingRepo.UpdateBalance(ctx, user.Contract, amount)
}

// SearchByName searches users by name
func (s *BillingService) SearchByName(ctx context.Context, name string) ([]models.BillingUser, error) {
	return s.billingRepo.SearchByName(ctx, name)
}

// SearchByAddress searches users by address
func (s *BillingService) SearchByAddress(ctx context.Context, address string) ([]models.BillingUser, error) {
	return s.billingRepo.SearchByAddress(ctx, address)
}

// GetBillingUser retrieves billing user by ID (alias for GetUserByID)
func (s *BillingService) GetBillingUser(ctx context.Context, id int64) (*models.BillingUser, error) {
	return s.GetUserByID(ctx, id)
}

// GetServices retrieves services for a user
func (s *BillingService) GetServices(ctx context.Context, userID int64) ([]models.BillingService, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var services []models.BillingService
	if user.Services != "" {
		// Assuming user.Services contains service names or is a single service description
		// We create a single service entry for now
		services = append(services, models.BillingService{
			UserID: userID,
			Name:   user.Services,
			Status: "active", // Placeholder status
		})
	}

	return services, nil
}
