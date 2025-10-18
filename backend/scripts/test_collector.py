#!/usr/bin/env python3
"""
Тестовый скрипт для проверки сбора новостей из RSS
"""
import asyncio
import sys

sys.path.insert(0, "/app")

from app.core.database import AsyncSessionLocal
from app.services.collector import NewsCollector


async def test_collector():
    """Тестирование NewsCollector"""
    print("🧪 Testing NewsCollector...")
    print("-" * 50)

    async with AsyncSessionLocal() as db:
        collector = NewsCollector(db)

        # Тест фильтрации по ключевым словам
        print("\n1️⃣ Testing keyword filtering:")

        crypto_text = "Bitcoin price reached $100,000 today"
        is_crypto = collector.filter_by_keywords(crypto_text, "crypto")
        print(f"   Crypto text: '{crypto_text}'")
        print(
            f"   Is crypto: {is_crypto} ✅"
            if is_crypto
            else f"   Is crypto: {is_crypto} ❌"
        )

        politics_text = "Kremlin announced new foreign policy"
        is_politics = collector.filter_by_keywords(politics_text, "politics")
        print(f"   Politics text: '{politics_text}'")
        print(
            f"   Is politics: {is_politics} ✅"
            if is_politics
            else f"   Is politics: {is_politics} ❌"
        )

        # Тест хэширования
        print("\n2️⃣ Testing content hashing:")

        title = "Test News Title"
        content = "Test news content"
        hash1 = collector.calculate_content_hash(title, content)
        hash2 = collector.calculate_content_hash(title, content)

        print(f"   Title: {title}")
        print(f"   Content: {content}")
        print(f"   Hash 1: {hash1}")
        print(f"   Hash 2: {hash2}")
        print(
            f"   Hashes match: {hash1 == hash2} ✅"
            if hash1 == hash2
            else f"   Hashes match: {hash1 == hash2} ❌"
        )

        # Тест сбора новостей
        print("\n3️⃣ Testing news collection:")
        print("   Collecting news from all sources...")

        try:
            result = await collector.collect_all()

            print(f"\n   ✅ Collection completed!")
            print(f"   Total collected: {result['total_collected']}")
            print(f"   Sources:")
            for source_name, count in result["sources"].items():
                print(f"      - {source_name}: {count} news")

        except Exception as e:
            print(f"\n   ❌ Error: {e}")

        print("\n" + "-" * 50)
        print("✨ Test completed!")


if __name__ == "__main__":
    asyncio.run(test_collector())
