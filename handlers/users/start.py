import logging
import asyncio

import asyncpg
from aiogram import types
from aiogram.dispatcher.filters.builtin import CommandStart
from aiogram.types import CallbackQuery


from keyboards.default.buttons import tel_button
from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp, db
from utils.db_api import database
from utils.format_number import format_number


@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    try:
        user = await db.add_user(
            full_name=message.from_user.full_name,
            username=message.from_user.username,
            telegram_id=message.from_user.id
        )
    except asyncpg.exceptions.UniqueViolationError:
        await db.select_user(telegram_id=message.from_user.id)

    await message.answer(text=f"Привіт, {message.from_user.full_name}!\n"
                              f"Оберіть зручну для вас мову!",
                         reply_markup=choice_lang
                         )


@dp.callback_query_handler(start_callback.filter(lang="ENG"))
async def eng_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    await db.set_lang("ENG", call.from_user.id)
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="You choice english")


@dp.callback_query_handler(start_callback.filter(lang="RU"))
async def ru_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    await db.set_lang("RU", call.from_user.id)
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="Вы выбрали русский")


@dp.callback_query_handler(start_callback.filter(lang="UA"))
async def ua_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    await db.set_lang("UA", call.from_user.id)
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="Ви обрали українську\n"
                                 "Тепер відпрате будь ласка свій контакт щоб знайти вас у нашому білінгу")
    await call.message.answer(text="Кнопка для цього знизу", reply_markup=tel_button)


@dp.message_handler(content_types=types.ContentType.CONTACT)
async def ua_tel_get(message: types.Message):
    tel = message.contact.phone_number
    tel = format_number(tel)
    await database.search_query(tel)
    print(database.data)
    if len(database.data) > 0:
        await message.answer(text=f"""
        Ваш username: {database.data[0]}\n
На вашому рахунку: {database.data[1]}\n
Ваш номер договору: {database.data[2]}\n
Ваше ФИО: {database.data[3]}\n
Стан послуги: {database.data[4]}\n
Ваш пакет: {database.data[5]}""")
    else:
        await message.answer(text="Ви не зареєстровані у нашому білінгу")
