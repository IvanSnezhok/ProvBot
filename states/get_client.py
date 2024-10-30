from aiogram.dispatcher.filters.state import State, StatesGroup


class Client(StatesGroup):
    Quest = State()


class Request(StatesGroup):
    Quest = State()

class ChatStates(StatesGroup):
    waiting_for_admin = State()  # Користувач очікує підключення адміна
    in_chat = State()  # Активний чат для користувача
    admin_in_chat = State()  # Активний чат для адміністратора
    rating = State()  # Стан очікування оцінки після завершення чату


