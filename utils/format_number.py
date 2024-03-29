import time

from loader import db
from utils.db_api import database


def format_number(telegram_number: str):
    telegram_number = "".join(num for num in telegram_number if num not in " +()")[2:]
    r = telegram_number[:3]
    r1 = telegram_number[3:6]
    r2 = telegram_number[6:]
    r2 = r2[:2]
    r3 = telegram_number[8:]
    result = "{}-{}-{}-{}".format(r, r1, r2, r3)
    return result


def unformat_number(telegram_number: str):
    telegram_number = "".join(num for num in telegram_number)
    num1 = telegram_number[:3]
    num2 = telegram_number[4:7]
    num3 = telegram_number[8:10]
    num4 = telegram_number[11:]
    return "{}{}{}{}".format(num1, num2, num3, num4)


def number(phone_number: str):
    num1 = phone_number[:3]
    num2 = phone_number[3:6]
    num3 = phone_number[6:8]
    num4 = phone_number[8:]
    return "{}-{}-{}-{}".format(num1, num2, num3, num4)


async def format_text_account(text: dict) -> str:
    if text['state'] == 'off':
        text['state'] = 'Заблоковано'
    elif text['state'] == 'pause':
        text['state'] = 'Пауза'
    else:
        text['state'] = 'Активна'
    result = f"""Ваш username: {text['name']}\nНа вашому рахунку: {text['balance']}\nВаш номер договору: {text['contract']}\nВаше ПІБ: {text['fio']}\nСтан послуги: {text['state']}\nВаш пакет: {text['paket']}\nВаша Адреса: {text['address']}"""
    return result


async def format_text_account_admin(text: dict) -> str:
    phone_bot = await db.get_phone_by_contract(text["contract"])
    if text['state'] == 'off':
        text['state'] = 'Заблоковано'
    elif text['state'] == 'pause':
        text['state'] = 'Пауза'
    else:
        text['state'] = 'Активна'
    result = f"""Ваш username: {text['name']}\nНа вашому рахунку: {text['balance']}\nВаш номер договору: {text['contract']}\nВаш ip: {text['ip']}\nВаше ПІБ: {text['fio']}\nСтан послуги: {text['state']}\nВаш пакет: {text['paket']}\nВаша Адреса: {text['address']}\nВаш телефон: {text['telefon']}\nТелефон що в боті: {phone_bot}"""
    return result
