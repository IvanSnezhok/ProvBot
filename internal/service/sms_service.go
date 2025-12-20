package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"provbot/internal/utils"
)

// SMSService handles SMS sending functionality
type SMSService struct {
	apiURL    string
	apiKey    string
	sender    string
	enabled   bool
	client    *http.Client
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	APIURL  string
	APIKey  string
	Sender  string
	Enabled bool
}

// SMSMessage represents an SMS message to send
type SMSMessage struct {
	Phone   string
	Message string
}

// SMSResponse represents response from SMS API
type SMSResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// NewSMSService creates a new SMS service
func NewSMSService(config *SMSConfig) *SMSService {
	if config == nil {
		config = &SMSConfig{
			Enabled: false,
		}
	}

	return &SMSService{
		apiURL:  config.APIURL,
		apiKey:  config.APIKey,
		sender:  config.Sender,
		enabled: config.Enabled,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsEnabled returns whether SMS service is enabled
func (s *SMSService) IsEnabled() bool {
	return s.enabled && s.apiKey != "" && s.apiURL != ""
}

// SendSMS sends an SMS message to the specified phone number
func (s *SMSService) SendSMS(phone, message string) error {
	if !s.IsEnabled() {
		utils.Logger.Warn("SMS service is disabled, skipping SMS send")
		return fmt.Errorf("SMS service is disabled")
	}

	// Format phone number (ensure it starts with country code)
	formattedPhone := s.formatPhone(phone)

	// Prepare request body (structure depends on SMS provider)
	requestBody := map[string]interface{}{
		"phone":   formattedPhone,
		"message": message,
		"sender":  s.sender,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("SMS API returned error: %s, body: %s", resp.Status, string(body))
	}

	// Parse response
	var smsResp SMSResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		utils.Logger.WithError(err).Warn("Failed to parse SMS response, but request was successful")
	}

	utils.Logger.Infof("SMS sent successfully to %s, ID: %s", formattedPhone, smsResp.ID)
	return nil
}

// formatPhone formats phone number for SMS API
func (s *SMSService) formatPhone(phone string) string {
	// Remove all non-digit characters
	var digits []rune
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	result := string(digits)

	// Ensure phone starts with country code (380 for Ukraine)
	if len(result) == 10 && result[0] == '0' {
		result = "38" + result
	} else if len(result) == 9 {
		result = "380" + result
	}

	return result
}

// SendBulkSMS sends SMS to multiple recipients
func (s *SMSService) SendBulkSMS(messages []SMSMessage) (int, int) {
	if !s.IsEnabled() {
		utils.Logger.Warn("SMS service is disabled, skipping bulk SMS send")
		return 0, len(messages)
	}

	successCount := 0
	failCount := 0

	for _, msg := range messages {
		if err := s.SendSMS(msg.Phone, msg.Message); err != nil {
			utils.Logger.WithError(err).Errorf("Failed to send SMS to %s", msg.Phone)
			failCount++
		} else {
			successCount++
		}
	}

	return successCount, failCount
}
