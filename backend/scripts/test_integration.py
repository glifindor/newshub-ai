#!/usr/bin/env python3
"""
Полный интеграционный тест: collect → analyze → post
"""
import asyncio
import sys
sys.path.insert(0, '/app')

from app.core.database import AsyncSessionLocal
from app.services.collector import NewsCollector
from app.services.ai_analyzer import AIAnalyzer
from app.services.telegram_poster import TelegramPoster
from app.core.logger import get_logger

logger = get_logger(__name__)


async def full_integration_test():
    """Полный тест pipeline"""
    print("=" * 70)
    print("🚀 NewsHub AI - Full Integration Test")
    print("=" * 70)
    
    async with AsyncSessionLocal() as db:
        # Step 1: Collection
        print("\n📥 STEP 1: Collecting news...")
        print("-" * 70)
        
        collector = NewsCollector(db)
        collection_result = await collector.collect_all()
        
        print(f"✅ Collected {collection_result['total_collected']} news items")
        for source, count in collection_result['sources'].items():
            print(f"   • {source}: {count}")
        
        # Step 2: Analysis
        print("\n🤖 STEP 2: Analyzing news with AI...")
        print("-" * 70)
        
        analyzer = AIAnalyzer(db)
        analysis_result = await analyzer.analyze_pending_news(limit=5)
        
        print(f"✅ Analyzed {analysis_result['total']} news items")
        print(f"   • Analyzed: {analysis_result['analyzed']}")
        print(f"   • Rejected: {analysis_result['rejected']}")
        print(f"   • Failed: {analysis_result['failed']}")
        
        # Step 3: Review
        print("\n📊 STEP 3: Reviewing analyzed news...")
        print("-" * 70)
        
        from sqlalchemy import select
        from app.models import NewsItem, NewsStatus
        
        result = await db.execute(
            select(NewsItem)
            .where(NewsItem.status == NewsStatus.ANALYZED)
            .order_by(NewsItem.relevance_score.desc())
            .limit(5)
        )
        news_items = result.scalars().all()
        
        if news_items:
            print(f"Found {len(news_items)} analyzed news items:\n")
            
            for i, news in enumerate(news_items, 1):
                print(f"{i}. {news.title[:60]}...")
                print(f"   📊 Relevance: {news.relevance_score}/10")
                print(f"   🏷️  Category: {news.category.value}")
                print(f"   📝 AI Summary: {news.ai_summary[:100]}...")
                print(f"   🔗 URL: {news.url}")
                print()
        else:
            print("⚠️  No analyzed news found")
        
        # Step 4: Posting (опционально)
        print("\n📤 STEP 4: Posting to Telegram...")
        print("-" * 70)
        
        user_input = input("Do you want to post news to Telegram? (y/N): ")
        
        if user_input.lower() == 'y':
            poster = TelegramPoster(db)
            
            # Сначала модерация
            moderation_result = await poster.handle_moderation_requests()
            print(f"✅ Moderation: {moderation_result['notified']} notifications sent")
            
            # Затем автопост
            post_result = await poster.post_analyzed_news(limit=3)
            print(f"✅ Posted {post_result['posted']} news items to Telegram")
            print(f"   • Failed: {post_result['failed']}")
        else:
            print("⏭️  Skipping Telegram posting")
        
        print("\n" + "=" * 70)
        print("✨ Integration test completed!")
        print("=" * 70)


if __name__ == "__main__":
    asyncio.run(full_integration_test())
