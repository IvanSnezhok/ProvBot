import logging

from aiogram import executor

from loader import dp, db
import middlewares, filters, handlers
from utils.notify_admins import on_startup_notify, on_shutdown_notify


async def on_startup(dispatcher):
    # Уведомляет про запуск
    logging.info("Создаем подключение к локальной ДБ")
    await db.create()
    logging.info("Создаем таблицу пользователей")
    await db.create_table_users()
    logging.info("Готово.")
    await on_startup_notify(dispatcher)


if __name__ == '__main__':
    executor.start_polling(dp, on_startup=on_startup, on_shutdown=on_shutdown_notify)
