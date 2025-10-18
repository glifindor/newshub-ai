"""
AI Analyzer Service - анализ новостей с помощью OpenRouter API
"""

import asyncio
import json
from datetime import datetime
from typing import Any, Dict, Optional

import httpx
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from tenacity import (retry, retry_if_exception_type, stop_after_attempt,
                      wait_exponential)

from app.core.config import settings
from app.core.logger import get_logger
from app.models import AITask, NewsChannel, NewsItem, NewsStatus

logger = get_logger(__name__)


class AIAnalyzer:
    """AI-анализатор новостей через OpenRouter"""

    def __init__(self, db: AsyncSession):
        self.db = db
        self.api_url = f"{settings.OPENROUTER_API_URL}/chat/completions"
        self.api_key = settings.OPENROUTER_API_KEY
        self.model = settings.OPENROUTER_MODEL

        # Альтернативные модели (fallback)
        self.fallback_models = [
            "anthropic/claude-3-haiku",
            "meta-llama/llama-3-70b-instruct",
            "google/gemini-pro",
        ]

    def create_analysis_prompt(
        self, title: str, content: str, channel: NewsChannel
    ) -> str:
        """Создать промпт для AI-анализа"""

        if channel == NewsChannel.CRYPTO:
            context = """Ты — профессиональный криптоаналитик и журналист. 
Анализируй новости о криптовалютах, блокчейне, DeFi, NFT для инвесторов и гиков."""
            insights_focus = (
                "что это значит для инвесторов, трейдеров и крипто-сообщества"
            )
        else:
            context = """Ты — политический аналитик и журналист. 
Анализируй новости о России, мировой политике, геополитике."""
            insights_focus = (
                "что это значит для России, мировой политики и международных отношений"
            )

        prompt = f"""{context}

Новость:
Заголовок: {title}
Содержание: {content}

Задача:
1. **Краткий тизер** (80-120 слов, русский язык) — перескажи суть новости увлекательно и понятно.
2. **AI-инсайты** (2-3 пункта) — {insights_focus}.
3. **Рейтинг релевантности** (0-10) — насколько новость важна и интересна для целевой аудитории.
4. **Хэштеги** (3-5 штук) — релевантные хэштеги для Telegram.

Верни ответ СТРОГО в формате JSON:
{{
    "teaser": "краткий тизер на русском",
    "insights": ["инсайт 1", "инсайт 2", "инсайт 3"],
    "relevance_score": 8,
    "hashtags": ["#Crypto", "#Bitcoin", "#AI"]
}}

Если новость НЕ релевантна теме или рейтинг < 7, верни relevance_score = 0."""

        return prompt

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=2, max=10),
        retry=retry_if_exception_type((httpx.HTTPError, httpx.TimeoutException)),
    )
    async def call_openrouter(
        self, prompt: str, model: Optional[str] = None
    ) -> Dict[str, Any]:
        """Вызов OpenRouter API с retry логикой"""

        model = model or self.model

        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
            "HTTP-Referer": settings.NEXT_PUBLIC_SITE_URL,
            "X-Title": "NewsHub AI",
        }

        payload = {
            "model": model,
            "messages": [{"role": "user", "content": prompt}],
            "temperature": 0.7,
            "max_tokens": 800,
        }

        logger.info(f"Calling OpenRouter API with model: {model}")

        try:
            async with httpx.AsyncClient(timeout=settings.AI_TIMEOUT_SECONDS) as client:
                response = await client.post(
                    self.api_url, json=payload, headers=headers
                )
                response.raise_for_status()

                data = response.json()
                content = data["choices"][0]["message"]["content"]

                # Попытка парсинга JSON из ответа
                # Иногда модель возвращает ```json ... ```, удаляем это
                if "```json" in content:
                    content = content.split("```json")[1].split("```")[0].strip()
                elif "```" in content:
                    content = content.split("```")[1].split("```")[0].strip()

                result = json.loads(content)

                logger.info(
                    f"OpenRouter API call successful, relevance: {result.get('relevance_score')}"
                )

                return result

        except httpx.HTTPStatusError as e:
            logger.error(
                f"OpenRouter API HTTP error: {e.response.status_code}",
                response=e.response.text,
            )
            raise
        except json.JSONDecodeError as e:
            logger.error(
                f"Failed to parse OpenRouter response as JSON: {e}", content=content
            )
            # Возвращаем дефолтное значение
            return {
                "teaser": "Ошибка анализа",
                "insights": [],
                "relevance_score": 0,
                "hashtags": [],
            }
        except Exception as e:
            logger.error(f"OpenRouter API error: {e}", exc_info=True)
            raise

    async def analyze_with_fallback(self, prompt: str) -> Dict[str, Any]:
        """Анализ с fallback на альтернативные модели"""

        # Пробуем основную модель
        try:
            return await self.call_openrouter(prompt, self.model)
        except Exception as e:
            logger.warning(f"Primary model failed: {e}, trying fallback models")

        # Пробуем fallback модели
        for fallback_model in self.fallback_models:
            try:
                logger.info(f"Trying fallback model: {fallback_model}")
                return await self.call_openrouter(prompt, fallback_model)
            except Exception as e:
                logger.warning(f"Fallback model {fallback_model} failed: {e}")
                continue

        # Если все модели провалились
        logger.error("All AI models failed, returning default analysis")
        return {
            "teaser": "Анализ временно недоступен",
            "insights": ["AI-анализ не удалось выполнить"],
            "relevance_score": 5,
            "hashtags": ["#News"],
        }

    async def analyze_news_item(self, news_id: str) -> bool:
        """Анализ конкретной новости"""
        logger.info(f"Analyzing news item: {news_id}")

        # Получение новости
        result = await self.db.execute(select(NewsItem).where(NewsItem.id == news_id))
        news = result.scalar_one_or_none()

        if not news:
            logger.error(f"News item not found: {news_id}")
            return False

        if news.status != NewsStatus.PENDING:
            logger.warning(
                f"News item {news_id} already analyzed, status: {news.status}"
            )
            return False

        # Создание AI задачи
        ai_task = AITask(
            news_item_id=news.id,
            task_type="analyze",
            status="processing",
            provider="openrouter",
            model=self.model,
        )
        self.db.add(ai_task)
        await self.db.commit()

        start_time = datetime.utcnow()

        try:
            # Создание промпта
            prompt = self.create_analysis_prompt(
                news.title, news.content, news.category
            )

            # AI-анализ
            analysis = await self.analyze_with_fallback(prompt)

            # Обновление новости
            news.ai_summary = analysis.get("teaser", "")
            news.ai_insights = "\n".join(
                [f"• {insight}" for insight in analysis.get("insights", [])]
            )
            news.ai_hashtags = analysis.get("hashtags", [])
            news.relevance_score = float(analysis.get("relevance_score", 0))

            # Определение статуса
            if news.relevance_score >= settings.AI_IMPORTANCE_THRESHOLD:
                news.status = NewsStatus.ANALYZED

                # Требует модерации если очень важная
                if news.relevance_score >= 8:
                    news.requires_moderation = True
            else:
                # Низкий рейтинг - отклоняем
                news.status = NewsStatus.REJECTED

            # Обновление AI задачи
            processing_time = (datetime.utcnow() - start_time).total_seconds() * 1000
            ai_task.status = "completed"
            ai_task.output_data = json.dumps(analysis, ensure_ascii=False)
            ai_task.processing_time_ms = int(processing_time)
            ai_task.completed_at = datetime.utcnow()

            await self.db.commit()

            logger.info(
                f"News analysis complete: {news_id}",
                relevance=news.relevance_score,
                status=news.status.value,
                time_ms=processing_time,
            )

            return True

        except Exception as e:
            logger.error(f"Error analyzing news {news_id}: {e}", exc_info=True)

            # Обновление AI задачи с ошибкой
            ai_task.status = "failed"
            ai_task.error_message = str(e)
            ai_task.completed_at = datetime.utcnow()

            await self.db.commit()

            return False

    async def analyze_pending_news(self, limit: int = 10) -> Dict[str, Any]:
        """Анализ всех новостей в статусе PENDING"""
        logger.info(f"Analyzing pending news (limit: {limit})")

        # Получение PENDING новостей
        result = await self.db.execute(
            select(NewsItem)
            .where(NewsItem.status == NewsStatus.PENDING)
            .order_by(NewsItem.created_at.desc())
            .limit(limit)
        )
        pending_news = result.scalars().all()

        logger.info(f"Found {len(pending_news)} pending news items")

        # Анализ каждой новости
        results = {
            "total": len(pending_news),
            "analyzed": 0,
            "rejected": 0,
            "failed": 0,
        }

        for news in pending_news:
            try:
                success = await self.analyze_news_item(str(news.id))

                if success:
                    # Обновляем статистику
                    await self.db.refresh(news)
                    if news.status == NewsStatus.ANALYZED:
                        results["analyzed"] += 1
                    elif news.status == NewsStatus.REJECTED:
                        results["rejected"] += 1
                else:
                    results["failed"] += 1

            except Exception as e:
                logger.error(f"Error processing news {news.id}: {e}")
                results["failed"] += 1

        logger.info(f"Analysis complete: {results}")

        return results


async def generate_image_for_news(news_id: str, db: AsyncSession) -> Optional[str]:
    """Генерация изображения для новости через Freepik API"""
    logger.info(f"Generating image for news: {news_id}")

    # Получение новости
    result = await db.execute(select(NewsItem).where(NewsItem.id == news_id))
    news = result.scalar_one_or_none()

    if not news:
        logger.error(f"News item not found: {news_id}")
        return None

    # Если уже есть изображение, пропускаем
    if news.image_url:
        logger.info(f"News {news_id} already has image")
        return news.image_url

    try:
        # Создание промпта для генерации изображения
        # Используем AI summary или заголовок
        image_prompt = news.ai_summary or news.title

        # Ограничиваем длину промпта
        if len(image_prompt) > 200:
            image_prompt = image_prompt[:200] + "..."

        headers = {
            "x-freepik-api-key": settings.FREEPIK_API_KEY,
            "Content-Type": "application/json",
        }

        payload = {
            "prompt": image_prompt,
            "num_images": 1,
            "image": {"size": "landscape_16_9"},
        }

        async with httpx.AsyncClient(timeout=60.0) as client:
            response = await client.post(
                f"{settings.FREEPIK_API_URL}/ai/text-to-image",
                json=payload,
                headers=headers,
            )
            response.raise_for_status()

            data = response.json()

            # Извлечение URL изображения
            if data.get("data") and len(data["data"]) > 0:
                image_url = data["data"][0].get("url")

                # Обновление новости
                news.image_url = image_url
                await db.commit()

                logger.info(f"Image generated for news {news_id}: {image_url}")
                return image_url

        return None

    except Exception as e:
        logger.error(f"Error generating image for news {news_id}: {e}", exc_info=True)
        return None
