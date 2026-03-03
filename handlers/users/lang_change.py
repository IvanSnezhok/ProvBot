from loader import dp, db

from aiogram.dispatcher.filters.builtin import Text
from aiogram import types
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton

from middlewares import _, __
from middlewares.language_middleware import invalidate_lang_cache
from keyboards.default.buttons import lang_change


@dp.message_handler(Text(equals=__("Змінити мову")))
async def change_lang(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    ban = await db.get_ban()
    if message.from_user.id in ban:
        await message.answer(
            _("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    else:
        msg = await message.answer(text=_("Оберіть мову"), reply_markup=lang_change)
        await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(Text(equals=["🇷🇺 RU", "🇺🇸 EN", "🇺🇦 UA"]))
async def changed_lang(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await db.set_lang(message.text[3:].lower(), message.from_user.id)
    invalidate_lang_cache(message.from_user.id)
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

    msg = await message.answer(
        text=_("Ви обрали {}\nТепер можете перейти у головне меню",
               locale=message.text[3:].lower()).format(message.text[3:]),
        reply_markup=return_button)
    await db.message("BOT", 10001, msg.html_text, msg.date)
