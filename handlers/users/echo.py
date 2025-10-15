import logging

from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Command
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton

from data.config import ADMINS
from keyboards.default.buttons import return_button, tel_button, client_request, unknown_request_button
from loader import dp, db
from utils.db_api.database import tel_by_group, account_show

from middlewares import _, __

from utils.db_api import database


@dp.message_handler(Command('phone'))
async def phones_message(message: types.Message, state: FSMContext):
    phones = await tel_by_group()
    await message.answer(phones, reply_markup=return_button)


@dp.message_handler(content_types=types.ContentTypes.ANY, state=None)
@dp.message_handler(Command(['help']))
async def bot_echo(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    tel = await db.select_tel(message.from_user.id)
    ban = await db.get_ban()

    print(ban)
    if await db.is_alarm(message.from_user.id):
        await database.search_query(tel)
        if len(database.data) > 0:
            message_alarm = await db.get_alarm_message(int(database.data[6]))
            await message.answer(message_alarm, reply_markup=return_button)
    elif message.from_user.id in ban:
        await message.answer(
            _("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    elif tel:
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
            # for admin in ADMINS:
            #     msg = await dp.bot.send_message(admin, text=_("Нове запитання з послуги клієнта\n")
        else:
            msg = await message.answer(
                text=_("Вказаний номер телефону не знайдено у нашому білінгу\n"
                       "Якщо ви бажаєте підключитися - залиште заявку на підключення натиснувши кнопку"),
                reply_markup=unknown_request_button)
            await db.message("BOT", 10001, msg.html_text, msg.date)
        try:
            for admin in ADMINS:
                answer_reply = InlineKeyboardMarkup()
                answer_reply.add(InlineKeyboardButton(text="Відповісти",
                                                      callback_data=f"answer {message.from_user.id}"))
                if message.content_type == 'text':
                    msg = await dp.bot.send_message(chat_id=admin,
                                                    text=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                         f"Текст сообщения: {message.text}\n"
                                                         f"Телефон: {await db.select_tel(message.from_user.id)}\n"
                                                         f"Номер договору: {database.data[2]}\n",
                                                    parse_mode='HTML',
                                                    reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'photo':
                    msg = await dp.bot.send_photo(chat_id=admin,
                                                  photo=message.photo[-1].file_id,
                                                  caption=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                          f"Текст сообщения: {message.caption}\n"
                                                          f"Телефон: {await db.select_tel(message.from_user.id)}\n"
                                                          f"Номер договору: {database.data[2]}\n",
                                                  parse_mode='HTML',
                                                  reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'document':
                    msg = await dp.bot.send_document(chat_id=admin,
                                                     document=message.document.file_id,
                                                     caption=f"Сообщения от пользователя:"
                                                             f" {message.from_user.full_name}\n"
                                                             f"Телефон: {await db.select_tel(message.from_user.id)}\n"
                                                             f"Номер договору: {database.data[2]}\n",
                                                     parse_mode='HTML',
                                                     reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'video':
                    msg = await dp.bot.send_video(chat_id=admin,
                                                  video=message.video.file_id,
                                                  caption=f"Сообщения от пользователя:"
                                                          f" {message.from_user.full_name}\n"
                                                          f"Телефон: {await db.select_tel(message.from_user.id)}\n"
                                                          f"Номер договору: {database.data[2]}\n",
                                                  parse_mode='HTML',
                                                  reply_markup=answer_reply)
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
        try:
            for admin in ADMINS:
                answer_reply = InlineKeyboardMarkup()
                answer_reply.add(InlineKeyboardButton(text="Відповісти",
                                                      callback_data=f"answer {message.from_user.id}"))
                if message.content_type == 'text':
                    msg = await dp.bot.send_message(chat_id=admin,
                                                    text=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                         f"Текст сообщения: {message.text}\n"
                                                         f"Пользователь без телефона",
                                                    parse_mode='HTML',
                                                    reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'photo':
                    msg = await dp.bot.send_photo(chat_id=admin,
                                                  photo=message.photo[-1].file_id,
                                                  caption=f"Сообщения от пользователя: {message.from_user.full_name}\n"
                                                          f"Пользователь без телефона",
                                                  parse_mode='HTML',
                                                  reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'document':
                    msg = await dp.bot.send_document(chat_id=admin,
                                                     document=message.document.file_id,
                                                     caption=f"Сообщения от пользователя:"
                                                             f" {message.from_user.full_name}\n"
                                                             f"Пользователь без телефона",
                                                     reply_markup=answer_reply)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                elif message.content_type == 'video':
                    msg = await dp.bot.send_video(chat_id=admin,
                                                  video=message.video.file_id,
                                                  caption=f"Сообщения от пользователя:"
                                                          f" {message.from_user.full_name}\n"

                                                          f"Пользователь без телефона",
                                                  parse_mode='HTML',
                                                  reply_markup=answer_reply)
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
