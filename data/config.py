from pathlib import Path

from environs import Env

# Теперь используем вместо библиотеки python-dotenv библиотеку environs
env = Env()
env.read_env()

BOT_TOKEN = env.str("BOT_TOKEN")  # Забираем значение типа str
ADMINS = env.list("ADMINS")  # Тут у нас будет список из админов
IP = env.str("ip")  # Тоже str, но для айпи адреса хоста

DB_USER = env.str("DB_USER")
DB_PASS = env.str("DB_PASS")
DB_NAME = env.str("DB_NAME")
DB_HOST = env.str("DB_HOST")

PROVIDER_TOKEN = env.str("PROVIDER_TOKEN")

I18N_DOMAIN = "prov_bot"
BASE_DIR = Path(__file__).parent
LOCALES_DIR = BASE_DIR/'locales'

BILL_USER = env.str("BILL_USER")
BILL_PASS = env.str("BILL_PASS")
BILL_NAME = env.str("BILL_NAME")
BILL_HOST = env.str("BILL_HOST")
BILL_PORT = env.str("BILL_PORT")
