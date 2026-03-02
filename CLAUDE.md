# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ProvBot is a Python/aiogram 2.x Telegram bot for an ISP (internet service provider). It handles user registration, billing integration, balance top-up via Telegram Payments, live support chat, admin panel, scheduled balance notifications, and SMS relay via Gmail API. A Go migration exists on the `dev2.0` branch.

## Build & Run Commands

```bash
python app.py                    # Run bot (long-polling mode)
docker build -t provbot .        # Build Docker image (python:3.9-slim-bullseye)
pip install -r requirements.txt  # Install dependencies
pybabel compile -d data/locales -D prov_bot  # Compile i18n translations
```

No tests exist in this codebase.

## Architecture

### Request Flow

```
Telegram API (long-polling) → executor.start_polling(dp)
  → Middleware chain: ThrottlingMiddleware → ACLMiddleware (i18n)
    → Decorator-based handlers matched by priority:
        1. IDFilter(ADMINS) commands  (panel.py — loaded first)
        2. /start, /help              (start.py)
        3. Text button matchers       (pay_bill.py, time_pay.py, etc.)
        4. FSM state handlers         (per-module)
        5. ContentType.ANY fallback   (echo.py — loaded last)
```

Handler registration order is determined by **import order** in `handlers/users/__init__.py`. The fallback handler in `echo.py` must load last.

### Dependency Wiring

`loader.py` creates all shared singletons imported by modules directly:
- `bot` — `Bot` instance (HTML parse mode)
- `dp` — `Dispatcher` with `MemoryStorage` (FSM states lost on restart)
- `db` — `Database` (PostgreSQL via asyncpg pool)
- `scheduler` — `AsyncIOScheduler` (APScheduler)

### Dual Database

- **PostgreSQL** (asyncpg pool) — bot's own data: `users`, `messages`, `alarm`, `bill_check`, `chats`, `user_clicks`, `settings`. Schema auto-created in `on_startup` via `CREATE TABLE IF NOT EXISTS`.
- **MySQL** (aiomysql, no pooling — new connection per query) — legacy billing system (cp1251 charset). **Readable and writable**: `pay_balance()` updates balance and inserts payment records, `t_pay()` creates temporary 24-hour credit entries.

### Key Files

| Purpose | Path |
|---|---|
| Entry point + startup | `app.py` |
| DI wiring (bot, dp, db, scheduler) | `loader.py` |
| Config (env vars via environs) | `data/config.py` |
| PostgreSQL access | `utils/db_api/postgresql.py` |
| MySQL billing access | `utils/db_api/database.py` |
| FSM state definitions | `states/get_client.py` |
| Handler import order | `handlers/users/__init__.py` |
| Admin panel (~700 lines) | `handlers/users/panel.py` |
| Payment flow | `handlers/users/pay_bill.py` |
| Support live chat | `handlers/users/support_chat.py` |
| Temporary credit | `handlers/users/time_pay.py` |
| Billing search utils | `utils/misc/find_in_bill.py` |
| SMS/Gmail relay | `utils/misc/sms_message.py` |
| Balance notifications | `utils/misc/debt_notification.py` |

### Key Patterns

- **Adding a handler**: create function with `@dp.message_handler(...)` or `@dp.callback_query_handler(...)` decorator in the appropriate module, ensure module is imported in `handlers/users/__init__.py` in the correct order
- **i18n**: `_("string")` for immediate translation, `__("string")` for lazy translation (used in keyboards). Locale files in `data/locales/{lang}/LC_MESSAGES/prov_bot.po`. `ACLMiddleware` reads user language from PostgreSQL on every update.
- **Admin detection**: `IDFilter(ADMINS)` filter, where `ADMINS` is a list of Telegram IDs from `.env`
- **FSM states**: defined in `states/get_client.py` as `StatesGroup` classes. Additional inline states used as bare strings in handlers (`'get_phone'`, `'invoice_payload'`, etc.)
- **Live support chat**: module-level `active_chats = {}` dict (userID → adminID) in `handlers/users/support_chat.py`. Admin uses `/connect <user_id>` to join, `/end_chat` to leave.
- **Payment stub toggle**: admin callback `toggle_payment_stub` → `db.toggle_payment_disabled()` in `settings` table → checked in `pay_bill.py` before processing payments

### Known Gotchas

- **Module-level mutable state in `utils/db_api/database.py`**: `data = []`, `plan = []`, `time_pay_b = []` are globals mutated by `search_query()` and `t_pay()`. Not thread-safe — concurrent calls can collide.
- **No MySQL connection pooling**: every function in `database.py` opens a new `aiomysql.connect()` and closes it after use.
- **SQL injection risk**: several MySQL queries use f-strings instead of parameterized queries (e.g., `pay_balance()`, `balance_change()`).
- **`check_contract_exists()`** fetches ALL contracts from MySQL to check membership — O(n) on every call.
- **MemoryStorage**: all FSM states are in-process RAM, lost on restart.

## Configuration

All config via `.env` file (loaded by environs). See `.env_orig.dist` for the full list. Key vars: `BOT_TOKEN`, `PROVIDER_TOKEN`, `ADMINS` (comma-separated Telegram IDs), `DB_*` (PostgreSQL), `BILL_*` (MySQL billing).

## Scheduler Jobs

- `send_message_sms` — every 3 minutes, polls Gmail API for unread emails, relays as Telegram messages or SMS via smsukraine.com.ua
- `schedule_debt_notification` — cron on 9th of month at 23:45, notifies users with low balance

## Language

Project documentation and comments are mixed Ukrainian/Russian. Code identifiers are in English. User-facing strings are translated via i18n (`_()` / `__()`).
