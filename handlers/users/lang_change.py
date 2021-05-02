from loader import dp, db

from aiogram.dispatcher.filters.builtin import Text
from aiogram import types
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton

from middlewares import _, __
from keyboards.default.buttons import lang_change


@dp.message_handler(Text(equals=__("Змінити мову")))
async def change_lang(message: types.Message):
    await message.answer(text=_("Оберіть мову"), reply_markup=lang_change)


@dp.message_handler(Text(equals="🇺🇦 UA"))
@dp.message_handler(Text(equals="🇺🇸 EN"))
@dp.message_handler(Text(equals="🇷🇺 RU"))
async def changed_lang(message: types.Message):
    await db.set_lang(message.text[3:].lower(), message.from_user.id)
    if message.text[3:] == "UA":
        return_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
            [
                KeyboardButton(text="Головне меню")
            ]
        ], one_time_keyboard=True)
    elif message.text[3:] == "EN":
        return_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
            [
                KeyboardButton(text="Main menu")
            ]
        ], one_time_keyboard=True)
    else:
        return_button = ReplyKeyboardMarkup(resize_keyboard=True, keyboard=[
            [
                KeyboardButton(text="Главное меню")
            ]
        ], one_time_keyboard=True)

    await message.answer(
        text=_("Ви обрали {}\n Тепер можете перейти у головне меню",
               locale=message.text[3:].lower()).format(message.text[3:]),
        reply_markup=return_button)
