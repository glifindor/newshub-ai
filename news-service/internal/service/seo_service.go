package service

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"go.uber.org/zap"
)

// SEOService интерфейс для работы с SEO оптимизацией
type SEOService interface {
	// GenerateMetaTags генерирует все SEO метатеги за один вызов
	GenerateMetaTags(ctx context.Context, title, content string) (seoTitle, seoDesc, keywords string)

	// OptimizeTitle оптимизирует заголовок для SEO (макс 60 символов)
	OptimizeTitle(title string) string

	// GenerateDescription создает SEO описание из контента (макс 160 символов)
	GenerateDescription(content string) string

	// ExtractKeywords извлекает ключевые слова из текста
	ExtractKeywords(text string, count int) []string
}

type seoService struct {
	logger *zap.Logger
}

// NewSEOService создает новый экземпляр SEO сервиса
func NewSEOService(logger *zap.Logger) SEOService {
	return &seoService{
		logger: logger,
	}
}

// GenerateMetaTags генерирует все SEO метатеги одновременно
func (s *seoService) GenerateMetaTags(ctx context.Context, title, content string) (string, string, string) {
	seoTitle := s.OptimizeTitle(title)
	seoDesc := s.GenerateDescription(content)
	keywords := s.ExtractKeywords(title+" "+content, 5)

	s.logger.Debug("Generated SEO meta tags",
		zap.String("original_title", title),
		zap.String("seo_title", seoTitle),
		zap.Int("desc_length", len(seoDesc)),
		zap.Strings("keywords", keywords),
	)

	return seoTitle, seoDesc, strings.Join(keywords, ", ")
}

// OptimizeTitle оптимизирует заголовок для SEO (рекомендуется макс 60 символов)
func (s *seoService) OptimizeTitle(title string) string {
	title = strings.TrimSpace(title)

	// Если заголовок уже оптимального размера
	if len(title) <= 60 {
		return title
	}

	// Обрезаем до 57 символов и добавляем многоточие
	optimized := title[:57] + "..."

	// Пытаемся обрезать по последней границе слова
	lastSpace := strings.LastIndex(optimized[:57], " ")
	if lastSpace > 30 { // Сохраняем минимум 30 символов
		optimized = title[:lastSpace] + "..."
	}

	return optimized
}

// GenerateDescription создает SEO описание из контента (макс 160 символов)
func (s *seoService) GenerateDescription(content string) string {
	// Удаляем HTML теги
	plainText := stripHTML(content)

	// Нормализуем пробелы
	plainText = regexp.MustCompile(`\s+`).ReplaceAllString(plainText, " ")
	plainText = strings.TrimSpace(plainText)

	if len(plainText) == 0 {
		return ""
	}

	if len(plainText) <= 160 {
		return plainText
	}

	// Обрезаем до 157 символов, оставляем место для "..."
	desc := plainText[:157]

	// Пытаемся обрезать по концу предложения
	if idx := strings.LastIndexAny(desc, ".!?"); idx > 100 {
		return strings.TrimSpace(plainText[:idx+1])
	}

	// Иначе обрезаем по последней границе слова
	lastSpace := strings.LastIndex(desc, " ")
	if lastSpace > 100 {
		desc = plainText[:lastSpace]
	}

	return desc + "..."
}

// ExtractKeywords извлекает наиболее важные ключевые слова из текста
func (s *seoService) ExtractKeywords(text string, count int) []string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	// Удаляем HTML и пунктуацию
	text = stripHTML(text)
	text = regexp.MustCompile(`[^\p{L}\p{N}\s]`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Разбиваем на слова
	words := strings.Fields(text)

	// Определяем стоп-слова (базовый набор для английского и русского)
	stopWords := map[string]bool{
		// English
		"the": true, "be": true, "to": true, "of": true, "and": true,
		"a": true, "in": true, "that": true, "have": true, "i": true,
		"it": true, "for": true, "not": true, "on": true, "with": true,
		"he": true, "as": true, "you": true, "do": true, "at": true,
		"this": true, "but": true, "his": true, "by": true, "from": true,
		"they": true, "we": true, "say": true, "her": true, "she": true,
		"or": true, "an": true, "will": true, "my": true, "one": true,
		"all": true, "would": true, "there": true, "their": true,
		"is": true, "are": true, "was": true, "were": true, "been": true,
		"has": true, "had": true, "can": true, "could": true, "may": true,
		// Russian
		"и": true, "в": true, "не": true, "на": true, "с": true,
		"по": true, "из": true, "к": true, "за": true, "у": true,
		"о": true, "от": true, "до": true, "для": true, "при": true,
		"это": true, "как": true, "его": true, "все": true, "что": true,
		"так": true, "же": true, "или": true, "был": true, "были": true,
		"быть": true, "есть": true, "чтобы": true, "можно": true, "только": true,
		"без": true, "где": true, "уже": true, "еще": true, "ни": true,
		"когда": true, "тот": true, "этот": true, "через": true, "под": true,
	}

	// Подсчитываем частоту слов (исключая стоп-слова и короткие слова)
	frequency := make(map[string]int)
	for _, word := range words {
		if len(word) > 3 && !stopWords[word] {
			frequency[word]++
		}
	}

	// Сортируем по частоте
	type wordFreq struct {
		word  string
		count int
	}

	var sorted []wordFreq
	for word, count := range frequency {
		sorted = append(sorted, wordFreq{word, count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	// Возвращаем топ N ключевых слов
	keywords := make([]string, 0, count)
	for i := 0; i < count && i < len(sorted); i++ {
		keywords = append(keywords, sorted[i].word)
	}

	return keywords
}

// stripHTML удаляет HTML теги из текста
func stripHTML(html string) string {
	// Удаляем теги script и style с содержимым
	re := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	html = re.ReplaceAllString(html, "")
	re = regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	html = re.ReplaceAllString(html, "")

	// Удаляем все остальные HTML теги
	re = regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(html, "")
}
