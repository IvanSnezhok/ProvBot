import asyncio
import imaplib
import email
import datetime

import aiohttp
from aiogram import types
from loader import dp, db
from utils.format_number import unformat_number, format_number, number


async def send_message(message: types.Message = None):
    # credentials
    username = "contact@infoaura.com.ua"

    # generated app password
    app_password = "^W7VCV<:'kKhBb9N"

    # https://www.systoolsgroup.com/imap/
    gmail_host = 'imap.gmail.com'

    # set connection
    mail = imaplib.IMAP4_SSL(gmail_host)

    # login
    mail.login(username, app_password)

    # select inbox
    mail.select("INBOX")

    # select specific mails
    _, selected_mails = mail.search(None, 'UNSEEN')

    # total number of mails from specific user
    print("Total Messages Unseen:", len(selected_mails[0].split()))

    # get all users from db
    await db.create()
    users = await db.select_all_users()
    users_phones = []
    users_id = []
    for i in range(len(users)):
        users_phones.append(unformat_number(str(users[i]['phone_number'])))
        users_id.append(users[i]['telegram_id'])

    # email to dict_list
    email_phone = []
    email_text = []
    for num in selected_mails[0].split():
        _, data = mail.fetch(num, '(RFC822)')
        _, bytes_data = data[0]

        # convert the byte data to message
        email_message = email.message_from_bytes(bytes_data)
        print("\n===========================================")

        # access data
        print("Subject: ", email_message["subject"])
        print("To:", email_message["to"])
        print("From: ", email_message["from"])
        print("Date: ", email_message["date"])
        for part in email_message.walk():
            if part.get_content_type() == "text/plain" or part.get_content_type() == "text/html":
                message = part.get_payload(decode=True)
                print("Message: \n", message.decode())
                print("==========================================\n")
                email_phone.append(email_message["subject"])
                email_text.append(message.decode())
                break

    for i in range(len(email_phone)):
        if email_phone[i] in users_phones:
            print(number(email_phone[i]))
            telegram_id = await db.select_id_by_phone(phone_number=number(email_phone[i]))
            print(telegram_id)
            await dp.bot.send_message(telegram_id, email_text[i])
            await db.message("BOT", 10001, email_text[i], datetime.datetime.now())
            print("Message sent via bot to:", telegram_id)
        else:
            print("Message sent via sms to:", email_phone[i])
            try:
                async with aiohttp.ClientSession() as session:
                    param = {'version': 'http',
                             'login': '380936425274',
                             "password": "iw79izvy",
                             'command': "send",
                             'from': 'IAura',
                             'to': f"{email_phone[i]}",
                             'message': f'{email_text[i]}\nt.me/infoaura_bot'}
                    async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                               params=param) as sms_get:
                        pass
            except Exception as e:
                print(e)
                pass

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(send_message())
