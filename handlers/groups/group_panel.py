from aiogram import types
from aiogram.dispatcher.filters import BoundFilter, IDFilter

from data.config import ADMINS
from filters import IsGroup
from loader import dp


@dp.message_handler(IsGroup(), IDFilter(ADMINS), commands=["panel"])
async def menu(message: types.Message):
    await message.answer("Group panel")