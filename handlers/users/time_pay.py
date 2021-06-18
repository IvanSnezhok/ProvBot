from aiogram import types
from aiogram.dispatcher.filters import Text

from keyboards.default.buttons import return_button
from loader import dp, db
from middlewares import _, __
from utils.db_api import database


@dp.message_handler(Text(equals=__("Тимчасовий платіж")))
async def time_pay(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    if await database.t_pay(await db.select_contract(message.from_user.id)):
        msg = await message.answer(text=_("Доступ в Інтернет розблоковано на 24 години!\n"
    "Рахунок поповнено на {} грн на 24 години! Тепер можете повернутись у головне меню").format(database.time_pay[0]),
                                   reply_markup=return_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)
    else:
        msg = await message.answer(text=_("Ви не можете використати тимчасовий платіж!\n"
                                          "Користуватись тимчасовим платежем можна раз на місяць!"),
                                   reply_markup=return_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)
