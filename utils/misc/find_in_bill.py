import aiomysql
import asyncio
import time

from data import config

loop = asyncio.get_event_loop()


async def find(contract=None, phone=None, name=None, address: list = None):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, charset="cp1251")
    cur = await conn.cursor()
    await cur.execute("SELECT name_street FROM `p_street`")
    streets = await cur.fetchall()
    streets = [street[0] for street in streets]
    print(streets)
    try:
        if address:
            if address[0] in streets:
                if len(address) <= 1:
                    for street in range(len(streets)):
                        if streets[street] == address[0]:
                            await cur.execute("SELECT name, balance, contract, fio, state, paket, telefon, ip "
                                              f"FROM `users` WHERE  street = %s", (street,))
                            result = await cur.fetchall()
                            break
                if len(address) <= 2:
                    for street in range(len(streets)):
                        if streets[street] == address[0]:
                            await cur.execute("SELECT name, balance, contract, fio, state, paket, telefon, ip "
                                              f"FROM `users` WHERE  street= %s AND house = %s", (street,
                                                                                                 address[1]))
                            result = await cur.fetchall()
                            break
                if len(address) <= 3:
                    for street in range(len(streets)):
                        if streets[street] == address[0]:
                            await cur.execute("SELECT name, balance, contract, fio, state, paket, telefon, ip "
                                              f"FROM `users` WHERE  street= %s AND house = %s"
                                              f" AND room = %s", (street, address[1], address[2]))
                            result = await cur.fetchall()
                            break
                else:
                    raise ValueError("Too many arguments")

    except TypeError as e:
        print(e)
    if contract:
        await cur.execute("SELECT name, balance, contract, fio, state, paket, telefon, street, house, room, ip, id "
                          "FROM `users` "
                          f"WHERE contract LIKE '{contract}%' ")
    if phone:
        await cur.execute("SELECT name, balance, contract, fio, state, paket, telefon, street, house, room, ip, id"
                          "FROM `users` "
                          f"WHERE telefon LIKE '%{phone}%' ")
    if name:
        sql = "SELECT name, balance, contract, fio, state, paket, telefon, street, house, room, ip, id " \
              "FROM `users` " \
              f"WHERE fio LIKE '%{name}%' "
        print(sql)
        await cur.execute(sql.encode('cp1251'))
    result = await cur.fetchall()
    print(result)

    if len(result) == 1:
        result = result[0]
        address = f"{streets[result[7] - 1]} {result[8]}, кв {result[9]}"
        result_dict = {'name': result[0],
                       'balance': result[1],
                       'contract': result[2],
                       'fio': result[3],
                       'state': result[4],
                       'paket': result[5],
                       'telefon': result[6],
                       'ip': result[10],
                       'id': result[11],
                       'address': address}
        await cur.execute('SELECT name FROM `plans2` WHERE id=%s', result_dict['paket'])
        paket = await cur.fetchall()
        paket = paket[0]
        result_dict['paket'] = paket[0]
        await cur.execute("select stop_time from netpause where mid =%s ORDER BY `stop_time` DESC LIMIT 1;",
                          result_dict['id'])
        net_pause_epoch = await cur.fetchall()
        if net_pause_epoch:
            net_pause_epoch = net_pause_epoch[0]
            if net_pause_epoch[0] > int(time.time()):
                result_dict['state'] = 'pause'
            else:
                pass
        await cur.close()
        conn.close()
        return result_dict
    else:
        temp_list = []
        for i in result:
            result_temp = i
            address = f"{streets[result_temp[7] - 1]} {result_temp[8]}, кв {result_temp[9]}"
            result_dict = {'name': result_temp[0],
                           'balance': result_temp[1],
                           'contract': result_temp[2],
                           'fio': result_temp[3],
                           'state': result_temp[4],
                           'paket': result_temp[5],
                           'telefon': result_temp[6],
                           'ip': result_temp[10],
                           'address': address}
            await cur.execute('SELECT name FROM `plans2` WHERE id=%s', result_dict['paket'])
            paket = await cur.fetchall()
            await cur.execute("select stop_time from netpause where mid =%s ORDER BY `stop_time` DESC LIMIT 1;",
                              result_dict['id'])
            net_pause_epoch = await cur.fetchall()
            if net_pause_epoch:
                net_pause_epoch = net_pause_epoch[0]
                if net_pause_epoch[0] > int(time.time()):
                    result_dict['state'] = 'pause'
                else:
                    pass
            temp_list.append(result_dict)
        await cur.close()
        conn.close()
        return temp_list


async def active_users(contract):
    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=loop, charset="cp1251")
    cur = await conn.cursor()
    await cur.execute("SELECT balance, state, paket, grp, ip , id"
                      "FROM `users` "
                      f"WHERE contract LIKE '{contract}%' ")
    result = await cur.fetchall()
    result = result[0]
    balance = result[0]
    state = result[1]
    paket = result[2]
    grp = result[3]
    ip = result[4]
    id = result[5]
    epoch = int(time.time())
    await cur.execute("select stop_time from netpause where mid =%s ORDER BY `stop_time` DESC LIMIT 1;", id)
    net_pause_epoch = await cur.fetchall()
    net_pause_epoch = net_pause_epoch[0]
    await cur.execute('SELECT name FROM `plans2` WHERE id=%s', paket)
    paket = await cur.fetchall()
    paket = paket[0]
    await cur.close()

    conn.close()
    if grp == 1:
        return "Блок."
    elif grp == 7 or net_pause_epoch[0] < epoch:
        return "Пауза"
    elif balance >= 0 and state == 'on':
        return "Активна"
    else:
        return "Блок."
    # if state == 'on' and balance >= 0 and grp != 1:
    #     return "Активна"
    # elif state == 'off' and balance >= and grp != 1:
    #     return "Блок."
    # elif state == 'off' and balance < 0 and grp != 1:
    #     return "Блок."
    # elif grp == 1:
    #     return "Блок."
    # elif grp == 7:
    #     return "Пауза"
    # else:
    #     return "Блок."
