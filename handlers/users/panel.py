import re

import aiohttp
import transliterate
from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import IDFilter, Text
from aiogram.types import ReplyKeyboardRemove, InlineKeyboardMarkup, InlineKeyboardButton

from data.config import ADMINS
from keyboards.default.admin import admin_keyboard, accept_message, back, accept_sms, accept_message_phone, \
    admin_account_menu, search_choice, back_inline, grp_choice, redact_alarm
from loader import dp, db
from utils.db_api.database import pay_balance, t_pay, users_with_alarm
from utils.format_number import unformat_number, number, format_text_account_admin
from utils.misc.find_in_bill import find
from utils.misc.sms_message import send_message_sms


@dp.message_handler(IDFilter(ADMINS), commands=['stats'], state='*')
async def stats(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await message.answer("Кількість користувачів: " + str(await db.count_users()) + "\n" +
                         "Кількість користувачів с договором: " + str(await db.count_contract_users()),
                         reply_markup=back)


@dp.message_handler(IDFilter(ADMINS), text="Назад", state='*')
@dp.message_handler(IDFilter(ADMINS))
async def admin_panel(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await state.finish()
    result = await find(contract=message.text)
    if result:
        if type(result) is dict:
            await message.answer(text=await format_text_account_admin(result), reply_markup=admin_account_menu)
            await state.set_data({"account": result['contract']})
        else:
            msg = await message.answer("Користувача не знайдено\n"
                                       "Панель адміністратора", reply_markup=admin_keyboard)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    else:
        msg = await message.answer("Панель адміністратора", reply_markup=admin_keyboard)
        await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS), text="back", state='*')
async def admin_panel(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await state.finish()
    msg = await call.message.answer("Панель адміністратора", reply_markup=admin_keyboard)
    await db.message(call.from_user.full_name, call.from_user.id, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS), Text(startswith='answer'), state='*')
async def admin_answer(call: types.CallbackQuery, state: FSMContext):
    await state.finish()
    user_id = call.data.split(' ')[1]
    print(user_id)
    await state.set_state('answer')
    await state.set_data({'user_id': user_id})
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.message.answer(f"Напишіть текст або надішліть фото/документ/відео"
                              f" для відправлення користувачу {user_id}", reply_markup=back)


@dp.callback_query_handler(IDFilter(ADMINS), state='*', text='message_history')
async def message_history_start(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    msg = await call.message.answer(text='Напишіть кількість останніх повідомлень що бажаєте передивитися\n'
                                         'Стандартна кількість 10 повідомлень')
    await db.message("BOT", 10001, msg.html_text, msg.date)
    await state.set_state('message_history')


@dp.message_handler(IDFilter(ADMINS), state='answer', content_types=['text', 'photo', 'document', 'video'])
async def admin_answer_text(message: types.Message, state: FSMContext):
    user_id = await state.get_data()
    print(user_id)
    await state.finish()
    if message.content_type == "text":
        await dp.bot.send_message(user_id['user_id'], message.text)
        for admin in ADMINS:
            msg = await dp.bot.send_message(admin,
                                            "Повідомлення:\n" + message.text + "\nДо:\n" + str(user_id['user_id']),
                                            reply_markup=back)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    elif message.content_type == "photo":
        await dp.bot.send_photo(user_id['user_id'], message.photo[-1].file_id)
        for admin in ADMINS:
            msg = await dp.bot.send_photo(admin, message.photo[-1].file_id, caption="До:\n" + str(user_id['user_id']),
                                          reply_markup=back)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    elif message.content_type == "document":
        await dp.bot.send_document(user_id['user_id'], message.document.file_id)
        for admin in ADMINS:
            msg = await dp.bot.send_document(admin, message.document.file_id, caption="До:\n" + str(user_id['user_id']),
                                             reply_markup=back)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    elif message.content_type == "video":
        await dp.bot.send_video(user_id['user_id'], message.video.file_id)
        for admin in ADMINS:
            msg = await dp.bot.send_video(admin, message.video.file_id, caption="До:\n" + str(user_id['user_id']),
                                          reply_markup=back)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS), text="admin_change_balance", state='*')
async def admin_change_balance(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    msg = await call.message.answer("Введіть сумму на яку хочете змінити баланс\n"
                                    "Баланс буде дорівнювати x + сумма котру вкажете, де х поточний баланс",
                                    reply_markup=back)
    await state.set_state('admin_change_balance')
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS), text="account_menu")
async def account_menu(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    msg = await call.message.answer("Виберіть як будемо шукати абонента", reply_markup=search_choice)
    await state.set_state('account_menu')
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS), text=['search_contract', 'search_phone', 'search_name', 'search_address'],
                           state='account_menu')
async def search_account(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    if call.data == 'search_contract':
        await state.set_state('search_contract')
        await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
        msg = await call.message.answer("Введіть номер договору")
        await db.message("BOT", 10001, msg.html_text, msg.date)
    elif call.data == 'search_phone':
        await state.set_state('search_phone')
        await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
        msg = await call.message.answer("Введіть номер телефону")
        await db.message("BOT", 10001, msg.html_text, msg.date)
    elif call.data == 'search_name':
        await state.set_state('search_name')
        await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
        msg = await call.message.answer("Введіть ім'я та/або прізвище")
        await db.message("BOT", 10001, msg.html_text, msg.date)
    elif call.data == 'search_address':
        await state.set_state('search_address')
        await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
        msg = await call.message.answer("Введіть вулицю, будинок та квартиру через пробіл")
        await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(IDFilter(ADMINS), state=["search_contract", "search_phone", "search_name", "search_address"])
async def account_menu_handler(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await state.set_data({"account": message.text})
    try:
        cur_state = await state.get_state()
        if cur_state == 'search_contract':
            result = await find(contract=message.text)
        elif cur_state == 'search_phone':
            result = await find(phone=message.text)
        elif cur_state == 'search_name':
            result = await find(name=message.text)
        elif cur_state == 'search_address':
            result = await find(address=message.text.split(' '))

        if type(result) is dict:
            await message.answer(text=await format_text_account_admin(result), reply_markup=admin_account_menu)
            await state.set_data({"account": result['contract']})
        elif type(result) is list and len(result) > 0:
            await state.set_state('account_menu_list')
            keyboard = types.InlineKeyboardMarkup(row_width=1)
            for i in result:
                keyboard.add(
                    types.InlineKeyboardButton(text=i['contract'] + " " + i['fio'],
                                               callback_data='account_menu_list' + i['contract']))
                await state.set_data({i['contract']: i})
            await message.answer(text="Виберіть договір", reply_markup=keyboard)
        else:
            await message.answer(text=f"Не вдалося знайти користувача",
                                 reply_markup=back)
    except Exception as e:
        print(e)
        await message.answer(text=f"Не вдалося знайти користувача з номером договору {message.text}", reply_markup=back)
        await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), state='account_menu_list')
async def account_menu_list(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await state.set_state('account_menu')
    result = await find(contract=call.data[17:])
    msg = await call.message.answer(text=await format_text_account_admin(result), reply_markup=admin_account_menu)
    await state.set_data({"account": result['contract']})
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(IDFilter(ADMINS), state='message_history')
async def message_history_get(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    contract = await state.get_data()
    user_id = await db.select_user_id_by_contract(contract['account'])
    messages = await db.get_message_history(user_id=user_id[0][0], message_count=int(message.text))
    msg = await message.answer(f'Останні {message.text} повідомлень від {contract["account"]}')
    print(messages)
    for i in messages:
        if i[0] is not None:
            msg1 = await message.answer(i[0])
            await db.message("BOT", 10001, msg1.html_text, msg1.date)
        else:
            msg2 = await message.answer('Пусте повідомлення (натискання на кнопку або відправка файлів)')
            await db.message("BOT", 10001, msg2.html_text, msg2.date)
    await db.message("BOT", 10001, msg.html_text, msg.date)
    msg3 = await message.answer('Закінчення повідомлень', reply_markup=back)
    await db.message("BOT", 10001, msg3.html_text, msg3.date)
    await state.finish()



@dp.callback_query_handler(IDFilter(ADMINS), text='admin_temporary_payment', state='*')
async def admin_temporary_payment(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    try:
        account = await state.get_data()
        account = account['account']
        if await t_pay(account):
            msg = await call.message.answer(text="Тимчасовий платіж увімкнений успішно для абонента " + account)
        else:
            msg = await call.message.answer(text="Тимчасовий платіж не вдалося ввімкнути для абонента " + account)
    except KeyError:
        contract = re.search(r'\b\d{8}\b', call.message.text)
        if await t_pay(contract[0]):
            msg = await call.message.answer(text="Тимчасовий платіж увімкнений успішно для абонента " + contract[0])
        else:
            msg = await call.message.answer(text="Тимчасовий платіж не вдалося ввімкнути для абонента " + contract[0])
    result = await find(contract=call.data[17:])
    msg1 = await call.message.answer(text=await format_text_account_admin(result), reply_markup=admin_account_menu)
    await state.set_data({"account": result['contract']})

    await db.message("BOT", 10001, msg.html_text, msg.date)
    await db.message("BOT", 10001, msg1.html_text, msg1.date)


@dp.message_handler(IDFilter(ADMINS), state="admin_change_balance")
async def admin_change_balance_handler(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    account = await state.get_data()
    account = account['account']
    print(account)
    try:
        await pay_balance(account, message.text)
        msg = await message.answer(text=f"Баланс користувача {account} поповнено на {message.text}", reply_markup=back)
        result = await find(contract=account)
        balance = result['balance']
        phone = result['telefon']
        if len(phone) > 13:
            phone = phone[:13]
            print(phone)
        print(phone)
        try:
            text = "Рахунок " + account + " поповнений на " + message.text + " На рахунку знаходиться " + str(balance)
            await send_message_sms(unformat_number(phone), text)
        except Exception as e:
            print(e)
        await state.finish()
    except Exception as e:
        print(e)
        msg = await message.answer(text=f"Не вдалося поповнити баланс користувача {account} на {message.text}",
                                   reply_markup=back)
        await state.finish()
    await db.message("BOT", 10001, msg.html_text, msg.date)


# Send via telegram or sms to user
@dp.callback_query_handler(IDFilter(ADMINS), lambda call: call.data == "panel_send_message")
async def send_message(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    text = "Напишіть мобільний телефон або телеграм ід користувача"
    await call.answer()
    msg = await dp.bot.send_message(call.from_user.id, text, reply_markup=ReplyKeyboardRemove())
    await db.message(call.from_user.full_name, call.from_user.id, text, msg.date)
    await state.set_state("send_message_phone")


@dp.message_handler(IDFilter(ADMINS), state="send_message_phone")
async def message_get_phone(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await state.update_data(phone=message.text)
    users = await db.select_all_users()
    users_phones = []
    users_id = []
    try:
        for i in range(len(users)):
            users_phones.append(unformat_number(str(users[i]['phone_number'])))
            users_id.append(users[i]['telegram_id'])
        if message.text in users_id:
            msg = await message.answer('ІД знайдений у боті')
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        elif message.text in users_phones:
            msg = await message.answer('Телефон знайдений у боті')
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        else:
            msg = await message.answer("Телефон або ІД не знайдений у боті")
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    except Exception as e:
        print(e)
        msg = await message.answer("Телефон або ІД не знайдений у боті\n"
                                   f"Виникла помилка при пошуку {e}")
        await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    msg = await message.answer("Напишіть текст повідомлення, або відправте фотографію/документ",
                               reply_markup=ReplyKeyboardRemove())
    await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    await state.set_state("send_message_text")


@dp.message_handler(IDFilter(ADMINS), state="send_message_text", content_types=['text', 'photo', 'document', 'video'])
async def message_get_text(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    if message.content_type == "text":
        await state.update_data(type='text')
        await state.update_data(text=message.text)
    elif message.content_type == "photo":
        await state.update_data(type='photo')
        await state.update_data(photo=message.photo[-1].file_id)
    elif message.content_type == "document":
        await state.update_data(type='document')
        await state.update_data(document=message.document.file_id)
    elif message.content_type == "video":
        await state.update_data(type='video')
        await state.update_data(video=message.video.file_id)
    data = await state.get_data()
    print(data["phone"])
    users = await db.select_all_users()
    users_phones = []
    users_id = []
    try:
        for i in range(len(users)):
            users_phones.append(unformat_number(str(users[i]['phone_number'])))
            users_id.append(users[i]['telegram_id'])
        if int(data['phone']) in users_id:
            msg = await message.answer('ІД знайдений у боті, відправляти?', reply_markup=accept_message)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        elif data['phone'] in users_phones:
            msg = await message.answer('Телефон знайдений у боті, відправляти?', reply_markup=accept_message_phone)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        else:
            msg = await message.answer("Телефон або ІД не знайдений у боті, відправити смс?",
                                       reply_markup=accept_sms)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
    except Exception as e:
        print(e)
        await state.update_data(phone=unformat_number(data['phone']))
        data = await state.get_data()
        print(data["phone"])
        for i in range(len(users)):
            users_phones.append(unformat_number(str(users[i]['phone_number'])))
            users_id.append(users[i]['telegram_id'])
        if int(data['phone']) in users_id:
            msg = await message.answer('ІД знайдений у боті, відправляти?', reply_markup=accept_message)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        elif data['phone'] in users_phones:
            msg = await message.answer('Телефон знайдений у боті, відправляти?', reply_markup=accept_message_phone)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)
        else:
            msg = await message.answer("Телефон або ІД не знайдений у боті, відправити смс?",
                                       reply_markup=accept_sms)
            await db.message(message.from_user.full_name, message.from_user.id, msg.html_text, msg.date)


@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_message_decline", state="send_message_text")
@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_message_accept", state="send_message_text")
async def message_send_accept(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    data = await state.get_data()
    telegram_id = data['phone']
    contract = await db.select_contract(int(telegram_id))
    contract = contract[0]['contract']
    print(telegram_id)
    await call.answer()
    type = data['type']
    if type == 'text':
        text = data['text']
        if call.data == "panel_send_message_accept":
            try:
                msg_u = await dp.bot.send_message(chat_id=int(telegram_id), text=text, parse_mode="HTML")
                await db.message("BOT", 10001, msg_u.html_text, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nТекст повідомлення: {text}"
                                                    f"\nКонтракт користувача: {contract}",
                                                    reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'photo':
        photo = data['photo']
        if call.data == "panel_send_message_accept":
            try:
                msg_u = await dp.bot.send_photo(chat_id=int(telegram_id), photo=photo)
                await db.message("BOT", 10001, msg_u.photo[-1].file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    photo_send = await dp.bot.send_photo(admin, photo=photo, caption="Фото відправлено",
                                                         reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, photo_send.photo[-1].file_id, photo_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'document':
        document = data['document']
        if call.data == "panel_send_message_accept":
            try:
                msg_u = await dp.bot.send_document(chat_id=int(telegram_id), document=document)
                await db.message("BOT", 10001, msg_u.document.file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    doc_send = await dp.bot.send_document(admin, document=document, caption="Документ відправлено",
                                                          reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, doc_send.document.file_id, doc_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'video':
        video = data['video']
        if call.data == "panel_send_message_accept":
            try:
                msg_u = await dp.bot.send_video(chat_id=int(telegram_id), video=video)
                await db.message("BOT", 10001, msg_u.video.file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    video_send = await dp.bot.send_video(admin, video=video, caption="Відео відправлено",
                                                         reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, video_send.video.file_id, video_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_message_decline_phone", state="send_message_text")
@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_message_accept_phone", state="send_message_text")
async def message_send_accept_phone(call: types.CallbackQuery, state: FSMContext):
    await call.answer()
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    phone = await state.get_data()
    phone = phone['phone']
    telegram_id = await db.select_id_by_phone(number(phone))
    telegram_id = telegram_id[0]['telegram_id']
    data = await state.get_data()
    contract = await db.select_contract(telegram_id)
    contract = contract[0]['contract']
    type = data['type']
    if type == 'text':
        text = data['text']
        if call.data == "panel_send_message_accept_phone":
            try:
                msg_u = await dp.bot.send_message(chat_id=int(telegram_id), text=text, parse_mode="HTML")
                await db.message("BOT", 10001, msg_u.html_text, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nТекст повідомлення: {text}"
                                                    f"\nКонтракт користувача: {contract}",
                                                    reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline_phone":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'photo':
        photo = data['photo']
        if call.data == "panel_send_message_accept_phone":
            try:
                msg_u = await dp.bot.send_photo(chat_id=int(telegram_id), photo=photo)
                await db.message("BOT", 10001, msg_u.photo[-1].file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    photo_send = await dp.bot.send_photo(admin, photo=photo, caption="Фото відправлено",
                                                         reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, photo_send.photo[-1].file_id, photo_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline_phone":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'document':
        document = data['document']
        if call.data == "panel_send_message_accept_phone":
            try:
                msg_u = await dp.bot.send_document(chat_id=int(telegram_id), document=document)
                await db.message("BOT", 10001, msg_u.document.file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    doc_send = await dp.bot.send_document(admin, document=document, caption="Документ відправлено",
                                                          reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, doc_send.document.file_id, doc_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline_phone":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    elif type == 'video':
        video = data['video']
        if call.data == "panel_send_message_accept_phone":
            try:
                msg_u = await dp.bot.send_video(chat_id=int(telegram_id), video=video)
                await db.message("BOT", 10001, msg_u.video.file_id, msg_u.date)
                for admin in ADMINS:
                    msg = await dp.bot.send_message(admin,
                                                    f"Повідомлення відправлено користувачу: {telegram_id}"
                                                    f"\nКонтракт користувача: {contract}")
                    video_send = await dp.bot.send_video(admin, video=video, caption="Відео відправлено",
                                                         reply_markup=back)
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                    await db.message("BOT", 10001, video_send.video.file_id, video_send.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення повідомлення: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_message_decline_phone":
            msg = await dp.bot.send_message(call.from_user.id, "Повідомлення не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_sms_decline", state="send_message_text")
@dp.callback_query_handler(IDFilter(ADMINS),
                           lambda call: call.data == "panel_send_sms_accept", state="send_message_text")
async def message_send_accept_sms(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    phone = await state.get_data()
    phone = phone['phone']
    data = await state.get_data()
    if data['type'] == 'text':
        text = data['text']
        await call.answer()
        if call.data == "panel_send_sms_accept":
            try:
                msg = await dp.bot.send_message(call.from_user.id, "Починаю відправку SMS...")
                await db.message("BOT", 10001, msg.html_text, msg.date)
                if len(text) >= 53:
                    print('transliterate')
                    email_text_t = transliterate.translit(text, 'uk', reversed=True)
                    text = f"{email_text_t}\nt.me/infoaura_bot"
                else:
                    print("no transliterate")
                    text = f"{text}\nt.me/infoaura_bot"
                async with aiohttp.ClientSession() as session:
                    print("Sending sms to:", phone,
                          "with text:", text,
                          "length of string:", len(text))
                    param = {'version': 'http',
                             'login': '380936425274',
                             "password": "iw79izvy",
                             'key': '6cf938587e0ed0d992566730169e82e229f097c7',
                             'command': "send",
                             'from': 'IAura',
                             'to': f"{phone}",
                             'message': f'{text}'}
                    async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                               params=param) as sms:
                        print("SMS: ", await sms.text())
                for admin in ADMINS:
                    msg_1 = await dp.bot.send_message(admin,
                                                      f"SMS відправлено"
                                                      f"\nТелефон: {phone}"
                                                      f"\nТекст: {text}",
                                                      reply_markup=back)
                    await db.message("BOT", 10001, msg_1.html_text, msg_1.date)
                await state.finish()
            except Exception as e:
                msg = await dp.bot.send_message(call.from_user.id, f"Помилка відправлення SMS: {e}",
                                                reply_markup=back)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.finish()
        elif call.data == "panel_send_sms_decline":
            msg = await dp.bot.send_message(call.from_user.id, "SMS не відправлено", reply_markup=back)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.finish()
    else:
        await call.answer()
        msg = await dp.bot.send_message(call.from_user.id, "SMS не відправлено, тільки текстові повідомлення",
                                        reply_markup=back)
        await db.message("BOT", 10001, msg.html_text, msg.date)
        await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), text='ban_account')
async def ban_account(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    await call.message.edit_text("Введіть Telegram ID користувача, якого потрібно заблокувати",
                                 reply_markup=back_inline)
    await state.set_state("ban_account")


@dp.message_handler(IDFilter(ADMINS), state="ban_account")
async def ban_id(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    try:
        await db.set_ban(message.text)
        await message.answer(f"Користувач з Telegram ID: {message.text} заблокований",
                             reply_markup=back)
    except Exception as e:
        await message.answer(f"Помилка: {e}", reply_markup=back)
    await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), text='unban_account')
async def unban_account(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    banned_id = await db.get_ban()
    if banned_id:
        await call.message.edit_text("Введіть Telegram ID користувача, якого потрібно розблокувати"
                                     "{} - {}".format(*banned_id),
                                     reply_markup=back_inline)
        await state.set_state("unban_account")
    else:
        await call.message.edit_text("Введіть Telegram ID користувача, якого потрібно розблокувати",
                                     reply_markup=back_inline)
        await state.set_state('unban_account')


@dp.message_handler(IDFilter(ADMINS), state="unban_account")
async def unban_id(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    try:
        await db.set_unban(message.text)
        await message.answer(f"Користувач з Telegram ID: {message.text} розблокований",
                             reply_markup=back)
    except Exception as e:
        await message.answer(f"Помилка: {e}", reply_markup=back)
    await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), text='register_alarm')
async def register_alarm(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    await call.message.answer("Виберіть группу аварії, або введіть номери договру котрих стосується аварія.\n"
                              "Наприклад:\n"
                              "10000001, 10000002, 10000003..", reply_markup=grp_choice)

    await state.set_state("register_alarm")


@dp.message_handler(IDFilter(ADMINS), state="register_alarm")
async def message_alarm_by_contract(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    try:
        alarm_message = "За Вашим підключенням зареєстрована аварійна ситуація. Перевірте термін усунення аварії пізніше."
        await db.insert_alarm(alarm_message, message.text)
        print(message.text)
        users = message.text.split(",")
        users_count = 0
        for user in users:
            try:
                if await db.set_alarm_for_users(user):
                    users_count += 1
            except Exception as e:
                print(e)
        await message.answer(f"Зареєстровано аварію для {users_count} користувачів", reply_markup=back)
    except Exception as e:
        await message.answer(f"Помилка: {e}", reply_markup=back_inline)
        await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), state="register_alarm")
async def message_alarm_grp(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    try:
        alarm_message = "За Вашим підключенням зареєстрована аварійна ситуація. Перевірте термін усунення аварії пізніше."
        await db.insert_alarm(alarm_message, call.data)
        print(call.data)
        users = await users_with_alarm(int(call.data))
        users_count = 0
        for user in users:
            try:
                if await db.set_alarm_for_users(user):
                    users_count += 1
            except Exception as e:
                print(e)
                pass
        await call.message.edit_text(f"Аварію успішно зареєстровано для {users_count} абонентів",
                                     reply_markup=back_inline)
        await state.finish()
    except Exception as e:
        await call.message.edit_text(f"Помилка: {e}", reply_markup=back_inline)
        await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), text='redactor_alarm')
async def redactor_alarm(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    await call.message.edit_text("Виберіть необхідну дію", reply_markup=redact_alarm)
    await state.set_state("redactor_alarm")


@dp.callback_query_handler(IDFilter(ADMINS), state='redactor_alarm')
async def redacting_alarm_state(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    await call.answer()
    records = await db.get_alarm()
    keyboard_alarms = InlineKeyboardMarkup(row_width=1)
    txt = ''
    print(records)
    if not records:
        await call.message.edit_text("Аварій не знайдено", reply_markup=back_inline)
        await state.finish()
    elif call.data == 'change_message_alarm':
        txt = 'Оберіть аварію для редагування за alarm_id\n'
        await state.set_state("redact_alarm")
    elif call.data == 'delete_alarm':
        txt = 'Оберіть аварію для видалення за alarm_id\n'
        await state.set_state("delete_alarm")
    if records:
        for record in records:
            txt += f"{record[0]} - {record[1]}\n"
            keyboard_alarms.add(InlineKeyboardButton(f"{record[0]}", callback_data=f"{record[0]}"))
        await call.message.edit_text(txt, reply_markup=keyboard_alarms)
    else:
        pass


@dp.callback_query_handler(IDFilter(ADMINS), state='redact_alarm')
async def redacting_alarm(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    try:
        await call.answer()
        await call.message.edit_text("Напишіть нове повідомлення що стосується аварії",
                                     reply_markup=back_inline)
        await state.set_state("redact_alarm_message")
        await state.update_data(alarm_id=call.data)
    except Exception as e:
        await call.message.edit_text(f"Помилка: {e}", reply_markup=back_inline)
        await state.finish()


@dp.callback_query_handler(IDFilter(ADMINS), state='delete_alarm')
async def deleting_alarm(call: types.CallbackQuery, state: FSMContext):
    await db.message(call.from_user.full_name, call.from_user.id, call.message.text, call.message.date)
    try:
        await call.answer()
        await db.delete_alarm(call.data)
        await call.message.edit_text(f"Аварію з ID {call.data} успішно видалено", reply_markup=back_inline)
        await state.finish()
    except Exception as e:
        await call.message.edit_text(f"Помилка: {e}", reply_markup=back_inline)
        await state.finish()


@dp.message_handler(IDFilter(ADMINS), state="redact_alarm_message")
async def message_alarm_redact(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    alarm_id = await state.get_data()
    await db.change_alarm_message(message.text, int(alarm_id['alarm_id']))
    await message.answer(f"Аварія з id: {alarm_id['alarm_id']}\nНовий текст аварії:\n{message.text}",
                         reply_markup=back)
    await state.finish()
