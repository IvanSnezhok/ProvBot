import logging

from aiogram import types
from aiogram.dispatcher.filters.builtin import CommandStart
from aiogram.types import CallbackQuery

from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp


@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    await message.answer(text=f"Привіт, {message.from_user.full_name}!\n"
                         f"Оберіть зручну для вас мову!",
                         reply_markup=choice_lang
                         )


@dp.callback_query_handler(start_callback.filter(lang="UA"))
async def ua_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.answer(text="Ви обрали Українську")


@dp.callback_query_handler(start_callback.filter(lang="ENG"))
async def eng_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.answer(text="You choice english")


@dp.callback_query_handler(start_callback.filter(lang="RU"))
async def ru_reply(call: CallbackQuery, callback_data: dict):
    await call.answer()
    logging.info(f"callback_data = {call.data}")
    logging.info(f"callback_data = {callback_data}")
    await call.message.answer(text="Вы выбрали русский")
