import logging

from aiogram import Dispatcher

from data.config import ADMINS

logger = logging.getLogger(__name__)


async def on_startup_notify(dp: Dispatcher):
    for admin in ADMINS:
        try:
            await dp.bot.send_message(admin, "Бот Запущен, попробуй /start")

        except Exception as err:
            logger.exception(err)


async def on_shutdown_notify(dp: Dispatcher):
    for admin in ADMINS:
        try:
            await dp.bot.send_message(admin, "Бот выключен, взаимодействие невозможно")

        except Exception as err:
            logger.exception(err)
