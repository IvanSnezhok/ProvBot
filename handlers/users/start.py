import logging

import asyncpg
from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters.builtin import CommandStart, Text, CommandHelp
from aiogram.types import CallbackQuery, ReplyKeyboardRemove, InlineKeyboardMarkup, InlineKeyboardButton

from data.config import ADMINS
from keyboards.default.buttons import tel_button, return_button, client_request, unknown_request_button
from keyboards.inline.callback_datas import start_callback
from keyboards.inline.start_keyboard import choice_lang
from loader import dp, db
from middlewares import _, __
from states.get_client import Client, Request
from utils.db_api import database
from utils.format_number import format_number


@dp.message_handler(CommandStart())
async def bot_start(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await message.answer(text=_("Привіт, {}!\n"
                                "Бот працює в тестовому режимі").format(message.from_user.full_name),
                         reply_markup=ReplyKeyboardRemove())
    try:
        await db.add_user(
            full_name=message.from_user.full_name,
            username=message.from_user.username,
            telegram_id=message.from_user.id
        )
    except asyncpg.exceptions.UniqueViolationError:
        await db.select_user(telegram_id=message.from_user.id)

    await message.answer(text=_("Оберіть зручну для Вас мову!"),
                         reply_markup=choice_lang
                         )


@dp.message_handler(CommandHelp())
async def help_message(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    msg = await message.answer(
        text=_("@infoaura_bot - бот мережі інтернет-провайдера Інфоаура.\n"
               "Бот призначений для доступного управлінням послугами, поповнення рахунку і виклику фахівця для вирішення локальних проблем Клієнта.\n"
               "Пропозиції та зауваження по роботі бота просимо писати на email: bot@infoaura.com.ua."),
        reply_markup=return_button)
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(state='get_phone')
async def get_phone_state(message: types.Message):
    text = _("Треба натиснути на кнопку щоб передати номер телефону")
    await message.answer(text=text, reply_markup=tel_button)
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)


@dp.callback_query_handler(start_callback.filter(lang=["RU", "UA", "EN"]))
async def lang_reply(call: CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await db.set_lang(call.data[7:].lower(), call.from_user.id)
    await call.answer()
    msg = await call.message.edit_text(
        text=_(
            "Ви обрали {}\nТепер відправте, будь ласка, свій номер телефону, щоб знайти Ваш обліковий запис у нашому білінгу",
            locale=call.data[7:].lower()).format(
            call.data[7:])
    )
    await db.message("BOT", 10001, msg.html_text, msg.date)
    msg1 = await call.message.answer(text=_("Кнопка для цього знизу", locale=call.data[7:].lower()),
                                     reply_markup=tel_button)
    await db.message("BOT", 10001, msg1.html_text, msg1.date)
    await state.set_state("get_phone")


@dp.message_handler(content_types=types.ContentType.CONTACT, state="get_phone")
async def ua_tel_get(message: types.Message, state: FSMContext):
    await state.finish()
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    tel = message.contact.phone_number
    tel = format_number(tel)
    await db.update_phone_number(tel, message.from_user.id)
    await database.search_query(tel)
    for admin in ADMINS:
        await dp.bot.send_message(admin, text=f"Новий клієнт: {message.from_user.full_name}, {message.from_user.id}\n"
                                              f"З номером телефону: {tel}\n")
        if len(database.data) > 0:
            await dp.bot.send_message(admin, text=f"Клієнт знайдений в білінгу, його номер договору "
                                                  f"{database.data[2]}\n")
        else:
            await dp.bot.send_message(admin, text=f"Клієнт не знайдений в білінгу\n")
    try:
        await db.set_contract(database.data[2], message.from_user.id)
    except IndexError:
        pass
    net_on = _("Увімкнено")
    net_off = _("Вимкнено")
    print(database.data)
    if len(database.data) > 0:

        net_pause = await database.check_net_pause(database.data[2])
        if net_pause is True and database.data[4] == "on":
            msg = await message.answer(text=_("Ваш username: {}\n"
                                              "На вашому рахунку: {}\n"
                                              "Ваш номер договору: {}\n"
                                              "Ваше ПІБ: {}\n"
                                              "Стан послуги: {}\n"
                                              "Ваш пакет: {}").format(
                database.data[0], database.data[1], database.data[2], database.data[3], net_on, database.data[5]),
                reply_markup=client_request)
            await db.message("BOT", 10001, msg.html_text, msg.date)

        else:
            msg = await message.answer(text=_("Ваш username: {}\n"
                                              "На вашому рахунку: {}\n"
                                              "Ваш номер договору: {}\n"
                                              "Ваше ПІБ: {}\n"
                                              "Стан послуги: {}\n"
                                              "Ваш пакет: {}").format(
                database.data[0], database.data[1], database.data[2], database.data[3], net_off, database.data[5]),
                reply_markup=client_request)
            await db.message("BOT", 10001, msg.html_text, msg.date)

    else:
        msg = await message.answer(
            text=_("Вказаний номер телефону не знайдено у нашому білінгу\n"
                   "Якщо ви бажаєте підключитися - залиште заявку на підключення натиснувши кнопку"),
            reply_markup=unknown_request_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(Text(equals=["Головне меню", "Главное меню", "Main menu"]), state="*")
async def main_menu(message: types.Message, state: FSMContext):
    await state.finish()
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    tel = await db.select_tel(user_id=message.from_user.id)
    await database.search_query(tel)
    try:
        await db.set_contract(database.data[2], message.from_user.id)
    except IndexError:
        pass
    net_on = __("Увімкнено")
    net_off = __("Вимкнено")
    ban = await db.get_ban()
    if message.from_user.id in ban:
        await message.answer(_("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    elif len(database.data) > 0:
        if await database.check_net_pause(database.data[2]) is True and database.data[4] == "on":
            if await db.is_alarm(message.from_user.id):
                message_alarm = await db.get_alarm_message(int(database.data[6]))
                if message_alarm is None:
                    message_alarm = await db.get_alarm_message(int(database.data[2]))
                await message.answer(f"<b>{message_alarm}</b>")
                msg = await message.answer(text=_("Ваш username: {}\n"
                                                  "На вашому рахунку: {}\n"
                                                  "Ваш номер договору: {}\n"
                                                  "Ваше ПІБ: {}\n"
                                                  "Стан послуги: {}\n"
                                                  "Ваш пакет: {}").format(
                    database.data[0], database.data[1], database.data[2], database.data[3], net_on, database.data[5]),
                    reply_markup=client_request)
                await db.message("BOT", 10001, msg.html_text, msg.date)
            else:
                msg = await message.answer(text=_("Ваш username: {}\n"
                                                  "На вашому рахунку: {}\n"
                                                  "Ваш номер договору: {}\n"
                                                  "Ваше ПІБ: {}\n"
                                                  "Стан послуги: {}\n"
                                                  "Ваш пакет: {}").format(
                    database.data[0], database.data[1], database.data[2], database.data[3], net_on, database.data[5]),
                    reply_markup=client_request)
                await db.message("BOT", 10001, msg.html_text, msg.date)
        else:
            if await db.is_alarm(message.from_user.id):
                message_alarm = await db.get_alarm_message(int(database.data[6]))
                if message_alarm is None:
                    message_alarm = await db.get_alarm_message(int(database.data[2]))
                await message.answer(f"<b>{message_alarm}</b>")
                msg = await message.answer(text=_("Ваш username: {}\n"
                                                  "На вашому рахунку: {}\n"
                                                  "Ваш номер договору: {}\n"
                                                  "Ваше ПІБ: {}\n"
                                                  "Стан послуги: {}\n"
                                                  "Ваш пакет: {}").format(
                    database.data[0], database.data[1], database.data[2], database.data[3], net_off, database.data[5]),
                    reply_markup=client_request)
                await db.message("BOT", 10001, msg.html_text, msg.date)
            else:
                msg = await message.answer(text=_("Ваш username: {}\n"
                                                  "На вашому рахунку: {}\n"
                                                  "Ваш номер договору: {}\n"
                                                  "Ваше ПІБ: {}\n"
                                                  "Стан послуги: {}\n"
                                                  "Ваш пакет: {}").format(
                    database.data[0], database.data[1], database.data[2], database.data[3], net_off, database.data[5]),
                    reply_markup=client_request)
                await db.message("BOT", 10001, msg.html_text, msg.date)
    else:
        msg = await message.answer(
            text=_("Вказаний номер телефону не знайдено у нашому білінгу\n"
                   "Якщо ви бажаєте підключитися - залиште заявку на підключення натиснувши кнопку"),
            reply_markup=unknown_request_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(Text(equals=__("Повідомити про проблему")))
async def request_for_ts(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    ban = await db.get_ban()
    tel = await db.select_tel(message.from_user.id)
    if await db.is_alarm(message.from_user.id):
        await database.search_query(tel)
        if len(database.data) > 0:
            message_alarm = await db.get_alarm_message(int(database.data[6]))
            if message_alarm is None:
                message_alarm = await db.get_alarm_message(int(database.data[2]))
            await message.answer(f"<b>{message_alarm}</b>")
    elif message.from_user.id in ban:
        await message.answer(
            _("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    else:
        msg = await message.answer(text=_("Введіть ваше ПІБ, номер телефону та опишіть проблему"),
                                   reply_markup=ReplyKeyboardRemove())
        await db.message("BOT", 10001, msg.html_text, msg.date)
        await Request.first()


@dp.message_handler(state=Request.Quest)
async def tech_support_message(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    answer = message.text
    user = await db.select_user_by_id(message.from_user.id)
    async with state.proxy() as data:
        data["Заявка"] = answer
        for admin in ADMINS:
            try:
                answer_reply = InlineKeyboardMarkup()
                answer_reply.add(InlineKeyboardButton(text="Відповісти",
                                                      callback_data=f"answer {message.from_user.id}"))
                msg = await dp.bot.send_message(admin, f"Завка на виклик майстра: {data['Заявка']}")
                msg1 = await dp.bot.send_message(admin, f"Користувач: {user[1]}"
                                                        f"\nНомер телефону: {user[5]}"
                                                        f"\nНомер договору: {user[6]}"
                                                        f"\nТелеграм ІД: {user[3]}",
                                                 reply_markup=answer_reply)
                await db.message("BOT", 10001, msg1.html_text, msg1.date)
                await db.message("BOT", 10001, msg.html_text, msg.date)

            except Exception as err:
                logging.exception(err)
    await state.reset_state()
    msg = await message.answer(text=_("Заявка в опрацюванні, чекайте зв'язку\n"
                                      "Можете повернутись у головне меню скориставшись кнопкою знизу"),
                               reply_markup=return_button)
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(Text(equals=__("Залишити заявку на підключення")))
async def get_client(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    ban = await db.get_ban()
    if message.from_user.id in ban:
        await message.answer(
            _("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    else:
        msg = await message.answer(
            text=_("Введіть ПІБ та номер телефону - ми зв'яжемось з Вами для підключення"),
            reply_markup=ReplyKeyboardRemove())
        await db.message("BOT", 10001, msg.html_text, msg.date)
        await Client.first()


@dp.message_handler(state=Client.Quest)
async def request_client(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    answer = message.text
    user = await db.select_user_by_id(message.from_user.id)
    async with state.proxy() as data:
        data["Заявка"] = answer
        for admin in ADMINS:
            try:
                answer_reply = InlineKeyboardMarkup()
                answer_reply.add(InlineKeyboardButton(text="Відповісти",
                                                      callback_data=f"answer {message.from_user.id}"))
                msg = await dp.bot.send_message(admin, f"Заявка на подключение: {data['Заявка']}")
                msg1 = await dp.bot.send_message(admin, f"Користувач: {user[1]}"
                                                        f"\nНомер телефону: {user[5]}"
                                                        f"\nНомер договору: {user[6]}"
                                                        f"\nТелеграм ІД: {user[3]}",
                                                 reply_markup=answer_reply)
                await db.message("BOT", 10001, msg1.html_text, msg1.date)
                await db.message("BOT", 10001, msg.html_text, msg.date)

            except Exception as err:
                logging.exception(err)
    await state.reset_state()
    msg = await message.answer(text=_("Заявка в опрацюванні, чекайте зв'язку\n"
                                      "Можете повернутись у головне меню скориставшись кнопкою знизу"),
                               reply_markup=return_button)
    await db.message("BOT", 10001, msg.html_text, msg.date)
