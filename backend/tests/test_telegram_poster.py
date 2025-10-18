"""
Unit тесты для TelegramPoster

Тестируем:
- Форматирование сообщений
- Rate limiting
- Retry логику
- Обработку ошибок
"""
import pytest
import asyncio
from datetime import datetime, timedelta
from unittest.mock import Mock, AsyncMock, patch
from uuid import uuid4

from app.services.telegram_poster import TelegramPoster, RateLimiter
from app.models import NewsItem, NewsStatus, NewsChannel


@pytest.fixture
def mock_db():
    """Mock AsyncSession"""
    db = AsyncMock()
    db.commit = AsyncMock()
    db.execute = AsyncMock()
    return db


@pytest.fixture
def mock_bot():
    """Mock Telegram Bot"""
    bot = AsyncMock()
    bot.send_message = AsyncMock()
    bot.send_photo = AsyncMock()
    bot.get_me = AsyncMock(return_value=Mock(
        id=123456789,
        first_name="NewsHub AI",
        username="NewsHubBot"
    ))
    return bot


@pytest.fixture
def sample_news():
    """Пример новости для тестов"""
    return NewsItem(
        id=uuid4(),
        title="Bitcoin достиг $100,000",
        content="Полный текст новости...",
        url="https://example.com/bitcoin-100k",
        category=NewsChannel.CRYPTO,
        status=NewsStatus.ANALYZED,
        ai_summary="Краткий тизер: Bitcoin впервые преодолел $100k благодаря институциональным инвестициям.",
        ai_insights=["Рост на 300%", "ETF стимулирует рынок", "Прогноз $150k"],
        ai_hashtags=["#Bitcoin", "#Crypto", "#ATH"],
        relevance_score=9.5,
        created_at=datetime.utcnow(),
    )


class TestRateLimiter:
    """Тесты для RateLimiter"""
    
    @pytest.mark.asyncio
    async def test_rate_limiter_basic(self):
        """Тест базовой функциональности rate limiter"""
        limiter = RateLimiter(max_messages=5, time_window=10)
        
        # Отправляем 5 сообщений - должно пройти быстро
        start = datetime.utcnow()
        for _ in range(5):
            await limiter.acquire()
        elapsed = (datetime.utcnow() - start).total_seconds()
        
        assert elapsed < 1  # Должно занять меньше секунды
        assert len(limiter.messages) == 5
    
    @pytest.mark.asyncio
    async def test_rate_limiter_blocking(self):
        """Тест блокировки при превышении лимита"""
        limiter = RateLimiter(max_messages=3, time_window=2)
        
        # Отправляем 3 сообщения
        for _ in range(3):
            await limiter.acquire()
        
        # 4-е сообщение должно заблокироваться
        start = datetime.utcnow()
        await limiter.acquire()
        elapsed = (datetime.utcnow() - start).total_seconds()
        
        assert elapsed >= 1.5  # Должно подождать минимум 1.5 секунды
    
    @pytest.mark.asyncio
    async def test_rate_limiter_cleanup(self):
        """Тест очистки старых записей"""
        limiter = RateLimiter(max_messages=5, time_window=1)
        
        # Отправляем 3 сообщения
        for _ in range(3):
            await limiter.acquire()
        
        # Ждём 1.5 секунды
        await asyncio.sleep(1.5)
        
        # Старые записи должны удалиться
        assert len(limiter.messages) == 0


class TestTelegramPoster:
    """Тесты для TelegramPoster"""
    
    def test_format_message_html(self, sample_news):
        """Тест форматирования сообщения в HTML"""
        poster = TelegramPoster(db=Mock())
        message = poster.format_message(sample_news, parse_mode="HTML")
        
        # Проверяем наличие основных элементов
        assert "🔐" in message  # Эмодзи
        assert "<b>Bitcoin достиг $100,000</b>" in message  # Заголовок
        assert "📝" in message  # Тизер
        assert "🔍" in message  # AI-инсайт
        assert "• Рост на 300%" in message  # Первый инсайт
        assert "🔗" in message  # Ссылка
        assert "#Bitcoin" in message  # Хэштеги
    
    def test_format_message_length_limit(self, sample_news):
        """Тест ограничения длины сообщения"""
        # Создаём очень длинную новость
        sample_news.ai_summary = "Очень длинный текст " * 500
        sample_news.ai_insights = ["Инсайт " + "x" * 1000] * 10
        
        poster = TelegramPoster(db=Mock())
        message = poster.format_message(sample_news)
        
        # Должна обрезаться до 4000 символов
        assert len(message) <= 4000
        assert message.endswith("...")
    
    def test_format_message_with_list_hashtags(self, sample_news):
        """Тест форматирования с хэштегами в виде списка"""
        sample_news.ai_hashtags = ["#Bitcoin", "#Crypto", "#ATH", "#BTC", "#Blockchain"]
        
        poster = TelegramPoster(db=Mock())
        message = poster.format_message(sample_news)
        
        # Должно быть максимум 5 хэштегов
        hashtag_count = message.count("#")
        assert hashtag_count <= 5
    
    def test_format_message_with_string_insights(self, sample_news):
        """Тест форматирования insights из строки"""
        sample_news.ai_insights = "Пункт 1\nПункт 2\nПункт 3"
        
        poster = TelegramPoster(db=Mock())
        message = poster.format_message(sample_news)
        
        # Должны быть bullet points
        assert "• Пункт 1" in message
        assert "• Пункт 2" in message
        assert "• Пункт 3" in message
    
    @pytest.mark.asyncio
    async def test_post_news_success(self, mock_db, mock_bot, sample_news):
        """Тест успешной публикации новости"""
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        # Mock send_message
        mock_message = Mock()
        mock_message.message_id = 12345
        mock_bot.send_message.return_value = mock_message
        
        # Публикуем
        result = await poster.post_news(sample_news, notify_admin=False)
        
        assert result is True
        assert sample_news.status == NewsStatus.PUBLISHED
        assert sample_news.telegram_message_id == 12345
        mock_bot.send_message.assert_called_once()
        mock_db.commit.assert_called_once()
    
    @pytest.mark.asyncio
    async def test_post_news_with_image(self, mock_db, mock_bot, sample_news):
        """Тест публикации новости с изображением"""
        sample_news.image_url = "https://example.com/image.jpg"
        
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        mock_message = Mock()
        mock_message.message_id = 12345
        mock_bot.send_photo.return_value = mock_message
        
        result = await poster.post_news(sample_news, notify_admin=False)
        
        assert result is True
        mock_bot.send_photo.assert_called_once()
        
        # Проверяем параметры вызова
        call_args = mock_bot.send_photo.call_args
        assert call_args.kwargs['photo'] == sample_news.image_url
        assert 'caption' in call_args.kwargs
    
    @pytest.mark.asyncio
    async def test_post_news_image_fallback(self, mock_db, mock_bot, sample_news):
        """Тест fallback на текст при ошибке загрузки изображения"""
        sample_news.image_url = "https://example.com/broken-image.jpg"
        
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        # send_photo падает, send_message работает
        mock_bot.send_photo.side_effect = Exception("Image load error")
        mock_message = Mock()
        mock_message.message_id = 12345
        mock_bot.send_message.return_value = mock_message
        
        result = await poster.post_news(sample_news, notify_admin=False)
        
        assert result is True
        mock_bot.send_photo.assert_called_once()
        mock_bot.send_message.assert_called_once()  # Fallback
    
    @pytest.mark.asyncio
    async def test_post_news_unknown_channel(self, mock_db, mock_bot, sample_news):
        """Тест обработки неизвестного канала"""
        sample_news.category = None  # Неизвестная категория
        
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        result = await poster.post_news(sample_news, notify_admin=False)
        
        assert result is False
        mock_bot.send_message.assert_not_called()
    
    @pytest.mark.asyncio
    async def test_post_news_telegram_error(self, mock_db, mock_bot, sample_news):
        """Тест обработки ошибки Telegram API"""
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        # Telegram API падает
        mock_bot.send_message.side_effect = Exception("Telegram API error")
        
        result = await poster.post_news(sample_news, notify_admin=False)
        
        assert result is False
        assert poster.stats['total_failed'] == 1
        assert "Telegram API error" in poster.stats['last_error']
    
    @pytest.mark.asyncio
    async def test_post_analyzed_news(self, mock_db, mock_bot):
        """Тест публикации всех проанализированных новостей"""
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        # Mock БД запрос
        mock_result = Mock()
        news_items = [
            NewsItem(
                id=uuid4(),
                title=f"News {i}",
                content="Content",
                url="https://example.com",
                category=NewsChannel.CRYPTO,
                status=NewsStatus.ANALYZED,
                requires_moderation=False,
            )
            for i in range(3)
        ]
        mock_result.scalars.return_value.all.return_value = news_items
        mock_db.execute.return_value = mock_result
        
        # Mock send_message
        mock_message = Mock()
        mock_message.message_id = 12345
        mock_bot.send_message.return_value = mock_message
        
        # Публикуем
        result = await poster.post_analyzed_news(limit=5)
        
        assert result['total'] == 3
        assert result['posted'] == 3
        assert result['failed'] == 0
        assert mock_bot.send_message.call_count == 3
    
    @pytest.mark.asyncio
    async def test_send_to_admin_moderation(self, mock_db, mock_bot, sample_news):
        """Тест отправки запроса на модерацию админу"""
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        mock_bot.send_message.return_value = Mock()
        
        result = await poster.send_to_admin(sample_news, action="moderation")
        
        assert result is True
        mock_bot.send_message.assert_called_once()
        
        # Проверяем параметры
        call_args = mock_bot.send_message.call_args
        assert call_args.kwargs['chat_id'] == poster.bot.send_message.call_args.kwargs['chat_id']
        assert "Требуется модерация" in call_args.kwargs['text']
        assert 'reply_markup' in call_args.kwargs  # Кнопки одобрить/отклонить


class TestIntegration:
    """Интеграционные тесты"""
    
    @pytest.mark.asyncio
    @pytest.mark.integration
    async def test_full_posting_flow(self, mock_db, mock_bot, sample_news):
        """Тест полного потока публикации"""
        poster = TelegramPoster(db=mock_db)
        poster.bot = mock_bot
        
        # Mock ответы
        mock_message = Mock()
        mock_message.message_id = 12345
        mock_bot.send_message.return_value = mock_message
        
        # 1. Публикация
        result = await poster.post_news(sample_news, notify_admin=True)
        assert result is True
        
        # 2. Проверка статуса в БД
        assert sample_news.status == NewsStatus.PUBLISHED
        assert sample_news.telegram_message_id == 12345
        
        # 3. Проверка уведомления админу
        assert mock_bot.send_message.call_count == 2  # Пост + уведомление админу
        
        # 4. Проверка статистики
        assert poster.stats['total_sent'] == 1
        assert poster.stats['total_failed'] == 0


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
