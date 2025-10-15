import base64
import datetime
import transliterate
import os.path
import json
import pickle

import aiohttp
from asyncio import sleep
from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from google.oauth2 import service_account

from loader import dp, db
from utils.format_number import unformat_number, number
from aiogram.utils.exceptions import UserDeactivated, BotBlocked

# Gmail API scope - читання повідомлень
SCOPES = ['https://www.googleapis.com/auth/gmail.modify']


async def get_gmail_service():
    """Створення сервісу Gmail API використовуючи service account."""
    creds = None

    # Спробуємо завантажити збережені токени
    if os.path.exists('token.pickle'):
        with open('token.pickle', 'rb') as token:
            creds = pickle.load(token)

    # Якщо немає валідних токенів, використовуємо service account
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            # Використовуємо service account
            creds = service_account.Credentials.from_service_account_file(
                'service-account.json', scopes=SCOPES)
            # Додаємо делегування, щоб діяти від імені вказаного користувача
            delegated_creds = creds.with_subject('contact@infoaura.com.ua')
            creds = delegated_creds

        # Зберігаємо токени для наступного використання
        with open('token.pickle', 'wb') as token:
            pickle.dump(creds, token)

    return build('gmail', 'v1', credentials=creds)


async def send_message_sms(phone: int = None, text: str = None):
    if phone is None or text is None:
        # Створюємо Gmail API сервіс
        service = await get_gmail_service()

        # Отримуємо непрочитані повідомлення
        results = service.users().messages().list(
            userId='me',
            q='is:unread'
        ).execute()

        messages = results.get('messages', [])

        print("Total Messages Unseen:", len(messages) if messages else 0)

        # Отримуємо всіх користувачів з бази даних
        await db.create()
        users = await db.select_all_users()
        users_phones = []
        users_id = []
        for i in range(len(users)):
            users_phones.append(unformat_number(str(users[i]['phone_number'])))
            users_id.append(users[i]['telegram_id'])

        # Списки для зберігання даних електронної пошти
        email_phone = []
        email_text = []

        if not messages:
            print("No unread messages found.")
            return

        for message in messages:
            msg = service.users().messages().get(
                userId='me',
                id=message['id'],
                format='full'
            ).execute()

            # Позначаємо як прочитане
            service.users().messages().modify(
                userId='me',
                id=message['id'],
                body={'removeLabelIds': ['UNREAD']}
            ).execute()

            # Отримуємо заголовки
            headers = msg['payload']['headers']
            subject = ""
            sender = ""
            date = ""

            for header in headers:
                if header['name'] == 'Subject':
                    subject = header['value']
                if header['name'] == 'From':
                    sender = header['value']
                if header['name'] == 'Date':
                    date = header['value']

            print("\n===========================================")
            print("Subject: ", subject)
            print("From: ", sender)
            print("Date: ", date)

            # Отримуємо текст повідомлення
            message_text = ""

            if 'parts' not in msg['payload']:
                if 'data' in msg['payload'].get('body', {}):
                    data = msg['payload']['body']['data']
                    message_text = base64.urlsafe_b64decode(data).decode('utf-8')
            else:
                parts = msg['payload']['parts']
                for part in parts:
                    if part['mimeType'] == 'text/plain' or part['mimeType'] == 'text/html':
                        if 'data' in part.get('body', {}):
                            data = part['body']['data']
                            message_text = base64.urlsafe_b64decode(data).decode('utf-8')
                            break

            print("Message: \n", message_text)
            print("==========================================\n")

            email_phone.append(subject)
            email_text.append(message_text)

        # Решта вашого коду залишається майже ідентичною
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
                    # Код для відправки SMS залишається незмінним
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

