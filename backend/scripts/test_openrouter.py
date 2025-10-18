#!/usr/bin/env python3
"""
Тестовый скрипт для проверки OpenRouter API
"""

import asyncio
import httpx
import os
from dotenv import load_dotenv

load_dotenv()


async def test_openrouter():
    """Тестирование OpenRouter API"""
    api_key = os.getenv("OPENROUTER_API_KEY")
    
    if not api_key:
        print("❌ OPENROUTER_API_KEY не установлен в .env!")
        return
    
    print("🔍 Тестирование OpenRouter API...")
    
    url = "https://openrouter.ai/api/v1/chat/completions"
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }
    
    data = {
        "model": "openai/gpt-3.5-turbo",
        "messages": [
            {
                "role": "user",
                "content": "Кратко опиши новость: 'Bitcoin достиг $100,000'"
            }
        ],
        "max_tokens": 100
    }
    
    try:
        async with httpx.AsyncClient(timeout=30.0) as client:
            response = await client.post(url, json=data, headers=headers)
            response.raise_for_status()
            
            result = response.json()
            message = result["choices"][0]["message"]["content"]
            
            print(f"✅ OpenRouter работает!")
            print(f"📝 Ответ: {message}")
            
    except httpx.HTTPStatusError as e:
        print(f"❌ HTTP ошибка: {e.response.status_code}")
        print(f"   Ответ: {e.response.text}")
    except Exception as e:
        print(f"❌ Ошибка: {e}")


if __name__ == "__main__":
    asyncio.run(test_openrouter())
