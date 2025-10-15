from aiogram import Dispatcher
#  from aiogram.contrib.middlewares.logging import LoggingMiddleware


from loader import dp
from .throttling import ThrottlingMiddleware
from .language_middleware import ACLMiddleware, setup_middleware


if __name__ == "middlewares":
    #  dp.middleware.setup(LoggingMiddleware())
    dp.middleware.setup(ThrottlingMiddleware())
    i18n = setup_middleware(dp)
    _ = i18n.gettext
    __ = i18n.lazy_gettext
