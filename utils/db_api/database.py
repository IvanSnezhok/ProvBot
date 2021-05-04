import aiomysql
import asyncio


loop = asyncio.get_event_loop()
data = []
plan = []


async def search_query(tel):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="root", password="password",
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
                                  user="root", password="password",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    await cur.execute("UPDATE users set balance = balance + 150 WHERE contract=%s", contract)
    await cur.close()
    conn.close()


async def pause_inet(contract):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="root", password="password",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    await cur.execute("")