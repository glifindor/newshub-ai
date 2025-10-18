#!/usr/bin/env python3
"""
Тестовый скрипт для проверки Telegram Bot
"""

import asyncio
import os
from dotenv import load_dotenv
from telegram import Bot

load_dotenv()


async def test_telegram():
    """Тестирование Telegram Bot"""
    bot_token = os.getenv("TELEGRAM_BOT_TOKEN")
    admin_chat_id = os.getenv("TELEGRAM_ADMIN_CHAT_ID")
    
    if not bot_token:
        print("❌ TELEGRAM_BOT_TOKEN не установлен в .env!")
        return
    
    print("🤖 Тестирование Telegram Bot...")
    
    try:
        bot = Bot(token=bot_token)
        
        # Получить информацию о боте
        bot_info = await bot.get_me()
        print(f"✅ Бот подключен!")
        print(f"   Имя: {bot_info.first_name}")
        print(f"   Username: @{bot_info.username}")
        print(f"   ID: {bot_info.id}")
        
        # Отправить тестовое сообщение админу
        if admin_chat_id:
            await bot.send_message(
                chat_id=admin_chat_id,
                text="🚀 NewsHub AI успешно подключен!\n\nБот готов к работе."
            )
            print(f"\n✉️  Тестовое сообщение отправлено в чат {admin_chat_id}")
        
    except Exception as e:
        print(f"❌ Ошибка: {e}")


if __name__ == "__main__":
    asyncio.run(test_telegram())
