import time
from typing import Tuple, Any, Dict, Optional

from aiogram import types
from aiogram.contrib.middlewares.i18n import I18nMiddleware
from data.config import I18N_DOMAIN, LOCALES_DIR
from loader import db, dp

_lang_cache: Dict[int, Tuple[str, float]] = {}
_LANG_CACHE_TTL = 300  # 5 minutes


async def get_lang(user_id):
    now = time.time()
    cached = _lang_cache.get(user_id)
    if cached and (now - cached[1]) < _LANG_CACHE_TTL:
        return cached[0]

    lang = await db.select_lang(user_id)
    if lang:
        _lang_cache[user_id] = (lang, now)
        return lang
    return None


def invalidate_lang_cache(user_id: int):
    _lang_cache.pop(user_id, None)


class ACLMiddleware(I18nMiddleware):
    async def get_user_locale(self, action, args):
        user = types.User.get_current()

        return await get_lang(user.id) or user.locale


def setup_middleware(dp):
    i18n = ACLMiddleware(I18N_DOMAIN, LOCALES_DIR)
    dp.middleware.setup(i18n)
    return i18n
