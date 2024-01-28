from aiogram.dispatcher.filters.state import State, StatesGroup


class Client(StatesGroup):
    Quest = State()


class Request(StatesGroup):
    Quest = State()
