package i18n

// GetENTranslations returns English translations
func GetENTranslations() map[string]string {
	return map[string]string{
		// Common
		"welcome":           "Welcome to the provider bot!",
		"error":             "An error occurred. Please try again later.",
		"unknown_command":   "Unknown command. Use /help for a list of commands.",
		"cancel":             "Cancel",
		"back":               "Back",
		"confirm":            "Confirm",
		"select":             "Select",
		
		// Start command
		"start_message":     "Hello! I'm the provider bot. I'll help you manage services, check balance, and contact support.\n\nPlease choose your preferred language!",
		"start_registered":  "You are already registered! Use /help for a list of commands.",
		"phone_request":     "You selected %s\nNow please send your phone number to find your account in our billing system",
		"send_phone":        "📱 Send phone number",
		"welcome_registered": "Welcome! You have been successfully registered. Use the menu to navigate.",
		"menu_topup":        "Top up account",
		"menu_support":      "Support chat",
		"menu_time_pay":     "Temporary payment",
		"menu_back":         "Back",
		
		// Payment
		"topup_notice":      "Please note that here you can only top up your personal account!",
		"topup_promotion":   "Special offer - top up your account for 6 months with one payment and get 10%% bonus!",
		"topup_custom_amount": "To top up your account, enter the top-up amount!\nFor example:\n250,\n500",
		"topup_enter_contract": "Enter the contract number you want to top up!",
		"no_contract":        "You don't have a saved contract. Please contact support.",
		"invalid_amount":     "Invalid amount format. Please enter a number.",
		"invalid_amount_minimum": "Minimum top-up amount is 0.1$ in UAH at NBU rate\nEnter another amount or return to main menu",
		"invalid_contract_format": "Invalid contract number format for top-up",
		"contract_not_found": "Contract not found in the system",
		"invoice_title":     "Top up %.2f UAH",
		"invoice_description": "Top up account %s by %.2f UAH",
		"invoice_label":     "Top up account %s",
		
		// Temporary payment
		"time_pay_success":   "Internet access unblocked for 24 hours!\nAccount topped up by %.2f UAH for 24 hours! You can now return to the main menu",
		"time_pay_failed":    "You cannot use temporary payment!\nTemporary payment can be used once per month!",
		"time_pay_ban":       "Hello! For inquiries, please use our support email support@infoaura.com.ua",
		
		// Profile
		"profile_title":     "Your Profile",
		"profile_id":        "ID: %d",
		"profile_username":  "Username: %s",
		"profile_name":      "Name: %s %s",
		"profile_language":  "Language: %s",
		"profile_balance":  "Balance: %.2f UAH",
		
		// Balance
		"balance_title":     "Your Balance",
		"balance_amount":    "Current balance: %.2f UAH",
		"balance_topup":     "Top up balance",
		"balance_topup_amount": "Enter amount to top up (minimum 10 UAH):",
		"balance_topup_success": "Balance successfully topped up by %.2f UAH",
		"balance_topup_failed": "Failed to top up balance. Please try again later.",
		
		// Services
		"services_title":    "Your Services",
		"services_list":     "Service list:",
		"services_none":     "You have no active services",
		"service_active":    "Active",
		"service_suspended": "Suspended",
		"service_inactive":  "Inactive",
		
		// Support
		"support_title":     "Support",
		"support_message":  "Enter your question or describe the problem:",
		"support_sent":     "Your request has been sent. We will contact you shortly.",
		"support_history":  "Request history",
		"support_already_active": "You are already in an active chat. Please wait for manager's response.",
		"support_waiting":  "You entered the support chat. Please wait for manager connection.",
		"support_waiting_manager": "Please wait for manager connection.",
		"support_manager_connected": "Manager connected to chat. You can now send your messages.",
		"support_chat_ended": "Chat ended. Thank you for your inquiry!",
		"support_already_connected": "This user is already in an active chat with another manager.",
		"support_no_active_chat": "You are not connected to any active chat.",
		
		// Language
		"language_title":    "Language Selection",
		"language_changed":  "Language changed to English",
		"language_ua":       "Українська",
		"language_en":       "English",
		
		// Admin
		"admin_menu":                      "Admin Panel",
		"admin_panel_title":               "Admin Panel",
		"admin_unauthorized":              "You don't have access to the admin panel",
		"admin_stats":                     "Total users: %d\nUsers with contracts: %d",
		"admin_users":                     "User Management",
		"admin_user_edit":                 "Edit User",
		"admin_search_menu":                "Select search method:",
		"admin_search_contract":           "Search by contract",
		"admin_search_phone":               "Search by phone",
		"admin_search_name":                "Search by name",
		"admin_search_address":             "Search by address",
		"admin_enter_contract":             "Enter contract number:",
		"admin_enter_phone":                "Enter phone number:",
		"admin_enter_name":                 "Enter name or surname:",
		"admin_enter_address":              "Enter address:",
		"admin_user_not_found":             "User not found",
		"admin_user_info_title":            "User Information",
		"admin_user_id":                    "User ID",
		"admin_username":                   "Username",
		"admin_balance":                    "Balance",
		"admin_status":                     "Status",
		"admin_contract":                   "Contract",
		"admin_change_balance":             "Change balance",
		"admin_temporary_payment":          "Temporary payment",
		"admin_answer_user":                "Answer user",
		"admin_enter_balance_amount":      "Enter amount to change balance:",
		"admin_invalid_amount":             "Invalid amount format",
		"admin_balance_updated":            "Balance changed by %.2f UAH. New balance: %.2f UAH",
		"admin_contract_not_found":         "Contract not found",
		"admin_temporary_payment_success":  "Temporary payment activated successfully. Amount: %.2f UAH",
		"admin_temporary_payment_failed":  "Failed to activate temporary payment",
		"admin_broadcast":                  "Broadcast Messages",
		"admin_broadcast_message":          "Enter message to broadcast:",
		"admin_broadcast_sent":             "Broadcast sent to %d users",
		"admin_enter_phone_or_userid":     "Enter phone number or user ID:",
		"admin_enter_message_text":         "Enter message text (or send photo/document/video):",
		"admin_message_sent":                "Message sent",
		"admin_enter_answer_text":          "Enter answer text:",
		"admin_answer_sent":                 "Answer sent",
		"admin_message_history":            "Message History",
		"admin_enter_message_count":        "Enter number of messages to view (1-100, default 10):",
		"admin_invalid_count":              "Invalid count. Enter a number from 1 to 100",
		"admin_no_messages":                "No messages found",
		"admin_users_found":                "Users found: %d",
		"admin_outage":                     "Outage Management",
		"admin_outage_create":              "Create Outage",
		"admin_outage_location":            "Enter outage location:",
		"admin_outage_description":         "Enter outage description:",
		"admin_outage_created":              "Outage created",
		
		// Outages
		"outage_title":      "Current Outages",
		"outage_none":       "No active outages",
		"outage_location":   "Location: %s",
		"outage_description": "Description: %s",
		"outage_status":     "Status: %s",
		
		// Help
		"help_title":        "Available commands:",
		"help_start":        "/start - Start working with the bot",
		"help_profile":      "/profile - View profile",
		"help_balance":      "/balance - Check balance and top up",
		"help_services":     "/services - View services",
		"help_support":      "/support - Contact support",
		"help_language":     "/language - Change language",
		"help_help":         "/help - Show this help",

		// Balance notifications
		"balance_notification_message": "⚠️ Warning! Your balance is low (%.2f UAH). Possible service blocking on the 12th. We recommend topping up your account.",
	}
}

