from aiogram.types import ReplyKeyboardMarkup, KeyboardButton

tel_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text="📱",
                       request_contact=True)
    ]
], one_time_keyboard=True)

request_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text="Головне меню"),
        KeyboardButton(text="Залишити заявку на підключення")
    ]
], one_time_keyboard=True)

client_request = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text="Залишити заявку на майтсра")
    ]
], one_time_keyboard=True)

return_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
    [
        KeyboardButton(text="Головне меню")
    ]
], one_time_keyboard=True)
