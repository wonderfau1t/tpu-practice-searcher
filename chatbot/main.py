import asyncio
import logging
import os
import asyncpg

from aiogram import Bot, Dispatcher, types
from aiogram.filters.command import Command
from aiogram import F

# Включаем логирование, чтобы не пропустить важные сообщения
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
# Объект бота
bot = Bot(token=os.getenv("TELEGRAM_BOT_TOKEN"))
DB_HOST = os.getenv("DB_HOST")
DB_NAME = os.getenv("POSTGRES_DB")
DB_USER = os.getenv("POSTGRES_USER")
DB_PASSWORD = os.getenv("POSTGRES_PASSWORD")
DB_URL = f"postgres://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:5432/{DB_NAME}"

# Диспетчер
dp = Dispatcher()
pool: asyncpg.Pool = None


async def init_db():
    global pool
    pool = await asyncpg.create_pool(DB_URL)


async def update_phone_number(user_id: int, new_phone_number: str):
    conn = await asyncpg.connect(DB_URL)
    try:
        await conn.execute(
            "UPDATE users SET phone_number = $1 WHERE id = $2;",
            new_phone_number,
            user_id
        )
    finally:
        await conn.close()


# Хэндлер на команду /start
@dp.message(Command("start"))
async def cmd_start(message: types.Message):
    await message.answer(
        "Привет! Это сервис для поиска практики для студентов ТПУ\n"
        "Переходите в сервис по кнопке \"Запустить\""
    )

@dp.message(F.contact)
async def handle_contact(message: types.Message):
    contact = message.contact
    logger.info(f"Принял контакт: {contact.phone_number}")
    await update_phone_number(message.from_user.id, contact.phone_number)
    await message.answer("Ваш номер телефона успешно подтвержден!")


@dp.message(F.text)
async def handle_text_messages(message: types.Message):
    await message.answer(
        "Я сейчас никак не реагрирую на сообщения\n"
        "Для работы с сервисом тебе необходимо перейти в mini-app по кнопке \"Запустить\""
    )


# Запуск процесса поллинга новых апдейтов
async def main():
    await init_db()
    await bot.delete_webhook(drop_pending_updates=True)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
