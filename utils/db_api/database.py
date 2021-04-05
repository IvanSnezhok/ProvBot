import aiomysql
import asyncio

loop = asyncio.get_event_loop()
data = []


async def search_query(tel):
    conn = await aiomysql.connect(host="localhost", port=3306,
                                  user="root", password="password",
                                  db="bill", loop=loop)
    cur = await conn.cursor()
    await cur.execute("SELECT name, balance, contract, fio, state, paket FROM `users` WHERE telefon=%s", tel)
    result = await cur.fetchall()
    data.clear()
    try:
        result = result[0]
        data.append(result[0])
        data.append(result[1])
        data.append(result[2])
        data.append(result[3])
        data.append(result[4])
        data.append(result[5])
    except IndexError:
        result = None
    await cur.close()
    conn.close()

