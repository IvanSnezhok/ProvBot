from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ReplyKeyboardMarkup, KeyboardButton

admin_keyboard = InlineKeyboardMarkup(row_width=2)
admin_keyboard.add(
    InlineKeyboardButton(text="Відправити повідомлення", callback_data="panel_send_message"),
    InlineKeyboardButton(text="Меню абонента", callback_data="account_menu"),
    InlineKeyboardButton(text="Заблокувати юзера", callback_data="ban_account"),
    InlineKeyboardButton(text="Розблокувати юзера", callback_data="unban_account"),
    InlineKeyboardButton(text="Зареэструвати аварію", callback_data="register_alarm"),
    InlineKeyboardButton(text="Редактор аварій", callback_data="redactor_alarm"))

accept_message = InlineKeyboardMarkup()
accept_message.add(
    InlineKeyboardButton(text="Так", callback_data="panel_send_message_accept"),
    InlineKeyboardButton(text="Ні", callback_data="panel_send_message_decline"))

accept_message_phone = InlineKeyboardMarkup()
accept_message_phone.add(
    InlineKeyboardButton(text="Так", callback_data="panel_send_message_accept_phone"),
    InlineKeyboardButton(text="Ні", callback_data="panel_send_message_decline_phone"))

accept_sms = InlineKeyboardMarkup()
accept_sms.add(
    InlineKeyboardButton(text="Так", callback_data="panel_send_sms_accept"),
    InlineKeyboardButton(text="Ні", callback_data="panel_send_sms_decline"))

back = ReplyKeyboardMarkup(resize_keyboard=True)
back.add(KeyboardButton(text="Назад"))

back_inline = InlineKeyboardMarkup()
back_inline.add(InlineKeyboardButton(text="Назад", callback_data="back"))

admin_account_menu = InlineKeyboardMarkup()
admin_account_menu.add(
    InlineKeyboardButton(text='Змінити баланс', callback_data='admin_change_balance'),
    InlineKeyboardButton(text='Тимчасовий платіж', callback_data='admin_temporary_payment'),
    InlineKeyboardButton(text='Змінити пакет', callback_data='admin_change_paket'),
    InlineKeyboardButton(text='Історія повідомлень', callback_data='message_history'))

search_choice = InlineKeyboardMarkup(row_width=2)
search_choice.add(InlineKeyboardButton(text='Номер договору', callback_data='search_contract'),
                  InlineKeyboardButton(text='Номер телефону', callback_data='search_phone'),
                  InlineKeyboardButton(text='Ім\'я', callback_data='search_name'),
                  InlineKeyboardButton(text='Адреса', callback_data='search_address'))

grp_choice = InlineKeyboardMarkup(row_width=2)
grp_choice.add(
    InlineKeyboardButton(text='"Удаленные"', callback_data='1'),
    InlineKeyboardButton(text='"Администраторы"', callback_data='2'),
    InlineKeyboardButton(text='"Сервера"', callback_data='3'),
    InlineKeyboardButton(text='"VIP"', callback_data='4'),
    InlineKeyboardButton(text='"Куренівка WhiteIPs"', callback_data='5'),
    InlineKeyboardButton(text='"Остановка"', callback_data='6'),
    InlineKeyboardButton(text='"Приостановка"', callback_data='7'),
    InlineKeyboardButton(text='"Савенка PON"', callback_data='8'),
    InlineKeyboardButton(text='"СМАРТ ИПТВ"', callback_data='9'),
    InlineKeyboardButton(text='"WiFi Киев"', callback_data='10'),
    InlineKeyboardButton(text='"WiFi Савенка"', callback_data='11'),
    InlineKeyboardButton(text='"Абоненты Киев"', callback_data='12'),
    InlineKeyboardButton(text='"Куренівка GreyIPs"', callback_data='13'),
    InlineKeyboardButton(text='"Поділ GreyIPs"', callback_data='14'),
    InlineKeyboardButton(text='"Катюжанка GreyIPs"', callback_data='15'))

redact_alarm = InlineKeyboardMarkup(row_width=2)
redact_alarm.add(
    InlineKeyboardButton(text='Видалити аварії', callback_data='delete_alarm'),
    InlineKeyboardButton(text='Змінити повідомлення', callback_data='change_message_alarm'))

