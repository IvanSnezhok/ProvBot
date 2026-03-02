import asyncio
import aiomysql
import logging
from datetime import datetime
from pytz import timezone
from loader import db, dp
from data import config
from utils.misc.sms_message import send_message_sms
from utils.format_number import number

logger = logging.getLogger(__name__)

async def notify_debtors():
    kyiv_time = datetime.now(timezone('Europe/Kiev'))
    logger.info("Starting debt notification at %s Kyiv time", kyiv_time.strftime('%Y-%m-%d %H:%M:%S'))

    conn = await aiomysql.connect(host=config.BILL_HOST, port=int(config.BILL_PORT),
                                  user=config.BILL_USER, password=config.BILL_PASS,
                                  db=config.BILL_NAME, loop=asyncio.get_event_loop())
    cursor = await conn.cursor()

    # Get tariffs
    await cursor.execute("SELECT i, name, price FROM plans2 WHERE name != '' AND i NOT IN (15, 21)")
    tariffs = await cursor.fetchall()

    for tariff in tariffs:
        tariff_id, tariff_name, tariff_price = tariff
        logger.info("Processing tariff %s; Price: %s GRN", tariff_name, tariff_price)

        query = f"""
        SELECT ip, telefon, fio, balance, contract 
        FROM users 
        WHERE balance - ({tariff_price} - {tariff_price}/100*start_day) < 0 
        AND {tariff_id} = paket 
        AND start_day >= 0 
        AND grp IN (5, 8, 10, 11, 12, 13, 14) 
        AND state = 'on'
        """

        await cursor.execute(query)
        debtors = await cursor.fetchall()

        for debtor in debtors:
            ip, phone, fio, balance, contract = debtor
            phone = ''.join(filter(str.isdigit, phone))
            debt = round(tariff_price - balance, 2)

            if debt > 0:
                logger.debug("Debtor %s:%s; tel:%s debt=%s grn", contract, ip, phone, debt)
                await send_notification(phone, contract, debt)

    await cursor.close()
    conn.close()

async def send_notification(phone, contract, debt):
    message = f"Завтра доступ до Інтернет за договором {contract} буде заблоковано. Поповніть рахунок мінімум на {debt} грн."
    
    try:
        # Використовуємо існуючу функцію send_message_sms
        result = await send_message_sms(number(phone), message)
        logger.info("Notification result for %s: %s", phone, result)
    except Exception as e:
        logger.error("Error sending notification to %s: %s", phone, e)

# Функція для запуску скрипта
async def run_debt_notification():
    await notify_debtors()

# Додайте цю функцію до scheduler у вашому app.py
def schedule_debt_notification(scheduler):
    # Виконувати завдання в останній день місяця о 19:00 за київським часом
    scheduler.add_job(run_debt_notification, "cron", day="last", hour=19, timezone="Europe/Kiev")