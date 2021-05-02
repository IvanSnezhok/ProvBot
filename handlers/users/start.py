import logging

import asyncpg
from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters.builtin import CommandStart, Text
from aiogram.types import CallbackQuery, ReplyKeyboardRemove

from data.config import ADMINS
from keyboards.default.buttons import tel_button, return_button, client_request, unknown_request_button, lang_change
from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp, db
from states.get_client import Client, Request
from utils.db_api import database
from utils.format_number import format_number
from middlewares import _, __


@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    await message.answer(text=_("Привіт, {}!\n").format(message.from_user.full_name),
                         reply_markup=ReplyKeyboardRemove())
    try:
        await db.add_user(
            full_name=message.from_user.full_name,
            username=message.from_user.username,
            telegram_id=message.from_user.id
        )
    except asyncpg.exceptions.UniqueViolationError:
        await db.select_user(telegram_id=message.from_user.id)

    await message.answer(text=_("Оберіть зручну для вас мову!"),
                         reply_markup=choice_lang
                         )


@dp.callback_query_handler(start_callback.filter(lang=["RU", "UA", "EN"]))
async def lang_reply(call: CallbackQuery, callback_data: dict):
    await db.set_lang(call.data[7:].lower(), call.from_user.id)
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(
        text=_("Ви обрали {}\n Тепер відпрате, будь ласка, свій контакт, щоб знайти вас у нашому білінгу",
               locale=call.data[7:].lower()).format(
            call.data[7:])
    )
    await call.message.answer(text=_("Кнопка для цього знизу", locale=call.data[7:].lower()), reply_markup=tel_button)


@dp.message_handler(content_types=types.ContentType.CONTACT)
async def ua_tel_get(message: types.Message):
    tel = message.contact.phone_number
    tel = format_number(tel)
    await db.update_phone_number(tel, message.from_user.id)
    await database.search_query(tel)
    logging.info(database.data)   # вывод результата поиска
    if len(database.data) > 0:
        await message.answer(text=_("Ваш username: {}\n"
                                    "На вашому рахунку: {}\n"
                                    "Ваш номер договору: {}\n"
                                    "Ваше ПІБ: {}\n"
                                    "Стан послуги: {}\n"
                                    "Ваш пакет: {}").format(
            database.data[0], database.data[1], database.data[2], database.data[3], database.data[4], database.data[5]),
            reply_markup=client_request)
    else:
        await message.answer(text=_("Ви не зареєстровані у нашому білінгу\n"
                                    "Якщо ви хочете підключитися можете залишити заявку на підключення натиснувши "
                                    "кнопку\n "),
                             reply_markup=unknown_request_button)


@dp.message_handler(Text(equals=["Головне меню", "Главное меню", "Main menu"]))
async def main_menu(message: types.Message):
    tel = await db.select_tel(user_id=message.from_user.id)
    await database.search_query(tel)
    try:
        await db.set_contract(database.data[2], message.from_user.id)
    except IndexError:
        pass
    if len(database.data) > 0:
        await message.answer(text=_("Ваш username: {}\n"
                                    "На вашому рахунку: {}\n"
                                    "Ваш номер договору: {}\n"
                                    "Ваше ПІБ: {}\n"
                                    "Стан послуги: {}\n"
                                    "Ваш пакет: {}").format(
            database.data[0], database.data[1], database.data[2], database.data[3], database.data[4], database.data[5]),
            reply_markup=client_request)
    else:
        await message.answer(text=_("Ви не зареєстровані у нашому білінгу\n"
                                    "Якщо ви хочете підключитися можете залишити заявку на підключення натиснувши "
                                    "кнопку\n "),
                             reply_markup=unknown_request_button,)


@dp.message_handler(Text(equals=__("Залишити заявку на виклик спеціаліста")))
async def request_for_ts(message: types.Message):
    await message.answer(text=_("Введіть ваше ПІБ, номер телефону та опишіть вашу проблему"),
                         reply_markup=ReplyKeyboardRemove())
    await Request.first()


@dp.message_handler(state=Request.Quest)
async def tech_support_message(message: types.Message, state: FSMContext):
    answer = message.text
    async with state.proxy() as data:
        data["Заявка"] = answer
        for admin in ADMINS:
            try:
                await dp.bot.send_message(admin, f"Заявка на майстра від клієнта: {data['Заявка']}")

            except Exception as err:
                logging.exception(err)
    await state.reset_state()
    await message.answer(text=_("Ваша заявка в опрацюванні, чекайте зв'язку\n"
                                "Можете повернутись у головне меню скориставшись кнопці знизу"),
                         reply_markup=return_button)


@dp.message_handler(Text(equals=__("Залишити заявку на підключення")))
async def get_client(message: types.Message):
    await message.answer(text=_("Введдіть ваше ПІБ та номер телефону, ми зв'яжемось з вами для обговорення вашого "
                                "підключення\n"), reply_markup=ReplyKeyboardRemove())
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
    await message.answer(text=_("Ваша заявка в опрацюванні, чекайте зв'язку\n"
                                "Можете повернутись у головне меню скориставшись кнопці знизу"),
                         reply_markup=return_button)

