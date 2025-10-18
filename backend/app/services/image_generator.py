"""
Image Generator Service - генерация изображений для новостей через Freepik AI Mystic API

Документация: https://developers.freepik.com/reference/ai-endpoints-image-generation
"""

import asyncio
import time
from typing import Optional

import httpx
from tenacity import (retry, retry_if_exception_type, stop_after_attempt,
                      wait_exponential)

from app.core.config import settings
from app.core.logger import get_logger

logger = get_logger(__name__)


class ImageGenerator:
    """
    Генератор изображений для новостей через Freepik AI Mystic API

    Features:
    - Генерация изображений по текстовому промпту
    - Поддержка разных разрешений (2k по умолчанию)
    - Retry логика с exponential backoff
    - Polling для получения результата
    """

    def __init__(self):
        self.api_key = settings.FREEPIK_API_KEY
        self.api_url = "https://api.freepik.com/v1/ai/mystic"
        self.timeout = 60.0

        if not self.api_key:
            logger.warning("FREEPIK_API_KEY not set, image generation disabled")

    @retry(
        retry=retry_if_exception_type((httpx.HTTPStatusError, httpx.RequestError)),
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=2, max=10),
    )
    async def generate_image(
        self,
        prompt: str,
        title: str = "",
        category: str = "general",
        resolution: str = "2k",
        aspect_ratio: str = "landscape_4_3",
        model: str = "realism",
    ) -> Optional[str]:
        """
        Генерирует изображение по промпту

        Args:
            prompt: Текстовое описание изображения
            title: Заголовок новости (для улучшения промпта)
            category: Категория новости (crypto, politics)
            resolution: Разрешение (2k, 4k, 8k)
            aspect_ratio: Соотношение сторон (landscape_4_3, square_1_1, portrait_3_4)
            model: Модель (realism, artistic, anime)

        Returns:
            URL сгенерированного изображения или None
        """
        if not self.api_key:
            logger.warning("Image generation skipped: no API key")
            return None

        # Улучшаем промпт с учетом категории
        enhanced_prompt = self._enhance_prompt(prompt, title, category)

        logger.info(f"Generating image with prompt: {enhanced_prompt[:100]}...")

        # Payload для генерации
        payload = {
            "prompt": enhanced_prompt,
            "resolution": resolution,
            "aspect_ratio": aspect_ratio,
            "model": model,
            "creative_detailing": 50,
            "engine": "automatic",
            "fixed_generation": False,
            "filter_nsfw": True,  # Фильтруем NSFW контент
        }

        headers = {
            "x-freepik-api-key": self.api_key,
            "Content-Type": "application/json",
        }

        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                # 1. Запускаем генерацию
                response = await client.post(
                    self.api_url, json=payload, headers=headers
                )
                response.raise_for_status()
                data = response.json()

                if "data" not in data:
                    logger.error(f"Invalid response from Freepik API: {data}")
                    return None

                # 2. Получаем ID задачи
                generation_id = data["data"]["id"]
                logger.info(f"Image generation started: {generation_id}")

                # 3. Polling для получения результата
                image_url = await self._poll_generation_result(
                    generation_id, headers, max_attempts=30, delay=2
                )

                if image_url:
                    logger.info(f"Image generated successfully: {image_url}")
                    return image_url
                else:
                    logger.warning(
                        f"Image generation failed for: {enhanced_prompt[:50]}"
                    )
                    return None

        except httpx.HTTPStatusError as e:
            logger.error(
                f"Freepik API error: {e.response.status_code} - {e.response.text}"
            )
            raise
        except httpx.RequestError as e:
            logger.error(f"Request error: {str(e)}")
            raise
        except Exception as e:
            logger.error(f"Unexpected error generating image: {str(e)}")
            return None

    async def _poll_generation_result(
        self,
        generation_id: str,
        headers: dict,
        max_attempts: int = 30,
        delay: int = 2,
    ) -> Optional[str]:
        """
        Polling для получения результата генерации

        Args:
            generation_id: ID задачи генерации
            headers: HTTP заголовки с API ключом
            max_attempts: Максимальное количество попыток
            delay: Задержка между попытками (секунды)

        Returns:
            URL изображения или None
        """
        status_url = f"{self.api_url}/{generation_id}"

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            for attempt in range(max_attempts):
                try:
                    response = await client.get(status_url, headers=headers)
                    response.raise_for_status()
                    data = response.json()

                    if "data" not in data:
                        logger.error(f"Invalid status response: {data}")
                        return None

                    status = data["data"].get("status")

                    if status == "completed":
                        # Генерация завершена успешно
                        images = data["data"].get("images", [])
                        if images and len(images) > 0:
                            return images[0].get("url")
                        else:
                            logger.error("No images in completed generation")
                            return None

                    elif status == "failed":
                        # Генерация провалилась
                        error = data["data"].get("error", "Unknown error")
                        logger.error(f"Image generation failed: {error}")
                        return None

                    elif status in ["pending", "processing"]:
                        # Ещё генерируется, ждём
                        logger.debug(
                            f"Generation in progress ({attempt + 1}/{max_attempts}): {status}"
                        )
                        await asyncio.sleep(delay)

                    else:
                        logger.warning(f"Unknown generation status: {status}")
                        await asyncio.sleep(delay)

                except httpx.HTTPStatusError as e:
                    logger.error(f"Error polling generation: {e.response.status_code}")
                    await asyncio.sleep(delay)
                except Exception as e:
                    logger.error(f"Error polling generation: {str(e)}")
                    await asyncio.sleep(delay)

        logger.warning(f"Generation polling timeout after {max_attempts * delay}s")
        return None

    def _enhance_prompt(self, base_prompt: str, title: str, category: str) -> str:
        """
        Улучшает промпт с учетом контекста новости

        Args:
            base_prompt: Базовый промпт
            title: Заголовок новости
            category: Категория (crypto, politics)

        Returns:
            Улучшенный промпт
        """
        # Префиксы для разных категорий
        category_prefixes = {
            "crypto": "Professional news image about cryptocurrency and blockchain technology: ",
            "politics": "Professional news image about politics and government: ",
            "general": "Professional news image: ",
        }

        prefix = category_prefixes.get(category.lower(), category_prefixes["general"])

        # Добавляем контекст из заголовка если промпт короткий
        if len(base_prompt) < 50 and title:
            enhanced = f"{prefix}{title}. {base_prompt}"
        else:
            enhanced = f"{prefix}{base_prompt}"

        # Добавляем общие требования к качеству
        enhanced += (
            ". High quality, professional photography, 4k, detailed, "
            "modern design, clean composition, suitable for news article"
        )

        # Ограничиваем длину (Freepik API обычно имеет лимит)
        if len(enhanced) > 1000:
            enhanced = enhanced[:997] + "..."

        return enhanced

    async def generate_for_news(
        self, title: str, summary: str, category: str
    ) -> Optional[str]:
        """
        Генерирует изображение специально для новости

        Args:
            title: Заголовок новости
            summary: Краткое описание
            category: Категория новости

        Returns:
            URL изображения или None
        """
        # Создаём промпт из заголовка и саммари
        prompt = f"{title}. {summary[:200]}"

        # Выбираем aspect ratio в зависимости от категории
        aspect_ratio = "landscape_4_3"  # Хорошо для новостей
        if category == "crypto":
            model = "realism"  # Реалистичный стиль для крипто-новостей
        else:
            model = "realism"  # Реалистичный для всех

        return await self.generate_image(
            prompt=prompt,
            title=title,
            category=category,
            aspect_ratio=aspect_ratio,
            model=model,
        )
