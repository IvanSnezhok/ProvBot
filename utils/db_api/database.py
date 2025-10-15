import asyncio
import logging
import time
from datetime import datetime
from typing import Dict

import aiomysql
import pytz as pytz

from data import config

loop = asyncio.get_event_loop()
data = []
plan = []
time_pay_b = []


async def search_query(tel):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, charset="cp1251")
    cur = await conn.cursor()
    await cur.execute("SELECT name, balance, contract, fio, state, paket, grp "
                      "FROM `users` "
                      f"WHERE telefon LIKE '{tel}%' ")
    result = await cur.fetchall()
    data.clear()
    plan.clear()
    try:
        result = result[0]
        plan.append(result[5])
        await cur.execute('SELECT name FROM `plans2` WHERE id=%s', result[5])
        paket = await cur.fetchall()
        paket = paket[0]
        data.append(result[0])
        data.append(result[1])
        data.append(result[2])
        data.append(result[3])
        data.append(result[4])
        data.append(paket[0])
        data.append(result[6])

    except IndexError:
        result = None
    await cur.close()
    conn.close()


async def account_show(account) -> Dict:
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, charset="cp1251")
    cur = await conn.cursor()
    await cur.execute("SELECT name, balance, contract, fio, state, paket "
                      "FROM `users` "
                      f"WHERE contract LIKE '{account}%' ")
    result = await cur.fetchall()
    result = result[0]
    result_dict = {'name': result[0],
                   'balance': result[1],
                   'contract': result[2],
                   'fio': result[3],
                   'state': result[4],
                   'paket': result[5]}
    await cur.execute('SELECT name FROM `plans2` WHERE id=%s', result_dict['paket'])
    paket = await cur.fetchall()
    paket = paket[0]
    result_dict['paket'] = paket[0]
    await cur.close()
    conn.close()
    return result_dict


async def balance(contract):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute("SELECT balance, paket FROM users WHERE contract=%s", contract)
    client_balance = await cur.fetchall()
    client_balance = client_balance[0]
    await cur.execute('SELECT price FROM `plans2` WHERE id=%s', client_balance[1])
    tariff = await cur.fetchall()
    tariff = tariff[0][0]
    inequality = tariff[0] - client_balance[0]
    if inequality >= tariff:
        return False
    else:
        return inequality


async def pay_balance_150(contract):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(f"UPDATE users set balance = balance + 150 WHERE contract={contract}")
    await cur.close()
    conn.close()


async def balance_change(contract, cash):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(f"UPDATE users set balance = {cash} WHERE contract={contract}")
    await cur.close()
    conn.close()


async def pay_balance(contract, payload):
    now_time = datetime.now(pytz.timezone('Europe/Kiev'))
    now_t = time.mktime(now_time.timetuple()) + now_time.microsecond / 1E6
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(
        f"SELECT paket, id FROM users WHERE contract={contract}")
    user = await cur.fetchall()
    try:
        user = user[0]
        paket = user[0]
        id = user[1]
    except IndexError:
        user = None
        logging.info("User not Found")
        return False
    if paket:
        await cur.execute(f"SELECT price FROM plans2 WHERE id = {paket}")
    price = await cur.fetchall()
    try:
        price = price[0]
    except IndexError:
        price = None
        logging.info("Price not find")
        return False
    await cur.execute(f"UPDATE users set balance = balance + {payload} WHERE contract={contract}")
    await cur.execute("INSERT INTO pays (mid,cash,time,admin,reason,coment) VALUES (%s, %s, %s, %s, %s, %s)",
                      (
                          (id),
                          (payload),
                          (now_t),
                          ('BOT'),
                          (str(now_t)),
                          ('Popolnenie via BOT')
                      ))
    await cur.execute(f"SELECT balance FROM users WHERE contract = {contract}")
    balance_n = await cur.fetchall()
    balance_n = balance_n[0][0]
    if balance_n >= price[0]:
        await cur.execute(f'UPDATE users set state="on" WHERE contract={contract}')
    await cur.close()
    conn.close()


async def t_pay(contract):  # Временный плтажеж
    now_time = datetime.now(pytz.timezone('Europe/Kiev'))
    now_t = time.mktime(now_time.timetuple()) + now_time.microsecond / 1E6
    next_t = now_t + 86400

    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, charset='cp1251')
    cur = await conn.cursor()
    await cur.execute(
        f"SELECT t_pay, paket, srvs, contract, fio, telefon, start_day, balance, id FROM users WHERE contract={contract}")
    user = await cur.fetchall()
    try:
        user = user[0]
        time_pay = user[0]
        paket = user[1]
        my_srvs = user[2]
        my_contract = user[3]
        my_fio = user[4]
        my_telefon = user[5]
        my_start_day = user[6]
        old_balance = user[7]
        id = user[8]
    except IndexError:
        user = None
        logging.info("User not Found")
        return False

    if time_pay == 0:
        if paket:
            await cur.execute(f"SELECT price FROM plans2 WHERE id = {paket}")
        price = await cur.fetchall()
        try:
            price = int(price[0][0])
            print(price)
        except IndexError:
            price = None
            logging.info("Price not find")
            return False

        if old_balance > 0:
            balance = old_balance + price
        else:
            balance = -old_balance + price
        await cur.execute(
            f"INSERT INTO pays (mid, cash, time, bonus, admin, reason,coment,flag) VALUES ({id}, {price}, '{next_t}', 'y','BOT','Platej sozdan {now_t}','Razblokirovan na 24 chasa', 't');")
        await cur.execute(f"UPDATE users SET balance={balance} WHERE contract={contract}")
        await cur.execute(f"UPDATE users SET state='on' WHERE contract={contract}")
        await cur.execute(f"UPDATE users SET t_pay=1 WHERE contract={contract}")
        await cur.commit()
        time_pay_b.clear()
        time_pay_b.append(balance)
        return True
    else:
        return False


async def check_net_pause(contract):
    # IF TRUE, NETPAUSE IS OFF
    # IF FALSE, NETPAUSE IS ON
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(f"SELECT id FROM users WHERE contract = {contract};")
    id = await cur.fetchall()
    id = id[0]
    await cur.execute("SELECT data FROM netpause WHERE mid = %s", id)
    net = await cur.fetchall()
    try:
        if net[0][0] == 0:
            return True
        elif net[0][0] == 1:
            return False
        else:
            return False
    except IndexError:
        return True


async def tel_by_group(group: int = 8):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(f"SELECT telefon FROM users WHERE grp = {group}")
    phone = await cur.fetchall()
    try:
        phone_array = []
        for i in phone:
            phone_array.append(i[0])
        return phone_array
    except Exception as e:
        print(e)
        pass

async def users_with_alarm(grp, street=None, street_number=None):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    users_with_alarm_contract = []
    if grp and street and street_number:
        pass
    elif grp:
        await cur.execute(f"SELECT contract FROM users WHERE grp = {grp}")
        users = await cur.fetchall()
        for user in users:
            users_with_alarm_contract.append(user[0])
    return users_with_alarm_contract


async def check_contract_exists(contract):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, use_unicode='cp1251')
    cur = await conn.cursor()
    await cur.execute(f"SELECT contract FROM users")
    contracts_in_db = await cur.fetchall()
    if (contract,) in contracts_in_db:
        return True
    else:
        return False
