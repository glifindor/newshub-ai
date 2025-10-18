"""
Telegram Poster Service - публикация новостей в Telegram каналы

Основные возможности:
- Автоматическая публикация новостей в каналы
- Rate limiting (20 сообщений/минуту)
- Exponential backoff при ошибках
- Уведомления админу о каждой публикации
- Поддержка Markdown форматирования
- Retry логика для Telegram API
"""

import asyncio
from collections import deque
from datetime import datetime, timedelta
from typing import Dict, List, Optional

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from telegram import Bot, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.constants import ParseMode
from telegram.error import NetworkError, RetryAfter, TimedOut
from tenacity import (
    before_sleep_log,
    retry,
    retry_if_exception_type,
    stop_after_attempt,
    wait_exponential,
)

from app.core.config import settings
from app.core.logger import get_logger
from app.models import NewsChannel, NewsItem, NewsStatus
from app.services.image_generator import ImageGenerator

logger = get_logger(__name__)


class RateLimiter:
    """Rate limiter для Telegram API (20 сообщений/минуту)"""

    def __init__(self, max_messages: int = 20, time_window: int = 60):
        self.max_messages = max_messages
        self.time_window = time_window
        self.messages = deque()

    async def acquire(self):
        """Ожидание перед отправкой следующего сообщения"""
        now = datetime.utcnow()

        # Удаляем старые записи
        while (
            self.messages
            and (now - self.messages[0]).total_seconds() > self.time_window
        ):
            self.messages.popleft()

        # Если превышен лимит, ждём
        if len(self.messages) >= self.max_messages:
            oldest = self.messages[0]
            wait_time = self.time_window - (now - oldest).total_seconds()
            if wait_time > 0:
                logger.warning(f"Rate limit reached, waiting {wait_time:.1f}s")
                await asyncio.sleep(wait_time)
                # Рекурсивно проверяем снова
                return await self.acquire()

        # Добавляем текущую отправку
        self.messages.append(now)


class TelegramPoster:
    """
    Постер новостей в Telegram каналы

    Основные функции:
    - post_news(news) - публикация одной новости
    - post_analyzed_news(limit) - публикация всех проанализированных
    - send_to_admin(news, action) - уведомление админу
    - handle_moderation_requests() - обработка запросов модерации

    Rate limiting: 20 сообщений/минуту
    Retry: 3 попытки с exponential backoff
    """

    def __init__(self, db: AsyncSession):
        self.db = db
        self.bot = Bot(token=settings.TELEGRAM_BOT_TOKEN)
        self.image_generator = ImageGenerator()

        # Rate limiter
        self.rate_limiter = RateLimiter(max_messages=20, time_window=60)

        # Маппинг каналов
        self.channel_mapping = {
            NewsChannel.CRYPTO: settings.TELEGRAM_CRYPTO_CHANNEL,
            NewsChannel.POLITICS: settings.TELEGRAM_POLITICS_CHANNEL,
        }

        # Эмодзи для категорий
        self.category_emoji = {
            NewsChannel.CRYPTO: "🔐",
            NewsChannel.POLITICS: "🏛️",
        }

        # Счётчики статистики
        self.stats = {
            "total_sent": 0,
            "total_failed": 0,
            "last_error": None,
        }

    def format_message(self, news: NewsItem, parse_mode: str = "HTML") -> str:
        """
        Форматирование сообщения для Telegram

        Формат (HTML):
        🔐 **Заголовок новости**

        📝 [Краткий тизер от AI]

        🔍 AI-инсайт:
        • Пункт 1
        • Пункт 2
        • Пункт 3

        🔗 Источник: [ссылка]

        #хэштег1 #хэштег2 #хэштег3

        Args:
            news: NewsItem объект
            parse_mode: "HTML" или "Markdown"

        Returns:
            Отформатированное сообщение
        """

        # Эмодзи для категории
        emoji = self.category_emoji.get(news.category, "📰")

        # Заголовок
        if parse_mode == "HTML":
            title = f"{emoji} <b>{news.title}</b>"
        else:
            title = f"{emoji} *{news.title}*"

        # Тизер (AI summary)
        teaser = ""
        if news.ai_summary:
            teaser = f"\n\n📝 {news.ai_summary}"

        # AI инсайты (форматирование списка с bullet points)
        insights = ""
        if news.ai_insights:
            if parse_mode == "HTML":
                insights = f"\n\n� <b>AI-инсайт:</b>"
            else:
                insights = f"\n\n🔍 *AI-инсайт:*"

            # Парсим insights (может быть строкой или списком)
            if isinstance(news.ai_insights, list):
                points = news.ai_insights
            elif isinstance(news.ai_insights, str):
                # Пытаемся разбить на пункты
                points = [p.strip() for p in news.ai_insights.split("\n") if p.strip()]
            else:
                points = []

            for point in points:
                # Убираем существующие bullet points
                point = point.lstrip("•-*").strip()
                insights += f"\n• {point}"

        # Ссылка на источник
        if parse_mode == "HTML":
            link = f"\n\n🔗 <a href='{news.url}'>Читать подробнее</a>"
        else:
            link = f"\n\n🔗 [Читать подробнее]({news.url})"

        # Хэштеги
        hashtags = ""
        if news.ai_hashtags:
            if isinstance(news.ai_hashtags, list):
                tags = news.ai_hashtags
            elif isinstance(news.ai_hashtags, str):
                # Парсим строку с хэштегами
                tags = [
                    tag.strip()
                    for tag in news.ai_hashtags.split()
                    if tag.startswith("#")
                ]
            else:
                tags = []

            if tags:
                hashtags = "\n\n" + " ".join(tags[:5])  # Максимум 5 хэштегов

        # Собираем всё вместе
        message = f"{title}{teaser}{insights}{link}{hashtags}"

        # Ограничение длины (Telegram лимит: 4096 для текста, 1024 для caption)
        max_length = 4000
        if len(message) > max_length:
            # Обрезаем тизер
            excess = len(message) - max_length
            if news.ai_summary and len(news.ai_summary) > 100:
                new_teaser_len = max(50, len(news.ai_summary) - excess - 3)
                teaser = f"\n\n📝 {news.ai_summary[:new_teaser_len]}..."
                message = f"{title}{teaser}{insights}{link}{hashtags}"

            # Если всё ещё длинно, обрезаем insights
            if len(message) > max_length:
                message = message[: max_length - 3] + "..."

        return message

    async def send_to_admin(self, news: NewsItem, action: str = "moderation") -> bool:
        """Отправить уведомление админу"""
        try:
            if action == "moderation":
                text = f"""
🚨 <b>Требуется модерация</b>

<b>Заголовок:</b> {news.title}

<b>Категория:</b> {news.category.value}
<b>Релевантность:</b> {news.relevance_score}/10

<b>AI-тизер:</b>
{news.ai_summary}

Одобрить эту новость для публикации?
"""
            else:
                text = f"ℹ️ {action}"

            # Кнопки для одобрения/отклонения
            keyboard = InlineKeyboardMarkup(
                [
                    [
                        InlineKeyboardButton(
                            "✅ Одобрить", callback_data=f"approve_{news.id}"
                        ),
                        InlineKeyboardButton(
                            "❌ Отклонить", callback_data=f"reject_{news.id}"
                        ),
                    ]
                ]
            )

            await self.bot.send_message(
                chat_id=settings.TELEGRAM_ADMIN_CHAT_ID,
                text=text,
                parse_mode=ParseMode.HTML,
                reply_markup=keyboard,
            )

            logger.info(f"Admin notification sent for news: {news.id}")
            return True

        except Exception as e:
            logger.error(f"Error sending admin notification: {e}", exc_info=True)
            return False

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=2, max=60),
        retry=retry_if_exception_type((RetryAfter, TimedOut, NetworkError)),
        before_sleep=before_sleep_log(logger, "WARNING"),
    )
    async def _send_with_retry(
        self, channel: str, message_text: str, image_url: Optional[str] = None
    ):
        """
        Отправка сообщения с retry логикой

        Retry стратегия:
        - 3 попытки максимум
        - Exponential backoff: 2s, 4s, 8s, ...
        - Retry только при RetryAfter, TimedOut, NetworkError

        Args:
            channel: ID или @username канала
            message_text: Текст сообщения
            image_url: URL изображения (опционально)

        Returns:
            telegram.Message object

        Raises:
            RetryAfter: если Telegram требует подождать
            TimedOut: если timeout
            NetworkError: если проблемы с сетью
        """
        # Rate limiting
        await self.rate_limiter.acquire()

        # Отправка с изображением или без
        if image_url:
            try:
                # Caption максимум 1024 символа
                caption = (
                    message_text[:1020] + "..."
                    if len(message_text) > 1024
                    else message_text
                )

                message = await self.bot.send_photo(
                    chat_id=channel,
                    photo=image_url,
                    caption=caption,
                    parse_mode=ParseMode.HTML,
                )
                return message

            except Exception as img_error:
                logger.warning(f"Failed to send image, fallback to text: {img_error}")
                # Fallback: отправка только текста

        # Отправка текста
        message = await self.bot.send_message(
            chat_id=channel,
            text=message_text,
            parse_mode=ParseMode.HTML,
            disable_web_page_preview=False,
        )
        return message

    async def post_news(self, news: NewsItem, notify_admin: bool = True) -> bool:
        """
        Опубликовать новость в Telegram канал

        Процесс:
        1. Проверка канала
        2. Генерация изображения (если нет)
        3. Форматирование сообщения
        4. Rate limiting
        5. Отправка с retry логикой
        6. Обновление БД
        7. Уведомление админу (если notify_admin=True)

        Args:
            news: NewsItem объект
            notify_admin: Отправить уведомление админу о публикации

        Returns:
            True если успешно, False если ошибка
        """
        logger.info(f"Posting news {news.id} to Telegram", category=news.category.value)

        try:
            # Определение канала
            channel = self.channel_mapping.get(news.category)
            if not channel:
                logger.error(f"Unknown channel for category: {news.category}")
                self.stats["total_failed"] += 1
                return False

            # Генерация изображения, если его ещё нет
            if not news.image_url and settings.FREEPIK_API_KEY:
                logger.info(f"Generating image for news {news.id}")
                try:
                    image_url = await self.image_generator.generate_for_news(
                        title=news.title,
                        summary=news.ai_summary or news.content[:200],
                        category=news.category.value,
                    )
                    if image_url:
                        news.image_url = image_url
                        await self.db.commit()
                        logger.info(f"✅ Image generated: {image_url}")
                    else:
                        logger.warning(f"⚠️ Image generation failed for news {news.id}")
                except Exception as img_error:
                    logger.error(f"Image generation error: {img_error}", exc_info=True)
                    # Продолжаем публикацию без изображения

            # Форматирование сообщения
            message_text = self.format_message(news, parse_mode="HTML")

            # Отправка с retry
            try:
                message = await self._send_with_retry(
                    channel=channel,
                    message_text=message_text,
                    image_url=news.image_url,
                )
            except RetryAfter as e:
                # Telegram просит подождать
                wait_seconds = e.retry_after
                logger.warning(f"Telegram RetryAfter: waiting {wait_seconds}s")
                await asyncio.sleep(wait_seconds)
                # Повторная попытка
                message = await self._send_with_retry(
                    channel=channel,
                    message_text=message_text,
                    image_url=news.image_url,
                )

            # Обновление новости в БД
            news.status = NewsStatus.PUBLISHED
            news.published_at = datetime.utcnow()
            news.telegram_message_id = message.message_id
            news.telegram_channel = channel

            await self.db.commit()

            # Статистика
            self.stats["total_sent"] += 1

            logger.info(
                f"✅ News posted successfully: {news.id}",
                channel=channel,
                message_id=message.message_id,
                title=news.title[:50],
            )

            # Уведомление админу
            if notify_admin:
                await self._notify_admin_about_post(news, channel, message.message_id)

            return True

        except Exception as e:
            logger.error(f"❌ Error posting news {news.id}: {e}", exc_info=True)
            self.stats["total_failed"] += 1
            self.stats["last_error"] = str(e)
            return False

    async def _notify_admin_about_post(
        self, news: NewsItem, channel: str, message_id: int
    ):
        """Уведомить админа о публикации новости"""
        try:
            notification_text = f"""
✅ <b>Новость опубликована</b>

<b>Канал:</b> {channel}
<b>Заголовок:</b> {news.title}

<b>Релевантность:</b> {news.relevance_score}/10
<b>Message ID:</b> {message_id}

<a href="https://t.me/{channel.lstrip('@')}/{message_id}">Перейти к посту</a>
"""

            await self.bot.send_message(
                chat_id=settings.TELEGRAM_ADMIN_CHAT_ID,
                text=notification_text,
                parse_mode=ParseMode.HTML,
                disable_web_page_preview=True,
            )

        except Exception as e:
            logger.error(f"Failed to notify admin about post: {e}")

    async def post_analyzed_news(self, limit: int = 5) -> dict:
        """Опубликовать все проанализированные новости"""
        logger.info(f"Posting analyzed news (limit: {limit})")

        # Получение ANALYZED новостей без требования модерации
        result = await self.db.execute(
            select(NewsItem)
            .where(
                NewsItem.status == NewsStatus.ANALYZED,
                NewsItem.requires_moderation == False,
            )
            .order_by(NewsItem.relevance_score.desc())
            .limit(limit)
        )
        news_items = result.scalars().all()

        logger.info(f"Found {len(news_items)} news items to post")

        results = {
            "total": len(news_items),
            "posted": 0,
            "failed": 0,
        }

        for news in news_items:
            try:
                success = await self.post_news(news)

                if success:
                    results["posted"] += 1
                    # Задержка между постами (антиспам)
                    await asyncio.sleep(2)
                else:
                    results["failed"] += 1

            except Exception as e:
                logger.error(f"Error processing news {news.id}: {e}")
                results["failed"] += 1

        logger.info(f"Posting complete: {results}")

        return results

    async def handle_moderation_requests(self) -> dict:
        """Обработать новости, требующие модерации"""
        logger.info("Checking for news requiring moderation")

        # Получение ANALYZED новостей с флагом модерации
        result = await self.db.execute(
            select(NewsItem)
            .where(
                NewsItem.status == NewsStatus.ANALYZED,
                NewsItem.requires_moderation == True,
            )
            .order_by(NewsItem.relevance_score.desc())
            .limit(10)
        )
        news_items = result.scalars().all()

        logger.info(f"Found {len(news_items)} news items requiring moderation")

        results = {
            "total": len(news_items),
            "notified": 0,
            "failed": 0,
        }

        for news in news_items:
            try:
                success = await self.send_to_admin(news, action="moderation")

                if success:
                    results["notified"] += 1
                    # Сбрасываем флаг, чтобы не отправлять повторно
                    news.requires_moderation = False
                    await self.db.commit()
                else:
                    results["failed"] += 1

            except Exception as e:
                logger.error(f"Error notifying admin about news {news.id}: {e}")
                results["failed"] += 1

        return results


async def approve_news_for_posting(news_id: str, db: AsyncSession) -> bool:
    """Одобрить новость для публикации (вызывается админом)"""
    logger.info(f"Approving news for posting: {news_id}")

    # Получение новости
    result = await db.execute(select(NewsItem).where(NewsItem.id == news_id))
    news = result.scalar_one_or_none()

    if not news:
        logger.error(f"News item not found: {news_id}")
        return False

    if news.status != NewsStatus.ANALYZED:
        logger.warning(f"News {news_id} cannot be approved, status: {news.status}")
        return False

    # Изменение статуса на APPROVED (будет опубликовано автоматически)
    news.status = NewsStatus.APPROVED
    await db.commit()

    # Публикация
    poster = TelegramPoster(db)
    success = await poster.post_news(news)

    return success


async def reject_news(news_id: str, db: AsyncSession) -> bool:
    """Отклонить новость"""
    logger.info(f"Rejecting news: {news_id}")

    result = await db.execute(select(NewsItem).where(NewsItem.id == news_id))
    news = result.scalar_one_or_none()

    if not news:
        return False

    news.status = NewsStatus.REJECTED
    await db.commit()

    return True
