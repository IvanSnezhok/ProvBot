import logging
import random
import re
import uuid

import aiogram.utils.exceptions
from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Text
from aiogram.types import ContentType, LabeledPrice

from data.config import ADMINS
from data.pays_item import P180, P1080, P200, P1200, P350, P2100
from keyboards.default.buttons import return_button
from loader import dp, db, bot
from middlewares import _, __
from utils.db_api import database
from utils.db_api.database import check_contract_exists
from utils.misc.pay_load import Pay


@dp.message_handler(Text(__("Поповнити рахунок")), state="*")
async def contract_pay(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await database.search_query(tel=await db.select_tel(user_id=message.from_user.id))
    ban = await db.get_ban()
    print(ban)

    if message.from_user.id in ban:
        await message.answer(_("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки "
                               "support@infoaura.com.ua"))
    try:
        elif database.data[5] == 'СТАНДАРТ(180грн).':
            msg = await message.answer(text=_("Зверніть увагу, що тут ви можете поповнити тільки свій особовий рахунок!"))
            await db.message("BOT", 10001, msg.html_text, msg.date)
            invoice_pay = P180.generate_invoice()
            await state.update_data({'bill_id_p180': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=180)
            msg1 = await message.answer(
                text=_("Діє акція - поповни рахунок на 6 місяців одним платежем та отримуй 10% від суми поповнення!"))
            await db.message("BOT", 10001, msg1.html_text, msg1.date)
            invoice_pay = P1080.generate_invoice()
            await state.update_data({'bill_id_p1080': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=1080)
        elif database.data[5] == 'PON-100(200грн)' or database.data[5] == 'VIP WIFI-200':
            msg = await message.answer(text=_("Зверніть увагу, що тут ви можете поповнити тільки свій особовий рахунок!"))
            invoice_pay = P200.generate_invoice()
            await state.update_data({'bill_id_p200': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=200)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            msg1 = await message.answer(
                text=_("Діє акція - поповни рахунок на 6 місяців одним платежем та отримуй 10% від суми поповнення!"))
            invoice_pay = P1200.generate_invoice()
            await state.update_data({'bill_id_p1200': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=1200)
            await db.message("BOT", 10001, msg1.html_text, msg1.date)
        elif database.data[5] == 'PON-300(350грн)':
            msg = await message.answer(text=_("Зверніть увагу, що тут ви можете поповнити тільки свій особовий рахунок!"))
            invoice_pay = P350.generate_invoice()
            await state.update_data({'bill_id_p350': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=350)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            msg1 = await message.answer(
                text=_("Діє акція - поповни рахунок на 6 місяців одним платежем та отримуй 10% від суми поповнення!"))
            invoice_pay = P2100.generate_invoice()
            await state.update_data({'bill_id_p2100': invoice_pay[' start_parameter']})
            await bot.send_invoice(message.from_user.id, **invoice_pay, payload=2100)
            await db.message("BOT", 10001, msg1.html_text, msg1.date)
        else:
            msg = await message.answer(text=_("Для поповнення рахунку введіть сумму поповненя!\n"
                                            "Наприклад:\n"
                                            "250,\n"
                                            "500\n"),
                                    reply_markup=return_button)
            await state.set_state('invoice_payload')
            await db.message("BOT", 10001, msg.html_text, msg.date)
    except IndexError:
        msg = await message.answer(text=_("Для поповнення рахунку введіть сумму поповненя!\n"
                                            "Наприклад:\n"
                                            "250,\n"
                                            "500\n"),
                                    reply_markup=return_button)
            await state.set_state('invoice_payload')
            await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(state='invoice_payload')
async def get_invoice_payload(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    msg = await message.answer(text=_("Введіть номер договору який хочете поповнити!"))
    await state.set_data({'payload': message.text})
    await state.set_state('invoice_contract')
    await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(state='invoice_contract')
@dp.message_handler(state='invalid_payload')
async def get_invoice_contract(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    data = await state.get_data()
    if await state.get_state() == 'invoice_contract':
        data = await state.get_data()
        payload = data['payload']
        amount_pay = int(payload) * 100
        random_id = f"create_invoice_{amount_pay}_{uuid.uuid4()}"
        await state.update_data({'contract': message.text, 'bill_id': random_id})
        contract = message.text
        if re.match(r'^\d{8}$', contract) and await check_contract_exists(contract):
            invoice = Pay(
                title=f"Поповнення на {payload} грн",
                description=f"Поповнення рахунку {contract} на {payload} гривень",
                currency="UAH",
                prices=[
                    LabeledPrice(
                        label=f"Поповнення рахунку {contract}",
                        amount=amount_pay
                    )
                ],
                start_parameter=random_id
            )
            for admin in ADMINS:
                await dp.bot.send_message(admin,
                                          f"Створено інвойс для користувача {message.from_user.id}:\n"
                                          f"Договір: {contract}\n"
                                          f"Сума: {payload} грн\n"
                                          f"ID інвойсу: {random_id}")
            try:
                await bot.send_invoice(message.from_user.id, **invoice.generate_invoice(), payload=str(amount_pay))
            except aiogram.utils.exceptions.CurrencyTotalAmountInvalid:
                msg = await message.answer('Мінімальна сумма поповнення від 0.1$ в гривнях за курсом НБУ\n'
                                           'Введіть іншу сумму поповнення або можете повернутися у головне меню',
                                           reply_markup=return_button)
                await db.message("BOT", 10001, msg.html_text, msg.date)
                await state.set_state('invalid_payload')
        else:
            msg = await message.answer('Невірно вказаний номер договору для поповнення\n',
                                       reply_markup=return_button)
            await db.message("BOT", 10001, msg.html_text, msg.date)
    else:
        payload = message.text
        amount_pay = int(payload) * 100
        random_id = f'{data["bill_id"]}'
        contract = f'{data["contract"]}'
        invoice = Pay(
            title=f"Поповнення на {payload} грн",
            description=f"Поповнення рахунку {contract} на {payload} гривень",
            currency="UAH",
            prices=[
                LabeledPrice(
                    label=f"Поповнення рахунку {contract}",
                    amount=amount_pay
                )
            ],
            start_parameter=random_id
        )
        try:
            await bot.send_invoice(message.from_user.id, **invoice.generate_invoice(), payload=str(amount_pay))
        except aiogram.utils.exceptions.CurrencyTotalAmountInvalid:
            msg = await message.answer('Мінімальна сумма поповнення від 0.1$ в гривнях за курсом НБУ\n'
                                       'Введіть іншу сумму поповнення або можете повернутися у головне меню',
                                       reply_markup=return_button)
            await db.message("BOT", 10001, msg.html_text, msg.date)
            await state.set_state('invalid_payload')



@dp.pre_checkout_query_handler(state='*')
async def process_pre_checkout(query: types.PreCheckoutQuery):
    await bot.answer_pre_checkout_query(pre_checkout_query_id=query.id,
                                        ok=True)


@dp.message_handler(content_types=ContentType.SUCCESSFUL_PAYMENT, state="*")
async def process_successful_pay(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    data = await state.get_data()
    try:
        if data['contract']:
            data = await state.get_data()
            payload = data['payload']
            contract = data['contract']
            bill_id = data['bill_id']
            await db.add_bill(bill_id, message.from_user.id, message.date, message.from_user.username, contract,
                              payload)
            await database.pay_balance(contract=contract, payload=payload)
            for admin in ADMINS:
                try:
                    msg = await dp.bot.send_message(chat_id=admin,
                                                    text=_("Користувач {} успішно поповнив рахунок "
                                                           "на {} {}").format(
                                                        contract, payload, message.successful_payment.currency)
                                                    )
                    await db.message("BOT", 10001, msg.html_text, msg.date)
                except Exception as err:
                    logging.exception(err)
        else:
            contract = await db.select_contract(message.from_user.id)
            contract = contract[0]
            payload = message.successful_payment.total_amount // 100
        if payload == 180:
            bill_id = data['bill_id_p180']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        elif payload == 1080:
            bill_id = data['bill_id_p1080']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        elif payload == 200:
            bill_id = data['bill_id_p200']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        elif payload == 1200:
            bill_id = data['bill_id_p1200']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        elif payload == 350:
            bill_id = data['bill_id_p350']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        elif payload == 2100:
            bill_id = data['bill_id_p2100']
            await db.add_bill(
                bill_id, message.from_user.id, message.date, message.from_user.username, contract, payload)
        await database.pay_balance(contract=contract[0], payload=payload)
        msg = await dp.bot.send_message(chat_id=message.from_user.id,
                                        text=__("Ваш рахунок поповнено на {} {}!").format(
                                            payload, message.successful_payment.currency),
                                        reply_markup=return_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)
    except Exception as e:
        logging.error(f"Помилка при обробці платежу: {e}")
        error_msg = await dp.bot.send_message(chat_id=message.from_user.id,
                                              text=__(
                                                  "На жаль, виникла помилка при обробці платежу. Будь ласка, зверніться до служби підтримки."),
                                              reply_markup=return_button)
        await db.message("BOT", 10001, error_msg.html_text, error_msg.date)
        for admin in ADMINS:
            await dp.bot.send_message(chat_id=admin,
                                      text=f"Помилка при обробці платежу для користувача {message.from_user.id}: {e}")