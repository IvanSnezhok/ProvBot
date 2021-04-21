import logging

import asyncpg
from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters.builtin import CommandStart, Text
from aiogram.types import CallbackQuery, ReplyKeyboardRemove

from data.config import ADMINS
from keyboards.default.buttons import tel_button, return_button, request_button, client_request
from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp, db
from states.get_client import Client, Request
from utils.db_api import database
from utils.format_number import format_number


@dp.message_handler(Text(equals="Головне меню"))
@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    await message.answer(text=f"Привіт, {message.from_user.full_name}!\n", reply_markup=ReplyKeyboardRemove())
    try:
        await db.add_user(
            full_name=message.from_user.full_name,
            username=message.from_user.username,
            telegram_id=message.from_user.id
        )
    except asyncpg.exceptions.UniqueViolationError:
        await db.select_user(telegram_id=message.from_user.id)

    await message.answer(text=f"Оберіть зручну для вас мову!",
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
    await db.update_phone_number(tel, message.from_user.id)
    await database.search_query(tel)
    # print(database.data) # вывод результата поиска
    if len(database.data) > 0:
        await message.answer(text=f"Ваш username: {database.data[0]}\n"
                                  f"На вашому рахунку: {database.data[1]}\n"
                                  f"Ваш номер договору: {database.data[2]}\n"
                                  f"Ваше ФИО: {database.data[3]}\n"
                                  f"Стан послуги: {database.data[4]}\n"
                                  f"Ваш пакет: {database.data[5]}", reply_markup=client_request)
    else:
        await message.answer(text="Ви не зареєстровані у нашому білінгу\n"
                                  "Якщо ви хочете підключитися можете залишити заявку на підключення натиснувши "
                                  "кнопку\n "
                                  "Або можете повернутись у головне меню, для цього натисніть кнопку знизу",
                             reply_markup=request_button)


@dp.message_handler(Text(equals="Залишити заявку на майтсра"))
async def request_for_ts(message: types.Message):
    await message.answer(text="Введіть ваше ФІО та номер телефону та опишіть вашу проблему",
                         reply_markup=ReplyKeyboardRemove())
    await Request.first()


@dp.message_handler(state=Request.Quest)
async def tech_support_message(message: types.Message, state: FSMContext):
    answer = message.text
    async with state.proxy() as data:
        data["Заявка"] = answer
        for admin in ADMINS:
            try:
                await dp.bot.send_message(admin, f"Заявка на майста від клієнта: {data['Заявка']}")

            except Exception as err:
                logging.exception(err)
    await state.reset_state()
    await message.answer(text="Ваша заявка в опрацюванні, чекайте зв'язку\n"
                              "Можете повернутись у головне меню по кнопці знизу", reply_markup=return_button)


@dp.message_handler(Text(equals="Залишити заявку на підключення"))
async def get_client(message: types.Message):
    await message.answer(text="Введдіть ваше ФІО та номер телефону, ми зв'яжемось з вами для обговорення вашого "
                              "підключення\n", reply_markup=ReplyKeyboardRemove())
    await Client.first()


@dp.message_handler(state=Client.Quest)
async def request_client(message: types.Message, state: FSMContext):
    answer = message.text
    async with state.proxy() as data:
        data["Заявка"] = answer
        for admin in ADMINS:
            try:
                await dp.bot.send_message(admin, f"Заявка на подключение: {data['Заявка']}")

            except Exception as err:
                logging.exception(err)
    await state.reset_state()
    await message.answer(text="Ваша заявка в опрацюванні, чекайте зв'язку\n"
                              "Можете повернутись у головне меню по кнопці знизу", reply_markup=return_button)
