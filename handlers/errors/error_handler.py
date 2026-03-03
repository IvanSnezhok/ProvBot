import logging
from aiogram.utils.exceptions import (Unauthorized, InvalidQueryID, TelegramAPIError,
                                      CantDemoteChatCreator, MessageNotModified, MessageToDeleteNotFound,
                                      MessageTextIsEmpty, RetryAfter,
                                      CantParseEntities, MessageCantBeDeleted)

logger = logging.getLogger(__name__)

from loader import dp


@dp.errors_handler()
async def errors_handler(update, exception):
    """
    Exceptions handler. Catches all exceptions within task factory tasks.
    :param dispatcher:
    :param update:
    :param exception:
    :return: stdout logging
    """

    if isinstance(exception, CantDemoteChatCreator):
        logger.debug("Can't demote chat creator")
        return True

    if isinstance(exception, MessageNotModified):
        logger.debug('Message is not modified')
        return True
    if isinstance(exception, MessageCantBeDeleted):
        logger.debug('Message cant be deleted')
        return True

    if isinstance(exception, MessageToDeleteNotFound):
        logger.debug('Message to delete not found')
        return True

    if isinstance(exception, MessageTextIsEmpty):
        logger.debug('MessageTextIsEmpty')
        return True

    if isinstance(exception, Unauthorized):
        logger.info("Unauthorized: %s", exception)
        return True

    if isinstance(exception, InvalidQueryID):
        logger.exception("InvalidQueryID: %s \nUpdate: %s", exception, update)
        return True

    if isinstance(exception, TelegramAPIError):
        logger.exception("TelegramAPIError: %s \nUpdate: %s", exception, update)
        return True
    if isinstance(exception, RetryAfter):
        logger.exception("RetryAfter: %s \nUpdate: %s", exception, update)
        return True
    if isinstance(exception, CantParseEntities):
        logger.exception("CantParseEntities: %s \nUpdate: %s", exception, update)
        return True
    
    logger.exception("Update: %s \n%s", update, exception)
