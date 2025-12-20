package i18n

// GetRUTranslations returns Russian translations
func GetRUTranslations() map[string]string {
	return map[string]string{
		// Common
		"welcome":         "Добро пожаловать в бот провайдера!",
		"error":           "Произошла ошибка. Попробуйте позже.",
		"unknown_command": "Неизвестная команда. Используйте /help для списка команд.",
		"cancel":          "Отменить",
		"back":            "Назад",
		"confirm":         "Подтвердить",
		"select":          "Выбрать",

		// Start command
		"start_message":      "Здравствуйте! Я бот провайдера. Я помогу вам управлять услугами, проверять баланс и обращаться в техподдержку.\n\nВыберите удобный для Вас язык!",
		"start_registered":   "Вы уже зарегистрированы! Используйте /help для списка команд.",
		"phone_request":      "Вы выбрали %s\nТеперь отправьте, пожалуйста, свой номер телефона, чтобы найти Вашу учетную запись в нашем биллинге",
		"send_phone":         "Отправить номер телефона",
		"welcome_registered": "Добро пожаловать! Вы успешно зарегистрированы. Используйте меню для навигации.",
		"menu_topup":         "Пополнить счет",
		"menu_support":       "Чат с тех. поддержкой",
		"menu_time_pay":      "Временный платеж",
		"menu_back":          "Назад",

		// Payment
		"topup_notice":            "Обратите внимание, что здесь вы можете пополнить только свой лицевой счет!",
		"topup_promotion":         "Действует акция - пополни счет на 6 месяцев одним платежом и получите 10%% от суммы пополнения!",
		"topup_custom_amount":     "Для пополнения счета введите сумму пополнения!\nНапример:\n250,\n500",
		"topup_enter_contract":    "Введите номер договора, который хотите пополнить!",
		"no_contract":             "У вас нет сохраненного договора. Обратитесь в поддержку.",
		"invalid_amount":          "Неверный формат суммы. Введите число.",
		"invalid_amount_minimum":  "Минимальная сумма пополнения от 0.1$ в гривнах по курсу НБУ\nВведите другую сумму пополнения или можете вернуться в главное меню",
		"invalid_contract_format": "Неверно указан номер договора для пополнения",
		"contract_not_found":      "Договор не найден в системе",
		"invoice_title":           "Пополнение на %.2f грн",
		"invoice_description":     "Пополнение счета %s на %.2f гривен",
		"invoice_label":           "Пополнение счета %s",

		// Temporary payment
		"time_pay_success": "Доступ в Интернет разблокирован на 24 часа!\nСчет пополнен на %.2f грн на 24 часа! Теперь можете вернуться в главное меню",
		"time_pay_failed":  "Вы не можете использовать временный платеж!\nИспользовать временный платеж можно раз в месяц!",
		"time_pay_ban":     "Добро пожаловать! Для обращения, пожалуйста, воспользуйтесь нашим email технической поддержки support@infoaura.com.ua",

		// Profile
		"profile_title":    "Ваш профиль",
		"profile_id":       "ID: %d",
		"profile_username": "Пользователь: %s",
		"profile_name":     "Имя: %s %s",
		"profile_language": "Язык: %s",
		"profile_balance":  "Баланс: %.2f грн",

		// Balance
		"balance_title":         "Ваш баланс",
		"balance_amount":        "Текущий баланс: %.2f грн",
		"balance_topup":         "Пополнить баланс",
		"balance_topup_amount":  "Введите сумму для пополнения (минимум 10 грн):",
		"balance_topup_success": "Баланс успешно пополнен на %.2f грн",
		"balance_topup_failed":  "Не удалось пополнить баланс. Попробуйте позже.",

		// Services
		"services_title":    "Ваши услуги",
		"services_list":     "Список услуг:",
		"services_none":     "У вас нет активных услуг",
		"service_active":    "Активна",
		"service_suspended": "Приостановлена",
		"service_inactive":  "Неактивна",
		"service_on":        "Включена",
		"service_off":       "Выключена",

		// Support
		"support_title":             "Техподдержка",
		"support_message":           "Введите ваш вопрос или опишите проблему:",
		"support_sent":              "Ваше обращение отправлено. Мы свяжемся с вами в ближайшее время.",
		"support_history":           "История обращений",
		"support_already_active":    "Вы уже в активном чате. Ожидайте ответа от менеджера.",
		"support_waiting":           "Вы вошли в чат поддержки. Ожидайте подключения менеджера.",
		"support_waiting_manager":   "Пожалуйста, ожидайте подключения менеджера.",
		"support_manager_connected": "Менеджер подключился к чату. Можете писать ваши сообщения.",
		"support_chat_ended":        "Чат завершен. Спасибо за обращение!",
		"support_already_connected": "Этот пользователь уже в активном чате с другим менеджером.",
		"support_no_active_chat":    "Вы не подключены ни к одному активному чату.",
		"support_connect_button":    "Подключиться к чату",
		"support_end_chat_button":   "Завершить чат",
		"support_admin_connected":   "Вы подключились к чату с пользователем. Можете отвечать на сообщения.",

		// Language
		"language_title":   "Выбор языка",
		"language_changed": "Язык изменен на русский",
		"language_ua":      "Українська",
		"language_en":      "English",
		"language_ru":      "Русский",

		// Admin
		"admin_menu":                      "Админ панель",
		"admin_panel_title":               "Админ панель",
		"admin_unauthorized":              "У вас нет прав доступа к админ панели",
		"admin_stats":                     "Количество пользователей: %d\nКоличество пользователей с договором: %d",
		"admin_users":                     "Управление пользователями",
		"admin_user_edit":                 "Редактировать пользователя",
		"admin_search_menu":               "Выберите способ поиска пользователя:",
		"admin_search_contract":           "Поиск по договору",
		"admin_search_phone":              "Поиск по телефону",
		"admin_search_name":               "Поиск по имени",
		"admin_search_address":            "Поиск по адресу",
		"admin_enter_contract":            "Введите номер договора:",
		"admin_enter_phone":               "Введите номер телефона:",
		"admin_enter_name":                "Введите имя или фамилию:",
		"admin_enter_address":             "Введите адрес:",
		"admin_user_not_found":            "Пользователь не найден",
		"admin_user_info_title":           "Информация о пользователе",
		"admin_user_id":                   "ID пользователя",
		"admin_username":                  "Имя пользователя",
		"admin_balance":                   "Баланс",
		"admin_status":                    "Статус",
		"admin_contract":                  "Договор",
		"admin_change_balance":            "Изменить баланс",
		"admin_temporary_payment":         "Временный платеж",
		"admin_answer_user":               "Ответить пользователю",
		"admin_enter_balance_amount":      "Введите сумму для изменения баланса:",
		"admin_invalid_amount":            "Неверный формат суммы",
		"admin_balance_updated":           "Баланс изменен на %.2f грн. Новый баланс: %.2f грн",
		"admin_contract_not_found":        "Договор не найден",
		"admin_temporary_payment_success": "Временный платеж успешно активирован. Сумма: %.2f грн",
		"admin_temporary_payment_failed":  "Не удалось активировать временный платеж",
		"admin_broadcast":                 "Рассылка сообщений",
		"admin_broadcast_message":         "Введите сообщение для рассылки:",
		"admin_broadcast_sent":            "Рассылка отправлена %d пользователям",
		"admin_enter_phone_or_userid":     "Введите номер телефона или ID пользователя:",
		"admin_enter_message_text":        "Введите текст сообщения (или отправьте фото/документ/видео):",
		"admin_message_sent":              "Сообщение отправлено",
		"admin_enter_answer_text":         "Введите текст ответа:",
		"admin_answer_sent":               "Ответ отправлен",
		"admin_message_history":           "История сообщений",
		"admin_enter_message_count":       "Введите количество сообщений для просмотра (1-100, по умолчанию 10):",
		"admin_invalid_count":             "Неверное количество. Введите число от 1 до 100",
		"admin_no_messages":               "Сообщений не найдено",
		"admin_users_found":               "Найдено пользователей: %d",
		"admin_outage":                    "Управление авариями",
		"admin_outage_create":             "Создать аварию",
		"admin_outage_location":           "Введите локацию аварии:",
		"admin_outage_description":        "Введите описание аварии:",
		"admin_outage_created":            "Авария создана",

		// Outages
		"outage_title":       "Текущие аварии",
		"outage_none":        "Активных аварий нет",
		"outage_location":    "Локация: %s",
		"outage_description": "Описание: %s",
		"outage_status":      "Статус: %s",
		"outage_warning":     "⚠️ ВНИМАНИЕ! По вашему адресу имеется авария:",

		// Help
		"help_title":    "Доступные команды:",
		"help_start":    "/start - Начать работу с ботом",
		"help_profile":  "/profile - Просмотреть профиль",
		"help_balance":  "/balance - Проверить баланс и пополнить",
		"help_services": "/services - Просмотреть услуги",
		"help_support":  "/support - Обратиться в техподдержку",
		"help_language": "/language - Изменить язык",
		"help_help":     "/help - Показать эту справку",

		// Balance notifications
		"balance_notification_message": "Внимание! Ваш баланс низкий (%.2f грн). Возможна блокировка услуг 12 числа. Рекомендуем пополнить счет.",

		// Connection request
		"connection_request":         "Оставить заявку на подключение",
		"connection_request_prompt":  "Введите ФИО и номер телефона - мы свяжемся с Вами для подключения",
		"connection_request_sent":    "Заявка в обработке, ожидайте связи\nМожете вернуться в главное меню нажав кнопку внизу",
		"report_problem":             "Сообщить о проблеме",
		"report_problem_prompt":      "Введите ваше ФИО, номер телефона и опишите проблему",

		// Connect friend
		"connect_friend":       "Подключить друга",
		"connect_friend_promo": "<b>Подключите друга к нашей сети интернет</b> - получите на свой счет сумму стоимости Вашего текущего тарифного плана. При оформлении заявки на подключение, лицу нужно указать реквизиты Вашего подключения (на выбор: номер договора %s или адрес подключения %s). После фактического подключения «Друга» к нашей сети Ваш лицевой счет будет автоматически пополнен на сумму стоимости Вашего текущего тарифного плана. При подключении также действует акция «Переход»",

		// User info display
		"user_info_format": "Ваш username: %s\nНа вашем счету: %.2f\nВаш номер договора: %s\nВаше ФИО: %s\nСостояние услуги: %s\nВаш пакет: %s",
		"not_found_billing": "Ваш номер телефона не найден в нашем биллинге\nЕсли вы хотите подключиться - оставьте заявку на подключение нажав кнопку",

		// Shop
		"menu_shop":               "Магазин",
		"shop_categories":         "Выберите категорию товаров:",
		"shop_products_in_category": "Товары в категории",
		"shop_no_products":        "В этой категории нет товаров",
		"shop_product_not_found":  "Товар не найден",
		"shop_price":              "Цена",
		"shop_order":              "Заказать",
		"shop_order_message":      "Для заказа товара '%s' обратитесь в техподдержку через чат или позвоните нам.",
	}
}
