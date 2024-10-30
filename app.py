import logging
from datetime import datetime, date

from aiogram import executor, Dispatcher

from loader import dp, db, scheduler
import middlewares, filters, handlers
from utils.db_api import database
from utils.misc.debt_notification import schedule_debt_notification
from utils.misc.sms_message import send_message_sms
from utils.notify_admins import on_startup_notify, on_shutdown_notify


async def on_startup(dispatcher):
    # Уведомляет про запуск
    logging.info("Создаем подключение к локальной ДБ")
    await db.create()
    logging.info("Создаем таблицу пользователей")
    await db.create_table_users()
    logging.info("Создаем таблицу сообщений")
    await db.create_table_msg()
    logging.info("Создаем таблицу сигналов")
    await db.create_table_alarm()
    logging.info("Создаем таблицу оплат")
    await db.create_table_bill_check()
    logging.info("Створюємо таблицю чатів")
    await db.create_table_chats()
    logging.info("Создаем таблицу кликов")
    await db.create_table_user_clicks()
    logging.info("Включаем уведомления по таймеру")
    logging.info("Готово.")
    await on_startup_notify(dispatcher)
    scheduler.add_job(send_message_sms, 'interval', minutes=3)
    schedule_debt_notification(scheduler)
    scheduler.start()


if __name__ == '__main__':
    executor.start_polling(dp, on_startup=on_startup, on_shutdown=on_shutdown_notify)
