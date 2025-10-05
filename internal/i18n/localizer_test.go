package i18n

import (
	"os"
	"testing"
)

func TestLanguageDetection(t *testing.T) {
	l := NewLocalizer()

	// 测试中文语言检测
	testCases := []struct {
		langCode string
		expected Language
	}{
		{"zh_CN", LanguageChinese},
		{"zh_TW", LanguageChinese},
		{"zh", LanguageChinese},
		{"en_US", LanguageEnglish},
		{"en_GB", LanguageEnglish},
		{"en", LanguageEnglish},
		{"fr_FR", LanguageEnglish}, // 不支持的语言应返回默认英文
		{"", LanguageEnglish},      // 空字符串应返回默认英文
	}

	for _, tc := range testCases {
		result := l.parseLanguageCode(tc.langCode)
		if result != string(tc.expected) {
			t.Errorf("parseLanguageCode(%s) = %s, want %s", tc.langCode, result, tc.expected)
		}
	}
}

func TestTranslation(t *testing.T) {
	l := NewLocalizer()

	// 测试中文翻译
	l.SetLanguage(LanguageChinese)
	if l.T("app.title") != "📸 照片视频整理工具 v1.0" {
		t.Errorf("Chinese translation failed")
	}

	// 测试英文翻译
	l.SetLanguage(LanguageEnglish)
	if l.T("app.title") != "📸 Photo Video Organizer v1.0" {
		t.Errorf("English translation failed")
	}

	// 测试不存在的键应返回键本身
	if l.T("nonexistent.key") != "nonexistent.key" {
		t.Errorf("Should return key when translation not found")
	}
}

func TestFormatTranslation(t *testing.T) {
	l := NewLocalizer()
	l.SetLanguage(LanguageChinese)

	// 测试带参数的翻译
	result := l.Tf("error.extract_date", "test error")
	expected := "无法提取日期: test error"
	if result != expected {
		t.Errorf("Tf() = %s, want %s", result, expected)
	}
}

func TestSystemLanguageDetection(t *testing.T) {
	// 保存原始环境变量
	origLang := os.Getenv("LANG")
	defer os.Setenv("LANG", origLang)

	// 测试中文环境
	os.Setenv("LANG", "zh_CN.UTF-8")
	l := NewLocalizer()
	if l.GetLanguage() != LanguageChinese {
		t.Errorf("Should detect Chinese from LANG=zh_CN.UTF-8")
	}

	// 测试英文环境
	os.Setenv("LANG", "en_US.UTF-8")
	l = NewLocalizer()
	if l.GetLanguage() != LanguageEnglish {
		t.Errorf("Should detect English from LANG=en_US.UTF-8")
	}
}
