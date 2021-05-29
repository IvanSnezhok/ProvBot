from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Command

from loader import dp, db

from middlewares import _, __


# Эхо хендлер, куда летят текстовые сообщения без указанного состояния
@dp.message_handler(content_types=types.ContentTypes.ANY, state=None)
@dp.message_handler(Command(['help']))
async def bot_echo(message: types.Message):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    msg = await message.answer(_("Для взаємодії з ботом вам потрібно натиснути кнопку"))
    await db.message("BOT", 10001, msg.html_text, msg.date)


# Эхо хендлер, куда летят ВСЕ сообщения с указанным состоянием
# @dp.message_handler(state="*", content_types=types.ContentTypes.ANY)
# async def bot_echo_all(message: types.Message, state: FSMContext):
#     await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
#     state = await state.get_state()
#     await message.answer(f"Эхо в состоянии <code>{state}</code>.\n"
#                          f"\nСодержание сообщения:\n"
#                          f"<code>{message}</code>")
