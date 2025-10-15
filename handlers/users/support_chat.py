from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Text

from loader import dp, db, bot
from states.get_client import SupportChat
from data.config import ADMINS
from keyboards.default.buttons import return_button
from middlewares import _

active_chats = {}

@dp.message_handler(Text(equals=__("Чат з тех. підтримкою")))
async def start_support_chat(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    
    # Перевірка, чи користувач вже в чаті
    if message.from_user.id in active_chats:
        await message.answer(_("Ви вже в активному чаті. Очікуйте відповіді від менеджера."))
        return

    await SupportChat.WaitingForSupport.set()
    await state.update_data(user_id=message.from_user.id)
    
    await message.answer(_("Ви увійшли в чат підтримки. Очікуйте підключення менеджера."), reply_markup=return_button)
    
    # Повідомлення для адміністраторів
    for admin in ADMINS:
        await bot.send_message(admin, f"Новий запит на чат підтримки від користувача {message.from_user.full_name} (ID: {message.from_user.id})")

@dp.message_handler(state=SupportChat.WaitingForSupport)
async def support_chat_waiting(message: types.Message, state: FSMContext):
    await message.answer(_("Будь ласка, очікуйте підключення менеджера."))

@dp.message_handler(commands=['connect'], user_id=ADMINS)
async def connect_to_chat(message: types.Message, state: FSMContext):
    args = message.get_args()
    if not args:
        await message.answer("Використання: /connect <user_id>")
        return

    user_id = int(args)
    if user_id not in active_chats:
        active_chats[user_id] = message.from_user.id
        await bot.send_message(user_id, _("Менеджер підключився до чату. Можете писати ваші повідомлення."))
        await message.answer(f"Ви підключилися до чату з користувачем {user_id}")
        
        # Повідомлення для інших адміністраторів
        for admin in ADMINS:
            if admin != message.from_user.id:
                await bot.send_message(admin, f"Адміністратор {message.from_user.full_name} підключився до чату з користувачем {user_id}")
    else:
        await message.answer("Цей користувач вже в активному чаті з іншим менеджером.")

@dp.message_handler(state=SupportChat.Chatting)
async def handle_support_message(message: types.Message, state: FSMContext):
    data = await state.get_data()
    user_id = data.get('user_id')
    admin_id = active_chats.get(user_id)

    if admin_id:
        await bot.send_message(admin_id, f"Повідомлення від користувача: {message.text}")
    else:
        await message.answer(_("На жаль, зв'язок з менеджером втрачено. Спробуйте почати чат знову."))
        await state.finish()

@dp.message_handler(user_id=ADMINS, state='*')
async def handle_admin_message(message: types.Message, state: FSMContext):
    for user_id, admin_id in active_chats.items():
        if admin_id == message.from_user.id:
            await bot.send_message(user_id, f"Відповідь від менеджера: {message.text}")
            return
    
    await message.answer("Ви не підключені до жодного активного чату.")

@dp.message_handler(commands=['end_chat'], user_id=ADMINS)
async def end_chat(message: types.Message, state: FSMContext):
    for user_id, admin_id in active_chats.items():
        if admin_id == message.from_user.id:
            del active_chats[user_id]
            await bot.send_message(user_id, _("Чат завершено. Дякуємо за звернення!"))
            await message.answer(f"Чат з користувачем {user_id} завершено.")
            return
    
    await message.answer("Ви не підключені до жодного активного чату.")