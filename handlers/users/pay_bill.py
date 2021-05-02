import logging

from aiogram import types
from aiogram.dispatcher.filters import Text
from aiogram.types import ContentType

from data.config import ADMINS
from data.pays_item import P150, P900, P200, P1200
from keyboards.default.buttons import return_button
from loader import dp, db, bot
from utils.db_api import database
from middlewares import _, __


@dp.message_handler(Text(__("Поповнити рахунок")))
async def contract_pay(message: types.Message):
    await database.search_query(tel=await db.select_tel(user_id=message.from_user.id))
    if database.data[5] == '150.':
        await message.answer(text=_("Зверніть увагу, що тут ви можете поповнити тільки свій особистий рахунок!"))
        await bot.send_invoice(message.from_user.id, **P150.generate_invoice(), payload=150)
        await message.answer(
            text=_("Діє акція - поповни рахунок на 6 місяців уперед та отримуй 10 % від суми поповнення!"))
        await bot.send_invoice(message.from_user.id, **P900.generate_invoice(), payload=900)
    elif database.data[5] == '200':
        await message.answer(text=_("Зверніть увагу, що тут ви можете поповнити тільки свій особистий рахунок!"))
        await bot.send_invoice(message.from_user.id, **P200.generate_invoice(),
                               payload=200)
        await message.answer(
            text=_("Діє акція - поповни рахунок на 6 місяців уперед та отримуй 10 % від суми поповнення!"))
        await bot.send_invoice(message.from_user.id, **P1200.generate_invoice(),
                               payload=1200)
    else:
        await message.answer(text=_("Вибачте, але для вашого тарифу не передбачено поповнення рахунку через бот"),
                             reply_markup=return_button)


@dp.pre_checkout_query_handler()
async def process_pre_checkout(query: types.PreCheckoutQuery):
    await bot.answer_pre_checkout_query(pre_checkout_query_id=query.id,
                                        ok=True)


@dp.message_handler(content_types=ContentType.SUCCESSFUL_PAYMENT)
async def process_successful_pay(message: types.Message):
    contract = await db.select_contract(message.from_user.id)
    contract = contract[0]
    payload = message.successful_payment.total_amount // 100
    await database.pay_balance(contract=contract[0], payload=payload)
    for admin in ADMINS:
        try:
            await dp.bot.send_message(chat_id=admin,
                                      text=_("Користувач {} поповнив рахунок "
                                             "на {} {}").format(
                                          contract[0], payload, message.successful_payment.currency)
                                      )
        except Exception as err:
            logging.exception(err)
    await dp.bot.send_message(chat_id=message.from_user.id,
                              text=__("Ваш рахунок поповнено на {} {}!").format(
                                payload, message.successful_payment.currency),
                              reply_markup=return_button)
