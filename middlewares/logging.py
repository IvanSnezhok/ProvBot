from aiogram import types
from aiogram.dispatcher.middlewares import BaseMiddleware
from datetime import datetime
import pytz

from loader import db

class LoggingMiddleware(BaseMiddleware):
    async def on_pre_process_message(self, message: types.Message, data: dict):
        # Отримуємо час за Києвом
        kyiv_tz = pytz.timezone('Europe/Kiev')
        current_time = datetime.now(kyiv_tz)
        
        # Записуємо повідомлення від користувача
        if message.text:
            await db.message(
                full_name=message.from_user.full_name,
                telegram_id=message.from_user.id,
                text=message.text,
                date=current_time
            )

    async def on_post_process_message(self, message: types.Message, results, data: dict):
        # Отримуємо час за Києвом для відповіді бота
        kyiv_tz = pytz.timezone('Europe/Kiev')
        current_time = datetime.now(kyiv_tz)
        
        # Записуємо відповідь бота
        if results and isinstance(results, types.Message):
            await db.message(
                full_name="BOT",
                telegram_id=10001,
                text=results.html_text if hasattr(results, 'html_text') else results.text,
                date=current_time
            )
