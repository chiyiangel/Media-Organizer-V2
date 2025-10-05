package i18n

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Language 支持的语言类型
type Language string

const (
	LanguageChinese Language = "zh" // 中文
	LanguageEnglish Language = "en" // 英文
)

// Localizer 本地化管理器
type Localizer struct {
	currentLang Language
	messages    map[Language]map[string]string
	mu          sync.RWMutex
}

// 全局本地化实例
var (
	globalLocalizer *Localizer
	once            sync.Once
)

// GetLocalizer 获取全局本地化实例
func GetLocalizer() *Localizer {
	once.Do(func() {
		globalLocalizer = NewLocalizer()
	})
	return globalLocalizer
}

// NewLocalizer 创建新的本地化管理器
func NewLocalizer() *Localizer {
	l := &Localizer{
		currentLang: LanguageEnglish, // 默认英文
		messages:    make(map[Language]map[string]string),
	}

	// 初始化语言资源
	l.initMessages()

	// 检测系统语言
	l.detectSystemLanguage()

	return l
}

// detectSystemLanguage 检测系统语言
func (l *Localizer) detectSystemLanguage() {
	// 检测顺序：LANG环境变量 -> LC_ALL -> LC_MESSAGES -> Windows语言设置

	// Unix系统环境变量
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		lang = os.Getenv("LC_MESSAGES")
	}

	// Windows系统可能需要额外检测
	if lang == "" {
		// Windows下通常LANG为空，可以通过其他方式检测
		lang = os.Getenv("LANGUAGE")
	}

	// 解析语言代码
	detectedLang := l.parseLanguageCode(lang)
	if detectedLang != "" {
		l.SetLanguage(Language(detectedLang))
	}
}

// parseLanguageCode 解析语言代码
func (l *Localizer) parseLanguageCode(langCode string) string {
	if langCode == "" {
		return string(LanguageEnglish) // 默认返回英文
	}

	// 转换为小写并处理常见格式
	langCode = strings.ToLower(langCode)

	// 提取主要语言代码 (例如 zh_CN -> zh, en_US -> en)
	parts := strings.FieldsFunc(langCode, func(c rune) bool {
		return c == '_' || c == '-' || c == '.'
	})

	if len(parts) > 0 {
		mainLang := parts[0]

		// 匹配支持的语言
		switch {
		case strings.HasPrefix(mainLang, "zh"): // 中文 (zh, zh_CN, zh_TW等)
			return string(LanguageChinese)
		case strings.HasPrefix(mainLang, "en"): // 英文 (en, en_US, en_GB等)
			return string(LanguageEnglish)
		}
	}

	// 不支持的语言默认返回英文
	return string(LanguageEnglish)
} // SetLanguage 设置当前语言
func (l *Localizer) SetLanguage(lang Language) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查语言是否支持
	if _, exists := l.messages[lang]; exists {
		l.currentLang = lang
	}
}

// GetLanguage 获取当前语言
func (l *Localizer) GetLanguage() Language {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.currentLang
}

// T 翻译文本 (Translate的简写)
func (l *Localizer) T(key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 获取当前语言的消息映射
	if langMessages, exists := l.messages[l.currentLang]; exists {
		if message, found := langMessages[key]; found {
			return message
		}
	}

	// 如果当前语言没有找到，尝试英文作为回退
	if l.currentLang != LanguageEnglish {
		if langMessages, exists := l.messages[LanguageEnglish]; exists {
			if message, found := langMessages[key]; found {
				return message
			}
		}
	}

	// 如果都没有找到，返回键本身
	return key
}

// Tf 带格式化的翻译 (Translate with format)
func (l *Localizer) Tf(key string, args ...interface{}) string {
	template := l.T(key)
	if len(args) > 0 {
		// 简单的字符串替换，你也可以使用更复杂的格式化
		result := template
		for i, arg := range args {
			placeholder := "{" + fmt.Sprintf("%d", i) + "}"
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", arg))
		}
		return result
	}
	return template
}

// 便捷函数：直接使用全局实例
func T(key string) string {
	return GetLocalizer().T(key)
}

func Tf(key string, args ...interface{}) string {
	return GetLocalizer().Tf(key, args...)
}

func SetLanguage(lang Language) {
	GetLocalizer().SetLanguage(lang)
}

func GetLanguage() Language {
	return GetLocalizer().GetLanguage()
}
