import logging
from datetime import datetime, date

from aiogram import executor, Dispatcher

from loader import dp, db, scheduler
import middlewares, filters, handlers
from utils.db_api import database
from utils.notify_admins import on_startup_notify, on_shutdown_notify


async def notify_clients(dp: Dispatcher):
    client_id = await db.contract()
    for i in client_id:
        telegram_id = 390616685
        contract = 10002131
        today = date.today()
        today = today.replace(day=12)
        today = today.strftime("%d.%m.%y")
        balance = await database.balance(contract)
        if balance is False:
            pass
        else:
            await dp.bot.send_message(telegram_id,
                                      middlewares._("Шановний клієнт! Доступ до інтернету за рахунком {} буде "
                                                    "заблокований "
                                        "{}"
                                        "Рекомендуємо поповнити баланс мінімум на {}").format(contract, today, balance))


def scheduler_jobs():
    scheduler.add_job(notify_clients, "cron", day="9", hour="22", minute="43", args=(dp,))


async def on_startup(dispatcher):
    # Уведомляет про запуск
    logging.info("Создаем подключение к локальной ДБ")
    await db.create()
    logging.info("Создаем таблицу пользователей")
    await db.create_table_users()
    logging.info("Создаем таблицу сообщений")
    await db.create_table_msg()
    logging.info("Включаем уведомления по таймеру")
    scheduler_jobs()
    logging.info("Готово.")
    await on_startup_notify(dispatcher)


if __name__ == '__main__':
    scheduler.start()
    executor.start_polling(dp, on_startup=on_startup, on_shutdown=on_shutdown_notify)
