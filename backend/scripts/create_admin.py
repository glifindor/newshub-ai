#!/usr/bin/env python3
"""
Скрипт для создания администратора системы
Запуск: python scripts/create_admin.py
"""

import asyncio
import sys
from getpass import getpass

# Добавляем путь к приложению
sys.path.insert(0, '/app')

from app.core.database import AsyncSessionLocal
from app.core.security import get_password_hash


async def create_admin():
    """Создать администратора"""
    print("🔧 Создание администратора NewsHub AI\n")
    
    # Ввод данных
    username = input("Username: ").strip()
    if not username:
        print("❌ Username не может быть пустым!")
        return
    
    email = input("Email: ").strip()
    if not email:
        print("❌ Email не может быть пустым!")
        return
    
    password = getpass("Password: ")
    if len(password) < 6:
        print("❌ Пароль должен быть минимум 6 символов!")
        return
    
    password_confirm = getpass("Confirm password: ")
    if password != password_confirm:
        print("❌ Пароли не совпадают!")
        return
    
    # Хэширование пароля
    hashed_password = get_password_hash(password)
    
    # TODO: Сохранение в БД (после создания модели User)
    print(f"\n✅ Админ создан!")
    print(f"   Username: {username}")
    print(f"   Email: {email}")
    print(f"   Password hash: {hashed_password[:20]}...")
    print("\n⚠️  ВАЖНО: Сохраните эти данные в безопасном месте!")


if __name__ == "__main__":
    asyncio.run(create_admin())
