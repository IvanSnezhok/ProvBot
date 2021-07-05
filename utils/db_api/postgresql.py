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

    async def select_lang(self, user_id):
        sql = "SELECT lang FROM users WHERE telegram_id=$1"
        return await self.execute(sql, user_id, execute=True, fetchval=True)

    async def select_contract(self, user_id):
        sql = "SELECT contract FROM users WHERE telegram_id=$1"
        return await self.execute(sql, user_id, execute=True, fetch=True)

    async def select_tel(self, user_id):
        sql = "SELECT phone_number FROM users WHERE telegram_id=$1 "
        return await self.execute(sql, user_id, execute=True, fetchval=True)

    async def count_users(self):
        sql = "SELECT COUNT(*) FROM Users"
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
