import imaplib
import email
import datetime
import transliterate

import aiohttp
from asyncio import sleep

from loader import dp, db
from utils.format_number import unformat_number, number
from aiogram.utils.exceptions import UserDeactivated, BotBlocked


async def send_message_sms(phone: int = None, text: str = None):
    if phone is None or text is None:
        # credentials
        username = "contact@infoaura.com.ua"

        # decode format
        decode_format = "utf-8"

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
                    print("Message: \n", message.decode(decode_format))
                    print("==========================================\n")
                    email_phone.append(email_message["subject"])
                    email_text.append(message.decode(decode_format))
                    break

        for i in range(len(email_phone)):
            await sleep(30)
            if email_phone[i] in users_phones:
                try:
                    divider = email_text[i].find('-----')
                    print(number(email_phone[i]))
                    telegram_id = await db.select_id_by_phone(phone_number=number(email_phone[i]))
                    telegram_id = telegram_id[0]['telegram_id']
                    print(telegram_id)
                    if divider:
                        await dp.bot.send_message(telegram_id, email_text[i][divider + 5:])
                        await db.message("BOT", 10001, email_text[i][divider + 5:], datetime.datetime.now())
                        print("Message sent via bot to:", telegram_id)
                    else:
                        await dp.bot.send_message(telegram_id, email_text[i])
                        await db.message("BOT", 10001, email_text[i], datetime.datetime.now())
                        print("Message sent via bot to:", telegram_id)
                except UserDeactivated as e:
                    print("Message sent via sms to:", email_phone[i])
                    try:
                        divider = email_text[i].find('-----')
                        if divider:
                            message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                        else:
                            message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                        print(len(message_sms_text), " length of  string")
                        if len(message_sms_text) >= 70:
                            print('transliterate')
                            if divider:
                                email_text_t = transliterate.translit(email_text[i][:divider], 'uk', reversed=True)
                            else:
                                email_text_t = transliterate.translit(email_text[i], 'uk', reversed=True)
                            message_sms_text = f"{email_text_t}\nt.me/infoaura_bot"
                        else:
                            print("no transliterate")
                            if divider:
                                message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                            else:
                                message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                        await sleep(60)
                        async with aiohttp.ClientSession() as session:
                            print("Sending sms to:", email_phone[i],
                                  "with text:", message_sms_text,
                                  "length of string:", len(message_sms_text))
                            param = {'version': 'http',
                                     'login': '380936425274',
                                     "pass": "iw79izvy",
                                     "key": '6cf938587e0ed0d992566730169e82e229f097c7',
                                     'command': "send",
                                     'from': 'IAura',
                                     'to': f"{email_phone[i]}",
                                     'message': f'{message_sms_text}'}
                            async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                                       params=param) as sms:
                                print("SMS: ", await sms.text())
                    except Exception as e:
                        print(e)
                    continue
                except BotBlocked:
                    print("Message sent via sms to:", email_phone[i])
                    try:
                        divider = email_text[i].find('-----')
                        if divider:
                            message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                        else:
                            message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                        print(len(message_sms_text), " length of  string")
                        if len(message_sms_text) >= 70:
                            print('transliterate')
                            if divider:
                                email_text_t = transliterate.translit(email_text[i][:divider], 'uk', reversed=True)
                            else:
                                email_text_t = transliterate.translit(email_text[i], 'uk', reversed=True)
                            message_sms_text = f"{email_text_t}\nt.me/infoaura_bot"
                        else:
                            print("no transliterate")
                            if divider:
                                message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                            else:
                                message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                        await sleep(60)
                        async with aiohttp.ClientSession() as session:
                            print("Sending sms to:", email_phone[i],
                                  "with text:", message_sms_text,
                                  "length of string:", len(message_sms_text))
                            param = {'version': 'http',
                                     'login': '380936425274',
                                     "password": "iw79izvy",
                                     'command': "send",
                                     'from': 'IAura',
                                     'to': f"{email_phone[i]}",
                                     'message': f'{message_sms_text}'}
                            async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                                       params=param) as sms:
                                print("SMS: ", await sms.text())
                    except Exception as e:
                        print(e)
                    continue
            else:
                print("Message sent via sms to:", email_phone[i])
                try:
                    divider = email_text[i].find('-----')
                    if divider:
                        message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                    else:
                        message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                    print(len(message_sms_text), " length of  string")
                    if len(message_sms_text) >= 70:
                        print('transliterate')
                        if divider:
                            email_text_t = transliterate.translit(email_text[i][:divider], 'uk', reversed=True)
                        else:
                            email_text_t = transliterate.translit(email_text[i], 'uk', reversed=True)
                        message_sms_text = f"{email_text_t}\nt.me/infoaura_bot"
                    else:
                        print("no transliterate")
                        if divider:
                            message_sms_text = f"{email_text[i][:divider]}\nt.me/infoaura_bot"
                        else:
                            message_sms_text = f"{email_text[i]}\nt.me/infoaura_bot"
                    await sleep(60)
                    async with aiohttp.ClientSession() as session:
                        print("Sending sms to:", email_phone[i],
                              "with text:", message_sms_text,
                              "length of string:", len(message_sms_text))
                        param = {'version': 'http',
                                 'login': '380936425274',
                                 "password": "iw79izvy",
                                 'command': "send",
                                 'from': 'IAura',
                                 'to': f"{email_phone[i]}",
                                 'message': f'{message_sms_text}'}
                        async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                                   params=param) as sms:
                            print("SMS: ", await sms.text())
                except Exception as e:
                    print(e)
                    pass
    else:
        users = await db.select_all_users()
        users_phones = []
        users_id = []
        for i in range(len(users)):
            users_phones.append(unformat_number(str(users[i]['phone_number'])))
            users_id.append(users[i]['telegram_id'])
            if int(phone) in users_id:
                await dp.bot.send_message(int(phone), text)
                await db.message("BOT", 10001, text, datetime.datetime.now())
                print("Message sent via bot to:", int(phone))
                result = f"Message sent via bot to: {int(phone)}"
                return result
            elif phone in users_phones:
                telegram_id = await db.select_id_by_phone(phone_number=number(phone))
                telegram_id = telegram_id[0]['telegram_id']
                msg = await dp.bot.send_message(telegram_id, text)
                await db.message("BOT", 10001, text, datetime.datetime.now())
                print("Message sent via bot to:", telegram_id)
                result = f"Message sent via bot to: {telegram_id}"
                return result
        else:
            print("Message sent via sms to:", phone)
            try:
                message_sms_text = f"{text}\nt.me/infoaura_bot"
                print(len(message_sms_text), " length of  string")
                if len(message_sms_text) >= 70:
                    print('transliterate')
                    email_text_t = transliterate.translit(text, 'uk', reversed=True)
                    message_sms_text = f"{email_text_t}\nt.me/infoaura_bot"
                else:
                    print("no transliterate")
                    message_sms_text = f"{text}\nt.me/infoaura_bot"
                await sleep(60)
                async with aiohttp.ClientSession() as session:
                    print("Sending sms to:", phone,
                          "with text:", message_sms_text,
                          "length of string:", len(message_sms_text))
                    param = {'version': 'http',
                             'login': '380936425274',
                             "password": "iw79izvy",
                             'command': "send",
                             'from': 'IAura',
                             'to': f"{phone}",
                             'message': f'{message_sms_text}'}
                    async with session.request('http', "https://smsukraine.com.ua/api/http.php",
                                               params=param) as sms:
                        print("SMS: ", await sms.text())
                        return "SMS: " + await sms.text()
            except Exception as e:
                print(e)
                return e

