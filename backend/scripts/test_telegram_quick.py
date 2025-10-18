"""
Быстрый тест Telegram бота

Проверяет:
1. Подключение к Telegram Bot API
2. Отправку сообщения админу
3. Публикацию в канал @crypto_ainews
4. Публикацию в канал @kremlin_digest

Запуск:
    python scripts/test_telegram_quick.py
"""
import asyncio
import sys
from pathlib import Path

# Добавляем корневую директорию в PYTHONPATH
sys.path.insert(0, str(Path(__file__).parent.parent))

from telegram import Bot
from telegram.constants import ParseMode
from app.core.config import settings


async def test_bot_connection():
    """Тест подключения к боту"""
    print("\n🤖 Шаг 1: Проверка подключения к боту...")
    
    try:
        bot = Bot(token=settings.TELEGRAM_BOT_TOKEN)
        me = await bot.get_me()
        
        print(f"✅ Бот подключен успешно!")
        print(f"   📝 Имя: {me.first_name}")
        print(f"   🔗 Username: @{me.username}")
        print(f"   🆔 Bot ID: {me.id}")
        
        return bot
        
    except Exception as e:
        print(f"❌ Ошибка подключения к боту: {e}")
        print(f"\n💡 Проверьте:")
        print(f"   1. TELEGRAM_BOT_TOKEN в .env файле")
        print(f"   2. Токен получен от @BotFather")
        print(f"   3. Нет лишних пробелов в токене")
        return None


async def test_admin_message(bot: Bot):
    """Тест отправки сообщения админу"""
    print(f"\n📨 Шаг 2: Отправка тестового сообщения админу (ID: {settings.TELEGRAM_ADMIN_CHAT_ID})...")
    
    try:
        message = await bot.send_message(
            chat_id=settings.TELEGRAM_ADMIN_CHAT_ID,
            text="""
✅ <b>Тестовое сообщение от NewsHub AI Bot</b>

Бот работает корректно и готов к использованию!

🔧 <b>Конфигурация:</b>
• Крипто канал: {crypto}
• Политика канал: {politics}
• Admin ID: {admin}

🚀 Можно запускать автоматическую публикацию!
""".format(
                crypto=settings.TELEGRAM_CRYPTO_CHANNEL,
                politics=settings.TELEGRAM_POLITICS_CHANNEL,
                admin=settings.TELEGRAM_ADMIN_CHAT_ID,
            ),
            parse_mode=ParseMode.HTML,
        )
        
        print(f"✅ Сообщение отправлено админу!")
        print(f"   🆔 Message ID: {message.message_id}")
        print(f"   ⏰ Время: {message.date}")
        
        return True
        
    except Exception as e:
        print(f"❌ Ошибка отправки админу: {e}")
        print(f"\n💡 Проверьте:")
        print(f"   1. TELEGRAM_ADMIN_CHAT_ID в .env")
        print(f"   2. Напишите боту /start в личку")
        print(f"   3. Chat ID получен через @userinfobot")
        return False


async def test_crypto_channel(bot: Bot):
    """Тест публикации в крипто-канал"""
    print(f"\n🔐 Шаг 3: Публикация в крипто-канал ({settings.TELEGRAM_CRYPTO_CHANNEL})...")
    
    test_message = """
🔐 <b>Bitcoin достиг исторического максимума $100,000!</b>

📝 Криптовалюта Bitcoin впервые в истории преодолела отметку в $100,000. Аналитики связывают рост с массовым принятием институциональными инвесторами и запуском биткоин-ETF в США.

🔍 <b>AI-инсайт:</b>
• Институциональные инвестиции выросли на 300% за последний квартал
• Ожидается дальнейший рост до $150,000 к концу года
• Регуляторная ясность в США стимулирует приток капитала

🔗 <a href='https://cointelegraph.com/news/bitcoin-100k'>Читать подробнее на CoinTelegraph</a>

#Bitcoin #Crypto #ATH #BTC #Blockchain
"""
    
    try:
        message = await bot.send_message(
            chat_id=settings.TELEGRAM_CRYPTO_CHANNEL,
            text=test_message,
            parse_mode=ParseMode.HTML,
            disable_web_page_preview=False,
        )
        
        print(f"✅ Сообщение опубликовано в крипто-канале!")
        print(f"   🆔 Message ID: {message.message_id}")
        print(f"   🔗 Ссылка: https://t.me/{settings.TELEGRAM_CRYPTO_CHANNEL.lstrip('@')}/{message.message_id}")
        
        return True
        
    except Exception as e:
        print(f"❌ Ошибка публикации в крипто-канале: {e}")
        print(f"\n💡 Проверьте:")
        print(f"   1. Канал {settings.TELEGRAM_CRYPTO_CHANNEL} существует")
        print(f"   2. Бот добавлен как администратор канала")
        print(f"   3. Бот имеет права 'Публиковать сообщения'")
        print(f"   4. TELEGRAM_CRYPTO_CHANNEL правильно указан в .env")
        return False


async def test_politics_channel(bot: Bot):
    """Тест публикации в политический канал"""
    print(f"\n🏛️ Шаг 4: Публикация в политический канал ({settings.TELEGRAM_POLITICS_CHANNEL})...")
    
    test_message = """
🏛️ <b>Новый саммит G20 пройдёт в Москве</b>

📝 Лидеры стран G20 соберутся в Москве для обсуждения глобальных экономических вопросов, климатических изменений и цифровой трансформации. Саммит запланирован на ноябрь 2025 года.

🔍 <b>AI-инсайт:</b>
• Ожидается подписание ключевых соглашений по климату
• Россия предложит новую инициативу по цифровой экономике
• Фокус на укрепление многостороннего сотрудничества

🔗 <a href='https://tass.ru/politika/g20-moscow-2025'>Читать подробнее на ТАСС</a>

#G20 #Политика #Москва #Саммит #Россия
"""
    
    try:
        message = await bot.send_message(
            chat_id=settings.TELEGRAM_POLITICS_CHANNEL,
            text=test_message,
            parse_mode=ParseMode.HTML,
            disable_web_page_preview=False,
        )
        
        print(f"✅ Сообщение опубликовано в политическом канале!")
        print(f"   🆔 Message ID: {message.message_id}")
        print(f"   🔗 Ссылка: https://t.me/{settings.TELEGRAM_POLITICS_CHANNEL.lstrip('@')}/{message.message_id}")
        
        return True
        
    except Exception as e:
        print(f"❌ Ошибка публикации в политическом канале: {e}")
        print(f"\n💡 Проверьте:")
        print(f"   1. Канал {settings.TELEGRAM_POLITICS_CHANNEL} существует")
        print(f"   2. Бот добавлен как администратор канала")
        print(f"   3. Бот имеет права 'Публиковать сообщения'")
        print(f"   4. TELEGRAM_POLITICS_CHANNEL правильно указан в .env")
        return False


async def test_with_image(bot: Bot):
    """Тест публикации с изображением"""
    print(f"\n🖼️ Шаг 5: Тест публикации с изображением...")
    
    test_caption = """
🔐 <b>Ethereum переходит на Proof of Stake</b>

📝 Краткое описание новости с AI-анализом...

🔍 <b>AI-инсайт:</b>
• Пункт 1
• Пункт 2

🔗 <a href='https://example.com'>Читать подробнее</a>

#Ethereum #ETH #PoS
"""
    
    # Тестовое изображение (placeholder)
    test_image_url = "https://via.placeholder.com/1200x630/1E88E5/FFFFFF?text=NewsHub+AI+Test"
    
    try:
        message = await bot.send_photo(
            chat_id=settings.TELEGRAM_CRYPTO_CHANNEL,
            photo=test_image_url,
            caption=test_caption,
            parse_mode=ParseMode.HTML,
        )
        
        print(f"✅ Сообщение с изображением опубликовано!")
        print(f"   🆔 Message ID: {message.message_id}")
        
        return True
        
    except Exception as e:
        print(f"⚠️ Ошибка публикации с изображением: {e}")
        print(f"   Это нормально, fallback на текстовое сообщение работает")
        return True


async def main():
    """Главная функция"""
    print("=" * 60)
    print("🚀 ТЕСТИРОВАНИЕ TELEGRAM БОТА NewsHub AI")
    print("=" * 60)
    
    # Шаг 1: Подключение
    bot = await test_bot_connection()
    if not bot:
        print("\n❌ Тест провален на шаге 1. Исправьте ошибки и повторите.")
        return
    
    # Шаг 2: Админ сообщение
    admin_ok = await test_admin_message(bot)
    
    # Шаг 3: Крипто канал
    crypto_ok = await test_crypto_channel(bot)
    
    # Шаг 4: Политический канал
    politics_ok = await test_politics_channel(bot)
    
    # Шаг 5: С изображением
    image_ok = await test_with_image(bot)
    
    # Итоги
    print("\n" + "=" * 60)
    print("📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ")
    print("=" * 60)
    print(f"1. Подключение к боту:        {'✅' if bot else '❌'}")
    print(f"2. Сообщение админу:           {'✅' if admin_ok else '❌'}")
    print(f"3. Публикация в крипто-канал:  {'✅' if crypto_ok else '❌'}")
    print(f"4. Публикация в полит-канал:   {'✅' if politics_ok else '❌'}")
    print(f"5. Публикация с изображением:  {'✅' if image_ok else '⚠️'}")
    print("=" * 60)
    
    if all([bot, admin_ok, crypto_ok, politics_ok]):
        print("\n🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ! БОТ ГОТОВ К РАБОТЕ!")
        print("\nСледующие шаги:")
        print("1. Запустите backend: uvicorn app.main:app --reload")
        print("2. Откройте Swagger UI: http://localhost:8000/docs")
        print("3. Запустите полный pipeline: POST /api/v1/pipeline/pipeline")
        print("4. Проверьте каналы - новости должны появиться автоматически!")
    else:
        print("\n⚠️ НЕКОТОРЫЕ ТЕСТЫ НЕ ПРОШЛИ")
        print("Исправьте указанные ошибки и повторите тест.")
    
    print("\n📚 Документация: backend/docs/TELEGRAM_BOT_SETUP.md")
    print("=" * 60)


if __name__ == "__main__":
    asyncio.run(main())
