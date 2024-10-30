from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.dispatcher.filters import Command
from aiogram.dispatcher.filters.state import StatesGroup, State

from keyboards.default.buttons import return_button
from loader import dp, db
from middlewares import _


class Contract(StatesGroup):
    get_id = State()
    text = State()


@dp.message_handler(Command('contact'))
async def get_id(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    if message.from_user.id == 390616685 or message.from_user.id == 133347159:
        client_id = await db.choose_contract()
        msg = await message.answer(text=_("Оберіть клієнта, написав його telegram_id\n"))
        await Contract.get_id.set()
        # for i in client_id:
        #     await message.answer(f"Імя: {i[0]}\n"
        #                          f"Телеграм айді: {i[1]}\n"
        #                          f"Контракт: {i[2]}")
        await db.message("BOT", 10001, msg.html_text, msg.date)
    else:
        msg = await message.answer(_("Эту команду могут использовать только администраторы"),
                                   reply_markup=return_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)


@dp.message_handler(state=Contract.get_id)
async def text(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    await state.update_data(telegram_id=message.text)
    msg = await message.answer(_("Тепер напишіть повідомлення для клієнта"))
    await db.message("BOT", 10001, msg.html_text, msg.date)
    await Contract.text.set()


@dp.message_handler(state=Contract.text)
async def contact(message: types.Message, state: FSMContext):
    await db.message(message.from_user.full_name, message.from_user.id, message.text, message.date)
    data = await state.get_data()
    try:
        msg = await dp.bot.send_message(data.get("telegram_id"), message.text)
        await db.message("BOT", 10001, msg.html_text, msg.date)
        await state.reset_state()
        msg1 = await message.answer("Повідомлення відправлено", reply_markup=return_button)
        await db.message("BOT", 10001, msg1.html_text, msg.date)
    except Exception as e:
        msg = await message.answer(_("Повідомлення не було доставлене, скоріш за все telegram id не правильний"),
                                   reply_markup=return_button)
        await db.message("BOT", 10001, msg.html_text, msg.date)
        await state.reset_state()
