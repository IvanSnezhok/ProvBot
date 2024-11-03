from aiogram import types
from aiogram.dispatcher.filters import Text

from data.config import ADMINS
from keyboards.default.buttons import return_button
from loader import dp, db
from middlewares import _, __
from utils.db_api import database


@dp.message_handler(Text(equals=__("Тимчасовий платіж")))
async def time_pay(message: types.Message):
    
    user = await db.select_contract(message.from_user.id)
    user = user[0]["contract"]
    ban = await db.get_ban()
    if message.from_user.id in ban:
        await message.answer(
            _("Вітаємо! Для звернення, будь-ласка, скористайтесь нашим email технічної підтримки support@infoaura.com.ua"))
    elif await database.t_pay(user):
        msg = await message.answer(text=_("Доступ в Інтернет розблоковано на 24 години!\n"
    "Рахунок поповнено на {} грн на 24 години! Тепер можете повернутись у головне меню").format(database.time_pay_b[0]),
                                   reply_markup=return_button)
        
        for admin in ADMINS:
            await dp.bot.send_message(admin, _("Користувач {} використав тимчасовий платіж!").format(user))
    else:
        msg = await message.answer(text=_("Ви не можете використати тимчасовий платіж!\n"
                                          "Користуватись тимчасовим платежем можна раз на місяць!"),
                                   reply_markup=return_button)
        
