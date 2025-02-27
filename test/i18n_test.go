package test

import (
	i18 "ipfast_server/internal/config/i18n"
	"log"
	"testing"
)

// 测试动态添加翻译内容
func TestAddTranslation(t *testing.T) {
	translation := i18.Translation{
		Lang:           "en",
		MessageID:      "test_message",
		TranslationStr: "测试消息",
	}
	i18.AddTranslation(translation, nil)
	// 检查翻译是否已添加
	localizer := i18.NewLocalizer("en")
	result := localizer.F("test_message")
	log.Println(result)
	if result != "测试消息" {
		t.Errorf("Expected '测试消息', but got '%s'", result)
	}
}

// 测试动态删除翻译内容
func TestRemoveTranslation(t *testing.T) {
	translation1 := i18.Translation{
		Lang:           "en",
		MessageID:      "test_message",
		TranslationStr: "测试消息",
	}
	i18.AddTranslation(translation1, nil)
	i18.RemoveTranslation(translation1)
	i18.SyncTranslations()

	// 检查翻译是否已删除
	localizer := i18.NewLocalizer("en")
	result := localizer.F("test_message")
	if result == "测试消息" {
		t.Errorf("预期消息ID翻译为 'test_message', but got '%s'", result)
	}
}

// 测试查询当前已添加的翻译内容
func TestGetTranslation(t *testing.T) {
	translation := i18.Translation{
		Lang:           "en",
		MessageID:      "test_message",
		TranslationStr: "测试消息",
	}
	i18.AddTranslation(translation, nil)
	translations := i18.GetTranslation()
	if len(translations) != 1 {
		t.Errorf("预期只有一条已经添加的翻译, 但是实际为 %d", len(translations))
	}

	if translations[0].MessageID != "test_message" {
		t.Errorf("预期消息ID为 'test_message', 但是实际为 '%s'", translations[0].MessageID)
	}
}
