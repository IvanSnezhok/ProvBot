from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton

from keyboards.inline.callback_datas import start_callback

choice_lang = InlineKeyboardMarkup(row_width=3,
                                   inline_keyboard=[
                                       [
                                           InlineKeyboardButton(
                                               text="🇺🇦 UA",
                                               callback_data=start_callback.new("UA")
                                           ),
                                           InlineKeyboardButton(
                                               text="🇺🇸 EN",
                                               callback_data=start_callback.new("EN")
                                           ),
                                           InlineKeyboardButton(
                                               text="🇷🇺 RU",
                                               callback_data=start_callback.new("RU")
                                           )
                                       ]
                                   ])