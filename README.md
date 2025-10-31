# ProvBot - Telegram Bot для Провайдера

Telegram бот для провайдера, написаний на Go, з підтримкою мультимовності, інтеграцією з білінговою системою та адмін панеллю.

## Функціонал

### Для користувачів:
- Реєстрація та управління профілем
- Перевірка балансу та поповнення через Telegram Payments
- Перегляд послуг та їх статусів
- Звернення до техпідтримки
- Перегляд поточних аварій
- Мультимовність (українська та англійська)

### Для адміністраторів:
- Розсилка повідомлень користувачам
- Створення та управління аваріями за локацією
- Редагування користувачів та їх статусів через білінгову систему
- Управління послугами користувачів

## Вимоги

- Go 1.21 або вище
- PostgreSQL 12+ (для даних бота)
- MySQL 5.7+ (для білінгової системи)
- Telegram Bot Token

## Встановлення

1. Клонуйте репозиторій:
```bash
git clone <repository-url>
cd ProvBot
```

2. Встановіть залежності:
```bash
go mod download
```

3. Створіть файл `.env` на основі `.env.example`:
```bash
cp .env.example .env
```

4. Налаштуйте змінні середовища в `.env`:
```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=provbot
POSTGRES_PASSWORD=your_password
POSTGRES_DB=provbot_db
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=billing_user
MYSQL_PASSWORD=your_password
MYSQL_DB=billing_db
ADMIN_TELEGRAM_IDS=123456789,987654321
```

5. Створіть базу даних PostgreSQL та виконайте міграції:
```bash
# Створіть базу даних
createdb provbot_db

# Виконайте міграції (вручну або через migrate tool)
psql -d provbot_db -f internal/database/migrations/001_init_schema.up.sql
```

6. Запустіть бота:
```bash
go run cmd/bot/main.go
```

## Структура проекту

```
ProvBot/
├── cmd/bot/              # Точка входу
├── internal/
│   ├── bot/              # Логіка бота та обробники
│   ├── handlers/         # Обробники команд
│   ├── models/           # Моделі даних
│   ├── database/         # Підключення до БД та міграції
│   ├── repository/       # Репозиторії для роботи з БД
│   ├── service/          # Бізнес-логіка
│   ├── i18n/             # Інтернаціоналізація
│   └── utils/            # Утиліти (config, logger)
├── configs/              # Конфігураційні файли
└── README.md
```

## Команди бота

### Користувацькі команди:
- `/start` - Почати роботу з ботом
- `/help` - Список доступних команд
- `/profile` - Переглянути профіль
- `/balance` - Перевірити баланс та поповнити
- `/services` - Переглянути послуги
- `/support` - Звернутися до техпідтримки
- `/language` - Змінити мову

### Адмін команди:
- `/admin` - Адмін панель
- `/broadcast` - Розсилка повідомлень
- `/outage` - Управління аваріями
- `/users` - Управління користувачами

## Конфігурація

### База даних

Бот використовує дві бази даних:
- **PostgreSQL** - для зберігання даних бота (користувачі, логи, аварії)
- **MySQL** - для інтеграції з білінговою системою

### Адміністратори

Додайте Telegram ID адміністраторів у змінну `ADMIN_TELEGRAM_IDS` через кому:
```env
ADMIN_TELEGRAM_IDS=123456789,987654321
```

## Розробка

### Додавання нових команд

1. Додайте обробник у `internal/bot/handler.go`:
```go
h.handlers["/newcommand"] = h.handleNewCommand
```

2. Реалізуйте функцію обробника:
```go
func (h *BotHandler) handleNewCommand(ctx *BotContext) error {
    // Логіка обробки
    return nil
}
```

### Додавання перекладів

Додайте нові ключі перекладів у файли:
- `internal/i18n/ua.go` - українська мова
- `internal/i18n/en.go` - англійська мова

## Логування

Всі повідомлення та події логуються:
- В консоль (структурований JSON формат)
- В базу даних PostgreSQL (таблиці `message_logs` та `bot_logs`)

## Міграції БД

Міграції знаходяться в `internal/database/migrations/`. Для виконання міграцій використовуйте інструмент `migrate` або виконуйте SQL файли вручну.

## Ліцензія

[Вкажіть ліцензію]

## Автор

[Ваше ім'я]

