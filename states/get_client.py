from aiogram.dispatcher.filters.state import State, StatesGroup


class Client(StatesGroup):
    Quest = State()


class Request(StatesGroup):
    Quest = State()


class SupportChat(StatesGroup):
    WaitingForSupport = State()
    Chatting = State()


class BroadcastStates(StatesGroup):
    WaitingContent = State()  # Очікування контенту
    PreviewConfirm = State()  # Підтвердження preview
    Broadcasting = State()  # Активна розсилка
