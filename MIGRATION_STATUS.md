# ProvBot: Python → Go Migration Status

## Overview

This document tracks the migration progress from the Python Telegram bot (branch `origin/ProvBot`) to the Go implementation (branches `master`/`dev2.0`).

- **Original**: Python (aiogram framework)
- **Target**: Go (go-telegram-bot-api)
- **Migration Date**: October 31, 2025 (commit `7bb16b7`)

---

## Function Mapping

### User Handlers

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `bot_start()` | handlers/users/start.py | `HandleStart()` | internal/handlers/user/start.go | ✅ Done |
| `lang_reply()` | handlers/users/start.py | `HandleLanguageCallback()` | internal/handlers/user/start.go | ✅ Done |
| `get_phone_state()` | handlers/users/start.py | (state handler) | internal/handlers/user/start.go | ✅ Done |
| `ua_tel_get()` | handlers/users/start.py | `HandleContact()` | internal/handlers/user/start.go | ✅ Done |
| `call_main_menu()` | handlers/users/start.py | `ShowMainMenu()` | internal/handlers/user/start.go | ✅ Done |
| `main_menu()` | handlers/users/start.py | `ShowMainMenu()` | internal/handlers/user/start.go | ✅ Done |
| `help_message()` | handlers/users/start.py | (in handler.go) | internal/bot/handler.go | ✅ Done |
| `request_for_ts()` | handlers/users/start.py | `HandleStartSupport()` | internal/handlers/support/chat.go | ✅ Done |
| `tech_support_message()` | handlers/users/start.py | `HandleSupportMessage()` | internal/handlers/support/chat.go | ✅ Done |
| `connect_friend()` | handlers/users/start.py | `HandleConnectFriend()` | internal/handlers/user/start.go | ✅ Done |
| `get_client()` | handlers/users/start.py | - | - | ❌ Missing |
| `request_client()` | handlers/users/start.py | `HandleConnectionRequest()` | internal/handlers/user/start.go | ✅ Done |

### Payment Handlers

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `contract_pay()` | handlers/users/pay_bill.py | `HandleTopUp()` | internal/handlers/user/pay_bill.go | ✅ Done |
| `get_invoice_payload()` | handlers/users/pay_bill.py | `HandleAmountInput()` | internal/handlers/user/pay_bill.go | ✅ Done |
| `get_invoice_contract()` | handlers/users/pay_bill.py | `HandleContractInput()` | internal/handlers/user/pay_bill.go | ✅ Done |
| `process_pre_checkout()` | handlers/users/pay_bill.py | (in handler.go) | internal/bot/handler.go | ✅ Done |
| `process_successful_pay()` | handlers/users/pay_bill.py | (in handler.go) | internal/bot/handler.go | ✅ Done |
| `time_pay()` | handlers/users/time_pay.py | `HandleTimePay()` | internal/handlers/user/time_pay.go | ✅ Done |

### Admin Handlers

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `admin_panel()` | handlers/users/panel.py | `HandleAdmin()` | internal/handlers/admin/panel.go | ✅ Done |
| `stats()` | handlers/users/panel.py | `HandleStats()` | internal/handlers/admin/panel.go | ✅ Done |
| `account_menu()` | handlers/users/panel.py | `HandleAccountMenu()` | internal/handlers/admin/users.go | ✅ Done |
| `search_account()` | handlers/users/panel.py | `HandleSearch*()` | internal/handlers/admin/users.go | ✅ Done |
| `account_menu_handler()` | handlers/users/panel.py | `Handle*Input()` | internal/handlers/admin/users.go | ✅ Done |
| `account_menu_list()` | handlers/users/panel.py | (in users.go) | internal/handlers/admin/users.go | ✅ Done |
| `admin_change_balance()` | handlers/users/panel.py | `HandleBalanceChange()` | internal/handlers/admin/billing.go | ✅ Done |
| `admin_change_balance_handler()` | handlers/users/panel.py | `HandleBalanceChangeInput()` | internal/handlers/admin/billing.go | ✅ Done |
| `admin_temporary_payment()` | handlers/users/panel.py | `HandleTemporaryPayment()` | internal/handlers/admin/billing.go | ✅ Done |
| `send_message()` | handlers/users/panel.py | `HandleSendMessage()` | internal/handlers/admin/broadcast.go | ✅ Done |
| `message_get_phone()` | handlers/users/panel.py | `HandlePhoneInput()` | internal/handlers/admin/broadcast.go | ✅ Done |
| `message_get_text()` | handlers/users/panel.py | `HandleMessageInput()` | internal/handlers/admin/broadcast.go | ✅ Done |
| `message_send_accept()` | handlers/users/panel.py | (in broadcast.go) | internal/handlers/admin/broadcast.go | ✅ Done |
| `message_send_accept_phone()` | handlers/users/panel.py | (in broadcast.go) | internal/handlers/admin/broadcast.go | ✅ Done |
| `admin_answer()` | handlers/users/panel.py | `HandleAnswer()` | internal/handlers/admin/broadcast.go | ✅ Done |
| `admin_answer_text()` | handlers/users/panel.py | `HandleAnswerInput()` | internal/handlers/admin/broadcast.go | ✅ Done |
| `message_history_start()` | handlers/users/panel.py | (in logs.go) | internal/handlers/admin/logs.go | ✅ Done |
| `message_history_get()` | handlers/users/panel.py | `HandleHistoryCount()` | internal/handlers/admin/logs.go | ✅ Done |
| `message_send_accept_sms()` | handlers/users/panel.py | - | - | ❌ Missing |

### Shop Handlers

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `show_categories()` | handlers/users/shop.py | - | - | 🚫 Removed |
| `show_category_products()` | handlers/users/shop.py | - | - | 🚫 Removed |
| `show_product_page()` | handlers/users/shop.py | - | - | 🚫 Removed |

> **Note:** Shop functionality was removed as not needed for this project.

### Support Handlers

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `HandleStartSupport()` | (from start.py) | `HandleStartSupport()` | internal/handlers/support/chat.go | ✅ Done |
| `HandleSupportMessage()` | (from start.py) | `HandleSupportMessage()` | internal/handlers/support/chat.go | ✅ Done |
| `HandleAdminConnect()` | - | `HandleAdminConnect()` | internal/handlers/support/chat.go | ✅ Done |
| `HandleAdminMessage()` | - | `HandleAdminMessage()` | internal/handlers/support/chat.go | ✅ Done |
| `HandleEndChat()` | - | `HandleEndChat()` | internal/handlers/support/chat.go | ✅ Done |

### Scheduler/Notifications

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `notify_clients()` | handlers/users/notify_balance.py | (scheduler) | internal/scheduler/notify_balance.go | ✅ Done |
| `scheduler_jobs()` | handlers/users/notify_balance.py | (in main.go) | cmd/bot/main.go | ✅ Done |

### Database/Service Functions

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `search_query()` | utils/db_api/database.py | `SearchUser()` | internal/service/billing_service.go | ✅ Done |
| `account_show()` | utils/db_api/database.py | `GetBillingUser()` | internal/service/billing_service.go | ✅ Done |
| `balance()` | utils/db_api/database.py | `GetBalance()` | internal/service/billing_service.go | ✅ Done |
| `pay_balance()` | utils/db_api/database.py | `PayBalance()` | internal/service/billing_service.go | ✅ Done |
| `t_pay()` | utils/db_api/database.py | `TemporaryPay()` | internal/service/billing_service.go | ✅ Done |
| `balance_change()` | utils/db_api/database.py | `UpdateBalance()` | internal/service/billing_service.go | ✅ Done |
| `check_net_pause()` | utils/db_api/database.py | - | - | ❌ Missing |
| `get_ban()` | utils/db_api/postgresql.py | `IsBanned()` | internal/repository/user_repo.go | ✅ Done |
| `is_alarm()` | utils/db_api/postgresql.py | `HasActiveOutageForUser()` | internal/service/outage_service.go | ✅ Done |
| `get_alarm_message()` | utils/db_api/postgresql.py | `GetOutageMessageForUser()` | internal/service/outage_service.go | ✅ Done |

### SMS Integration

| Python Function | Python File | Go Function | Go File | Status |
|----------------|-------------|------------|---------|--------|
| `send_message_sms()` | utils/misc/sms_message.py | `SendSMS()` | internal/service/sms_service.go | ✅ Done |

---

## Translations (i18n)

| Language | Python File | Go File | Status |
|----------|-------------|---------|--------|
| Ukrainian (UA) | data/locales/- | internal/i18n/ua.go | ✅ Done |
| English (EN) | data/locales/en/ | internal/i18n/en.go | ✅ Done |
| Russian (RU) | data/locales/ru/ | - | 🚫 Removed |

> **Note:** Russian language support was removed.

---

## Missing Features Summary

### High Priority
1. **Network status check** (`check_net_pause()`) - Shows "On/Off" status in menu

### Medium Priority
2. **Get client feature** (`get_client()`) - Show client info

### Completed
- ✅ Ban system integration - IsBanned() integrated into pay_bill.go and time_pay.go

### Completed/Removed
- ~~Russian translations~~ - Removed (not needed)
- ~~Shop functionality~~ - Removed (not needed)
- ✅ Alarm/outage notifications - Replaced by OutageService
- ✅ SMS integration - Implemented in sms_service.go
- ✅ "Connect friend" feature - Implemented
- ✅ "Connection request" - Implemented

---

## Implementation Progress

- [x] Core bot framework
- [x] User registration and language selection
- [x] Main menu and navigation
- [x] Payment handling (Telegram Payments)
- [x] Temporary payment feature
- [x] Support chat system
- [x] Admin panel
- [x] User search (by contract, phone, name, address)
- [x] Balance management
- [x] Message broadcasting
- [x] Message history
- [x] Scheduled balance notifications
- [x] Logging middleware
- [x] Outage/Alarm notifications (OutageService)
- [x] "Connect friend" feature
- [x] Connection request form
- [x] SMS service
- [x] Ban system integration in payment handlers
- [ ] Network status display (`check_net_pause`)
- [ ] Get client feature

---

## File Structure Comparison

### Python (origin/ProvBot)
```
ProvBot/
├── app.py                          # Entry point
├── loader.py                       # Bot initialization
├── handlers/
│   ├── users/
│   │   ├── start.py               # Registration, menu
│   │   ├── panel.py               # Admin panel
│   │   ├── pay_bill.py            # Payments
│   │   ├── time_pay.py            # Temporary payment
│   │   ├── shop.py                # Shop
│   │   ├── contact.py             # Contact
│   │   └── notify_balance.py      # Notifications
│   └── groups/
│       └── group_panel.py         # Group features
├── keyboards/                      # Keyboards
├── middlewares/                    # Middleware
├── states/                         # FSM states
├── utils/
│   ├── db_api/
│   │   ├── database.py            # MySQL (billing)
│   │   └── postgresql.py          # PostgreSQL (bot)
│   └── misc/
│       └── sms_message.py         # SMS
└── data/
    └── locales/                   # Translations
```

### Go (master/dev2.0)
```
ProvBot/
├── cmd/bot/main.go                # Entry point
├── internal/
│   ├── bot/
│   │   ├── handler.go             # Main handler
│   │   ├── middleware.go          # Middleware
│   │   └── context.go             # Bot context
│   ├── handlers/
│   │   ├── user/
│   │   │   ├── start.go           # Registration, menu
│   │   │   ├── pay_bill.go        # Payments
│   │   │   └── time_pay.go        # Temporary payment
│   │   ├── admin/
│   │   │   ├── panel.go           # Admin panel
│   │   │   ├── users.go           # User search
│   │   │   ├── billing.go         # Balance management
│   │   │   ├── broadcast.go       # Broadcasting
│   │   │   └── logs.go            # Logs
│   │   └── support/
│   │       └── chat.go            # Support chat
│   ├── repository/
│   │   ├── user_repo.go           # User repository
│   │   ├── billing_repo.go        # Billing repository
│   │   └── log_repo.go            # Log repository
│   ├── service/
│   │   ├── user_service.go        # User service
│   │   ├── billing_service.go     # Billing service
│   │   └── notification_service.go # Notifications
│   ├── models/                     # Data models
│   ├── state/                      # State management
│   ├── i18n/                       # Translations
│   ├── database/                   # DB connections
│   ├── scheduler/                  # Scheduled tasks
│   └── utils/                      # Utilities
└── deploy/                         # Deployment scripts
```

---

## Notes

- Migration started: October 31, 2025
- Last updated: December 28, 2025
- Current branch: `dev2.0`
- Shop and Russian language support removed as not needed
