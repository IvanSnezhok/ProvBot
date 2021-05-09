import asyncio
import logging
import time

import aiomysql

loop = asyncio.get_event_loop()
data = []
plan = []
time_pay = []


async def search_query(tel):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="MySQL", password="M,srHEkK38VB)}5e",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    await cur.execute("SELECT name, balance, contract, fio, state, paket FROM `users` WHERE telefon=%s", tel)
    result = await cur.fetchall()
    data.clear()
    plan.clear()
    try:
        result = result[0]
        plan.append(result[5])
        await cur.execute("SELECT name FROM `plans2` WHERE id=%s", result[5])
        paket = await cur.fetchall()
        paket = paket[0]
        data.append(result[0])
        data.append(result[1])
        data.append(result[2])
        data.append(result[3])
        data.append(result[4])
        data.append(paket[0])

    except IndexError:
        result = None
    await cur.close()
    conn.close()


async def pay_balance_150(contract):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="MySQL", password="M,srHEkK38VB)}5e",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    await cur.execute("UPDATE users set balance = balance + 150 WHERE contract=%s", contract)
    await cur.close()
    conn.close()


async def pay_balance(contract, payload):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="root", password="password",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    execute = payload, contract
    await cur.execute("UPDATE users set balance = balance + %s WHERE contract=%s", execute)
    await cur.close()
    conn.close()


async def t_pay(contract):  # Временный плтажеж
    contract = contract[0]
    contract = contract[0]
    now_time = time.time()
    now_t = time.ctime(now_time)
    next_t = time.time() + 86400

    conn = await aiomysql.connect(host="localhost", port=3306,  # подключение к бд
                                  user="MySQL", password="M,srHEkK38VB)}5e",
                                  db="bill", loop=loop)
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
            await cur.execute("""SELECT price FROM plans2 WHERE id = %s""", paket)
        price = await cur.fetchall()
        try:
            price = price[0]
        except IndexError:
            price = None
            logging.info("Price not find")
            return False

        if old_balance > 0:
            price = paket
            balance = old_balance + price
        else:
            price = paket
            balance = -old_balance + price

        await cur.execute(f"UPDATE users SET t_pay=1 WHERE contract={contract}")
        await cur.execute(f"""INSERT INTO pays (mid,cash,time,bonus,admin,reason,coment,flag)
                VALUES
                ({id},{price},{next_t},'y','timepays','Platej sozdan {now_t}','Razblokirovan na 24 chasa', 't')""")
        await cur.execute(f"UPDATE users SET balance={balance} WHERE contract={contract}")
        await cur.execute(f"UPDATE users SET state='on' WHERE contract={contract}")
        time_pay.append(balance)
        return True
    else:
        return False
