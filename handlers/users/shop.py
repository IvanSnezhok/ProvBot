from aiogram import types
from aiogram.dispatcher import FSMContext
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ReplyKeyboardMarkup, KeyboardButton
from collections import defaultdict

from loader import dp, db
from utils.shop_parser import get_products
from middlewares import _, __

# Зберігаємо стан для пагінації та категорій
class ShopStates:
    current_page = defaultdict(int)  # user_id -> current_page
    items_per_page = 1
    category_dict = {}  # index -> category_name

@dp.message_handler(text=__("🛒 Магазин"))
async def show_categories(message: types.Message):
    products = await get_products()

    # Отримуємо унікальні категорії
    categories = set(product['category'] for product in products)
    
    # Створюємо клавіатуру з категоріями
    keyboard = InlineKeyboardMarkup(row_width=1)
    
    # Створюємо словник для зберігання відповідності індексів та категорій
    category_dict = {str(i): cat for i, cat in enumerate(categories)}
    
    # Зберігаємо словник в базі даних або в пам'яті для подальшого використання
    ShopStates.category_dict = category_dict
    
    for index, category in category_dict.items():
        keyboard.add(InlineKeyboardButton(
            text=__(category.split(' > ')[-1]),  # Беремо останню частину шляху категорії
            callback_data=f"cat_{index}"  # Використовуємо індекс замість повної назви категорії
        ))
    
    await message.answer(__("Оберіть категорію товарів:"), reply_markup=keyboard)

@dp.callback_query_handler(lambda c: c.data.startswith('cat_'))
async def show_category_products(callback: types.CallbackQuery):
    category_index = callback.data.replace('cat_', '')
    category = ShopStates.category_dict.get(category_index)
    
    if not category:
        await callback.answer(__("Категорія не знайдена"))
        return
    
    products = [p for p in await get_products() if p['category'] == category]

    if not products:
        await callback.answer(__("В цій категорії зараз немає товарів"))
        return
    
    # Скидаємо сторінку для користувача
    ShopStates.current_page[callback.from_user.id] = 0
    
    await show_product_page(callback.message, products, callback.from_user.id)
    await callback.answer()

async def show_product_page(message: types.Message, products: list, user_id: int):
    current_page = ShopStates.current_page[user_id]
    total_pages = len(products)
    
    if current_page >= total_pages:
        current_page = 0
        ShopStates.current_page[user_id] = 0
    elif current_page < 0:
        current_page = total_pages - 1
        ShopStates.current_page[user_id] = current_page
    
    product = products[current_page]
    
    # Створюємо клавіатуру
    keyboard = InlineKeyboardMarkup(row_width=2)
    
    # Додаємо кнопку "Замовити"
    keyboard.add(InlineKeyboardButton(
        text=__("🛍 Замовити"),
        url=product['link']
    ))
    
    # Додаємо навігаційні кнопки
    nav_buttons = []
    
    # Знаходимо індекс категорії
    category_index = None
    for idx, cat in ShopStates.category_dict.items():
        if cat == product['category']:
            category_index = idx
            break
            
    if category_index is not None:
        if current_page > 0:
            nav_buttons.append(InlineKeyboardButton("⬅️", callback_data=f"prev_{category_index}"))
        if current_page < total_pages - 1:
            nav_buttons.append(InlineKeyboardButton("➡️", callback_data=f"next_{category_index}"))
        
        if nav_buttons:
            keyboard.add(*nav_buttons)
    
    # Додаємо кнопки навігації по меню
    keyboard.add(InlineKeyboardButton(__("⬅️ До категорій"), callback_data="show_categories"))
    keyboard.add(InlineKeyboardButton(__("🏠 Головне меню"), callback_data="return_main"))
    
    # Формуємо підпис з номером сторінки
    caption = (
        f"{product['title']}\n\n"
        f"Ціна: {product['price']}\n\n"
        f"📄 {current_page + 1}/{total_pages}"
    )
    
    try:
        await message.edit_media(
            types.InputMediaPhoto(
                media=product['image_link'],
                caption=caption
            ),
            reply_markup=keyboard
        )
    except Exception:
        await message.answer_photo(
            photo=product['image_link'],
            caption=caption,
            reply_markup=keyboard
        )

@dp.callback_query_handler(lambda c: c.data == "show_categories")
async def return_to_categories(callback: types.CallbackQuery):
    await show_categories(callback.message)
    await callback.answer()

@dp.callback_query_handler(lambda c: c.data.startswith(('next_', 'prev_')))
async def navigate_products(callback: types.CallbackQuery):
    action, category_index = callback.data.split('_', 1)
    category = ShopStates.category_dict.get(category_index)
    
    if not category:
        await callback.answer(__("Категорія не знайдена"))
        return
        
    products = [p for p in await get_products() if p['category'] == category]

    if action == 'next':
        ShopStates.current_page[callback.from_user.id] += 1
    else:
        ShopStates.current_page[callback.from_user.id] -= 1
    
    await show_product_page(callback.message, products, callback.from_user.id)
    await callback.answer()

# Додаємо обробник для логування кліків по посиланню "Замовити"
@dp.callback_query_handler(lambda c: c.data.startswith('http'))
async def log_order_click(callback_query: types.CallbackQuery):
    await db.log_user_click(callback_query.from_user.id, callback_query.data)
    await callback_query.answer(_("Переходимо на сайт для замовлення..."))
