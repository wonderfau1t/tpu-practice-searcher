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
bot = Bot(token="7650435197:AAFWerS9j6ikk2MFlQ6raTgRFin2BL3FbQ0")
DATABASE_URL = os.getenv("DATABASE_URL")
# Диспетчер
dp = Dispatcher()
pool: asyncpg.Pool = None


async def init_db():
    global pool
    pool = await asyncpg.create_pool(DATABASE_URL)


async def update_phone_number(user_id: int, new_phone_number: str):
    conn = await asyncpg.connect(DATABASE_URL)
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
    logger.info("Отреагировал на команту старт")


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
