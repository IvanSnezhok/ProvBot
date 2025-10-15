from typing import Union

import asyncpg
from asyncpg import Connection
from asyncpg.pool import Pool

from data import config


class Database:

    def __init__(self):
        self.pool: Union[Pool, None] = None

    async def create(self):
        self.pool = await asyncpg.create_pool(
            user=config.DB_USER,
            password=config.DB_PASS,
            host=config.DB_HOST,
            database=config.DB_NAME
        )

    async def execute(self, command, *args,
                      fetch: bool = False,
                      fetchval: bool = False,
                      fetchrow: bool = False,
                      execute: bool = False
                      ):
        async with self.pool.acquire() as connection:
            connection: Connection
            async with connection.transaction():
                if fetch:
                    result = await connection.fetch(command, *args)
                elif fetchval:
                    result = await connection.fetchval(command, *args)
                elif fetchrow:
                    result = await connection.fetchrow(command, *args)
                elif execute:
                    result = await connection.execute(command, *args)
            return result

    async def create_table_users(self):
        sql = """
        CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        full_name VARCHAR(255) NOT NULL,
        username VARCHAR(255) NULL,
        telegram_id BIGINT NOT NULL UNIQUE, 
        lang VARCHAR(5) NULL,
        phone_number VARCHAR(255) NULL, 
        contract VARCHAR(255) NULL      
        );
        """
        await self.execute(sql, execute=True)

    async def create_table_msg(self):
        sql = """
        CREATE TABLE IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY,
        full_name VARCHAR(255) NOT NULL,
        telegram_id BIGINT NOT NULL,
        date TIMESTAMP NOT NULL,
        message VARCHAR(255) NULL
        ); 
        """
        await self.execute(sql, execute=True)

    async def create_table_alarm(self):
        sql = """
        CREATE TABLE IF NOT EXISTS alarm (
        alarm_id SERIAL PRIMARY KEY,
        message VARCHAR(255) NOT NULL,
        grp_alarm VARCHAR(255) NOT NULL,
        street varchar(255),
        street_number varchar(255)
        );
        """
        await self.execute(sql, execute=True)

    async def create_table_bill_check(self):
        sql = """
        CREATE TABLE IF NOT EXISTS bill_check (
        bill_id SERIAL PRIMARY KEY,
        telegram_id BIGINT NOT NULL,
        date TIMESTAMP NOT NULL,
        username varchar(255),
        contract varchar(255),
        pay_amount varchar(255)
        );
        """
        await self.execute(sql, execute=True)

    async def create_table_chats(self):
        sql = """
        CREATE TABLE IF NOT EXISTS chats (
            id SERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL,
            admin_id BIGINT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            closed_at TIMESTAMP,
            status VARCHAR(255) NULL,
            rating INTEGER NULL
        );
        """
        await self.execute(sql, execute=True)

    async def create_table_user_clicks(self):
        sql = """
        CREATE TABLE IF NOT EXISTS user_clicks (
        id SERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL,
        link VARCHAR(255) NOT NULL,
        click_time TIMESTAMP NOT NULL
        );
        """
        await self.execute(sql, execute=True)

    async def add_chat(self, user_id: int, admin_id: int = None, status: str = 'waiting'):
        sql = "INSERT INTO chats (user_id, admin_id, start_time, status) VALUES($1, $2, NOW(), $3) RETURNING id"
        return await self.execute(sql, user_id, admin_id, status, fetchval=True)

    async def update_chat(self, chat_id: int, admin_id: int = None, status: str = None):
        sql = "UPDATE chats SET admin_id = $1, status = $2 WHERE id = $3"
        await self.execute(sql, admin_id, status, chat_id, execute=True)

    async def end_chat(self, user_id: int):
        sql = "UPDATE chats SET closed_at = NOW(), status = 'closed' WHERE user_id = $1 AND status = 'active'"
        await self.execute(sql, user_id, execute=True)

    @staticmethod
    def format_args(sql, parameters: dict):
        sql += " AND ".join([
            f"{item} = ${num}" for num, item in enumerate(parameters.keys(),
                                                          start=1)
        ])
        return sql, tuple(parameters.values())

    async def add_user(self, full_name, username, telegram_id):
        sql = "INSERT INTO users (full_name, username, telegram_id) VALUES($1, $2, $3) returning *"
        return await self.execute(sql, full_name, username, telegram_id, fetchrow=True)

    async def message(self, full_name, telegram_id, message, date):
        sql = "INSERT INTO messages (full_name, telegram_id, message, date) VALUES ($1, $2, $3, $4)"
        return await self.execute(sql, full_name, telegram_id, message, date, execute=True)

    async def select_all_users(self):
        sql = "SELECT * FROM Users"
        return await self.execute(sql, fetch=True)

    async def select_user(self, **kwargs):
        sql = "SELECT * FROM Users WHERE "
        sql, parameters = self.format_args(sql, parameters=kwargs)
        return await self.execute(sql, *parameters, fetchrow=True)

    async def select_user_by_id(self, user_id):
        sql = "SELECT * FROM Users WHERE telegram_id=$1"
        return await self.execute(sql, user_id, fetchrow=True)

    async def select_lang(self, user_id):
        sql = "SELECT lang FROM users WHERE telegram_id=$1"
        return await self.execute(sql, user_id, execute=True, fetchval=True)

    async def select_contract(self, user_id):
        sql = "SELECT contract FROM users WHERE telegram_id=$1"
        return await self.execute(sql, user_id, execute=True, fetch=True)

    async def select_user_id_by_contract(self, contract):
        sql = 'SELECT telegram_id FROM users WHERE contract=$1'
        return await self.execute(sql, contract, fetch=True)

    async def select_tel(self, user_id):
        sql = "SELECT phone_number FROM users WHERE telegram_id=$1 "
        return await self.execute(sql, user_id, execute=True, fetchval=True)

    async def count_users(self):
        sql = "SELECT COUNT(*) FROM Users"
        return await self.execute(sql, fetchval=True)

    async def count_contract_users(self):
        sql = "SELECT COUNT(*) FROM users WHERE contract IS NOT NULL"
        return await self.execute(sql, fetchval=True)

    async def update_user_username(self, username, telegram_id):
        sql = "UPDATE Users SET username=$1 WHERE telegram_id=$2"
        return await self.execute(sql, username, telegram_id, execute=True)

    async def update_phone_number(self, phone_number, telegram_id):
        sql = "UPDATE Users SET phone_number=$1 WHERE telegram_id=$2"
        return await self.execute(sql, phone_number, telegram_id, execute=True)

    async def delete_users(self):
        await self.execute("DELETE FROM Users WHERE TRUE", execute=True)

    async def drop_users(self):
        await self.execute("DROP TABLE Users", execute=True)

    async def set_lang(self, lang, telegram_id):
        sql = "UPDATE users SET lang=$1 WHERE telegram_id=$2"
        await self.execute(sql, lang, telegram_id, execute=True)

    async def set_contract(self, contract, telegram_id):
        sql = "UPDATE users SET contract=$1 WHERE telegram_id=$2"
        await self.execute(sql, contract, telegram_id, execute=True)

    async def choose_contract(self):
        sql = "SELECT full_name, telegram_id, contract FROM users"
        return await self.execute(sql, execute=True, fetch=True)

    async def select_id_by_phone(self, phone_number):
        sql = "SELECT telegram_id FROM users WHERE phone_number=$1"
        return await self.execute(sql, phone_number, execute=True, fetch=True)

    async def get_phone_by_contract(self, contract):
        sql = "SELECT phone_number FROM users WHERE contract=$1"
        return await self.execute(sql, contract, fetchval=True)

    async def set_ban(self, telegram_id):
        sql = "UPDATE users SET ban=TRUE WHERE telegram_id=$1"
        return await self.execute(sql, int(telegram_id), execute=True)

    async def set_unban(self, telegram_id):
        sql = "UPDATE users SET ban=FALSE WHERE telegram_id=$1"
        return await self.execute(sql, int(telegram_id), execute=True)

    async def get_ban(self):
        sql = "SELECT telegram_id, contract FROM users WHERE ban=TRUE"
        rec_ban = await self.execute(sql, fetch=True)
        ban_list = [(i[0], i[1]) for i in rec_ban]
        return ban_list

    async def insert_alarm(self, message, grp_alarm):
        if grp_alarm is None:
            sql = "INSERT INTO alarm (message) VALUES ($1)"
            return await self.execute(sql, message, execute=True)
        else:
            sql = "INSERT INTO alarm (message, grp_alarm) VALUES ($1, $2)"
            return await self.execute(sql, message, grp_alarm, execute=True)

    async def get_alarm(self):
        sql = "SELECT * FROM alarm"
        return await self.execute(sql, fetch=True)

    async def change_alarm_message(self, message, alarm_id):
        sql = "UPDATE alarm SET message=$1 WHERE alarm_id=$2"
        return await self.execute(sql, message, alarm_id, execute=True)

    async def delete_alarm(self, alarm_id):
        sql = "DELETE FROM alarm WHERE alarm_id=$1"
        return await self.execute(sql, int(alarm_id), execute=True)

    async def set_alarm_for_users(self, contract):
        sql = "UPDATE users SET alarm=TRUE WHERE contract=$1"
        return await self.execute(sql, contract, execute=True)

    async def is_alarm(self, telegram_id):
        sql = "SELECT contract FROM users WHERE telegram_id=$1 and alarm=TRUE"
        alarm = await self.execute(sql, telegram_id, fetchval=True)
        if alarm:
            return True
        else:
            return False

    async def get_alarm_message(self, grp_alarm):
        sql = "SELECT message FROM alarm WHERE grp_alarm LIKE '%$1%'"
        return await self.execute(sql, str(grp_alarm), fetchval=True)

    async def add_bill(self, bill_id, telegram_id, date, username, contract, pay_amount):
        sql = ("INSERT INTO bill_check (bill_id, telegram_id, date, username, contract, pay_amount) VALUES"
               " ($1, $2, $3, $4, $5, $6)")
        return await self.execute(sql, bill_id, telegram_id, date, username, contract, pay_amount, execute=True)

    async def get_message_history(self, user_id, message_count: int = 10):
        sql = ("SELECT message FROM messages WHERE telegram_id=$1 ORDER BY date DESC LIMIT $2")
        return await self.execute(sql, user_id, int(message_count), fetch=True)

    async def log_user_click(self, user_id: int, link: str):
        sql = "INSERT INTO user_clicks (user_id, link, click_time) VALUES ($1, $2, NOW())"
        return await self.execute(sql, user_id, link, execute=True)

    async def is_user_in_chat(self, user_id: int):
        sql = "SELECT * FROM chats WHERE user_id = $1 AND status != 'closed'"
        return bool(await self.execute(sql, user_id, fetchrow=True))

    async def create_chat(self, user_id: int, admin_id: int):
        sql = "INSERT INTO chats (user_id, admin_id, status) VALUES ($1, $2, 'active')"
        await self.execute(sql, user_id, admin_id, execute=True)

    async def get_admin_for_user(self, user_id: int):
        sql = "SELECT admin_id FROM chats WHERE user_id = $1 AND status = 'active'"
        result = await self.execute(sql, user_id, fetchrow=True)
        return result['admin_id'] if result else None

    async def get_user_for_admin(self, admin_id: int):
        sql = "SELECT user_id FROM chats WHERE admin_id = $1 AND status = 'active'"
        result = await self.execute(sql, admin_id, fetchrow=True)
        return result['user_id'] if result else None

    async def close_chat(self, user_id: int, admin_id: int):
        sql = "UPDATE chats SET status = 'closed', closed_at = NOW() WHERE user_id = $1 AND admin_id = $2 AND status = 'active'"
        await self.execute(sql, user_id, admin_id, execute=True)

    async def save_rating(self, user_id: int, rating: int):
        sql = "UPDATE chats SET rating = $2 WHERE user_id = $1 AND status = 'closed'"
        await self.execute(sql, user_id, rating, execute=True)
