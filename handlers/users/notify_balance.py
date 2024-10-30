import datetime

from data import config
from loader import scheduler
from aiogram import Dispatcher
from loader import db
from utils.db_api import database
from middlewares import _


async def notify_clients(dp: Dispatcher):
    client_id = await db.contract()
    for i in client_id:
        telegram_id = i[0]
        contract = i[1]
        today = datetime.date.today()
        today = today.replace(day=12)
        today = today.strftime("%d.%m.%y")
        balance = await database.balance(contract)
        if balance is False:
            pass
        else:
            await dp.bot.send_message(telegram_id,
                                      _("Шановний клієнт! Доступ до інтернету за рахунком {} буде заблокований "
                                        "{}"
                                        "Рекомендуємо поповнити баланс мінімум на {}").format(contract, today, balance))


def scheduler_jobs():
    scheduler.add_job(notify_clients, "cron", day="9", hour="23", minute="45")
