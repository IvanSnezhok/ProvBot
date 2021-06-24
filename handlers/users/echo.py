import logging

from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Command

from data.config import ADMINS
from keyboards.default.buttons import return_button, tel_button, client_request, unknown_request_button
from loader import dp, db

from middlewares import _, __


# Эхо хендлер, куда летят текстовые сообщения без указанного состояния
from utils.db_api import database


@dp.message_handler(content_types=types.ContentTypes.ANY, state=None)
@dp.message_handler(Command(['help']))
async def bot_echo(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    tel = await db.select_tel(message.from_user.id)
    if tel:
        await database.search_query(tel)
        if len(database.data) > 0:
            msg = await message.answer(text=_("Ваш username: {}\n"
                                              "На вашому рахунку: {}\n"
                                              "Ваш номер договору: {}\n"
                                              "Ваше ПІБ: {}\n"
                                              "Стан послуги: {}\n"
                                              "Ваш пакет: {}").format(
                database.data[0], database.data[1], database.data[2], database.data[3], database.data[4],
                database.data[5]),
                reply_markup=client_request)
            await db.message("BOT", 10001, msg.html_text, msg.date)
        else:
            msg = await message.answer(
                text=_("Вказаний номер телефону не знайдено у нашому білінгу\n"
                       "Якщо ви бажаєте підключитися - залиште заявку на підключення натиснувши кнопку"),
                reply_markup=unknown_request_button)
            await db.message("BOT", 10001, msg.html_text, msg.date)
        for admin in ADMINS:
            try:
                msg = await dp.bot.send_message(chat_id=admin,
                                                text=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                     f"Текст сообщения: {message.text}\n"
                                                     f"Телефон: {await db.select_tel(message.from_user.id)}")
                await db.message("BOT", 10001, msg.html_text, msg.date)
            except Exception as err:
                logging.exception(err)
    else:
        await database.search_query(tel)
        if len(database.data) > 0:
            msg = await message.answer(text=_("Ваш username: {}\n"
                                              "На вашому рахунку: {}\n"
                                              "Ваш номер договору: {}\n"
                                              "Ваше ПІБ: {}\n"
                                              "Стан послуги: {}\n"
                                              "Ваш пакет: {}").format(
                database.data[0], database.data[1], database.data[2], database.data[3], database.data[4],
                database.data[5]),
                reply_markup=client_request)
            await db.message("BOT", 10001, msg.html_text, msg.date)
        else:
            msg = await message.answer(
                text=_("Вказаний номер телефону не знайдено у нашому білінгу\n"
                       "Якщо ви бажаєте підключитися - залиште заявку на підключення натиснувши кнопку"),
                reply_markup=unknown_request_button)
            await db.message("BOT", 10001, msg.html_text, msg.date)
        for admin in ADMINS:
            try:
                msg = await dp.bot.send_message(chat_id=admin,
                                                text=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                     f"Текст сообщения: {message.text}\n"
                                                     f"Пользователь без телефона")
                await db.message("BOT", 10001, msg.html_text, msg.date)
            except Exception as err:
                logging.exception(err)

# Эхо хендлер, куда летят ВСЕ сообщения с указанным состоянием
# @dp.message_handler(state="*", content_types=types.ContentTypes.ANY)
# async def bot_echo_all(message: types.Message, state: FSMContext):
#     await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
#     state = await state.get_state()
#     await message.answer(f"Эхо в состоянии <code>{state}</code>.\n"
#                          f"\nСодержание сообщения:\n"
#                          f"<code>{message}</code>")
