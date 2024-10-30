from aiogram.types import ReplyKeyboardMarkup, KeyboardButton

from middlewares import __

tel_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Відправити номер телефону"),
                       request_contact=True),
        KeyboardButton(text=__("Поповнити рахунок"))
    ]
], one_time_keyboard=True)

request_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Головне меню")),
        KeyboardButton(text=__("Залишити заявку на підключення"))
    ]
], one_time_keyboard=True)

unknown_request_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Залишити заявку на підключення")),
        KeyboardButton(text=__("Змінити мову")),
        KeyboardButton(text=__("Поповнити рахунок"))
    ]
])

client_request = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Тимчасовий платіж")),
        KeyboardButton(text=__("Поповнити рахунок")),
        KeyboardButton(text=__("Змінити мову"))
    ],
    [
        KeyboardButton(text=__("Повідомити про проблему")),
        KeyboardButton(text=__("Підключити друга"))
    ],
    [
        KeyboardButton(text=__("🛒 Магазин"))
    ]
], one_time_keyboard=True)

return_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Головне меню"))
    ]
], one_time_keyboard=True)

lang_change = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text="🇺🇦 UA"),
        KeyboardButton(text="🇺🇸 EN"),
        KeyboardButton(text="🇷🇺 RU")

    ]
], one_time_keyboard=True)

time_pay = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text=__("Тимчасовий платіж"))
    ]
], one_time_keyboard=True)
