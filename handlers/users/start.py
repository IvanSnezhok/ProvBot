import logging

from aiogram import types
from aiogram.dispatcher.filters.builtin import CommandStart
from aiogram.types import CallbackQuery


from keyboards.default.buttons import tel_button
from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp
from utils.format_number import format_number


@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    await message.answer(text=f"Привіт, {message.from_user.full_name}!\n"
                              f"Оберіть зручну для вас мову!",
                         reply_markup=choice_lang
                         )


@dp.callback_query_handler(start_callback.filter(lang="ENG"))
async def eng_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="You choice english")


@dp.callback_query_handler(start_callback.filter(lang="RU"))
async def ru_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="Вы выбрали русский")


@dp.callback_query_handler(start_callback.filter(lang="UA"))
async def ua_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.edit_text(text="Ви обрали українську\n"
                                 "Тепер відпрате будь ласка свій контакт щоб знайти вас у нашому білінгу")
    await call.message.answer(text="Кнопка для цього знизу", reply_markup=tel_button)


@dp.message_handler(content_types=types.ContentType.CONTACT)
async def ua_tel_get(message: types.Message):
    tel = message.contact
    tel = format_number(tel.phone_number)
    print(tel)
    await message.answer(text=f"Ваш номер {tel}")