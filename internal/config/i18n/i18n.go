package i18n

import (
	"encoding/json"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type MyLocalizer struct {
	*i18n.Localizer
}

const (
	currentBundlePath = "locales/admin/"
	defaultBundlePath = "locales/"
)

var languageFileMap = map[string]string{
	"en":    "en.json",
	"zh_CN": "zh_CN.json",
	"zh_TW": "zh_TW.json",
}
var (
	// 动态翻译内容缓存
	translationCache sync.Map

	// 当前使用的翻译 Bundle 在默认的翻译 Bundle 基础上添加了动态翻译内容
	currentBundle *i18n.Bundle

	// 默认的翻译 Bundle 作为备份
	defaultBundle *i18n.Bundle
	// 是否需要重新加载当前 Bundle
	needsReload bool
	// 互斥锁
	mu sync.Mutex
)

type Translation struct {
	Lang           string `json:"lang"`
	MessageID      string `json:"message_id"`
	TranslationStr string `json:"translation_str"`
}

func LoadMessageFile(path string, bundle *i18n.Bundle, format string, unmarshalFunc i18n.UnmarshalFunc) *i18n.Bundle {
	if format != "" && unmarshalFunc != nil {
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	}
	for _, file := range languageFileMap {
		bundle.MustLoadMessageFile(path + file)
	}
	return bundle
}
func deepCopyBundle(bundle **i18n.Bundle) *i18n.Bundle {
	newBundle := i18n.NewBundle(language.English)
	data, _ := json.Marshal(bundle)
	json.Unmarshal(data, &newBundle)
	return newBundle
}

func Setup() {
	defaultBundle = i18n.NewBundle(language.English)
	defaultBundle = LoadMessageFile(defaultBundlePath, defaultBundle, "json", json.Unmarshal)
	currentBundleN := *defaultBundle
	// currentBundle = deepCopyBundle(&defaultBundle)
	currentBundle = &currentBundleN
	currentBundle = LoadMessageFile(currentBundlePath, currentBundle, "", nil)
}

func NewLocalizer(lang string) MyLocalizer {
	return MyLocalizer{i18n.NewLocalizer(currentBundle, lang)}
}

func (localizer MyLocalizer) F(msg string) string {
	translated, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: msg,
	})
	if err != nil {
		translated = msg
	}
	return translated
}

// 动态添加翻译内容
func AddTranslation(translation Translation, bundle *i18n.Bundle) {
	if bundle == nil {
		currentBundle.AddMessages(language.Make(translation.Lang), &i18n.Message{
			ID:    translation.MessageID,
			Other: translation.TranslationStr,
		})
		translationCache.Store(translation.MessageID, translation)
		return
	}
	bundle.AddMessages(language.Make(translation.Lang), &i18n.Message{
		ID:    translation.MessageID,
		Other: translation.TranslationStr,
	})
}

// 动态删除翻译内容（通过设置标志位实现）
func RemoveTranslation(translation Translation) {
	translationCache.Delete(translation.MessageID)
	mu.Lock()
	needsReload = true
	mu.Unlock()
}

// 查询当前已添加的翻译内容
func GetTranslation() (res []Translation) {
	res = make([]Translation, 0)
	translationCache.Range(func(key, value interface{}) bool {
		translation := value.(Translation)
		res = append(res, translation)
		return true
	})
	return
}

// 同步翻译内容（在适当的时间调用）
func SyncTranslations() {
	mu.Lock()
	if needsReload {
		reloadCurrentBundle()
		needsReload = false
	}
	mu.Unlock()
}

// 同步翻译内容
func reloadCurrentBundle() {
	newBundle := deepCopyBundle(&defaultBundle)
	translationCache.Range(func(key, value interface{}) bool {
		translation := value.(Translation)
		AddTranslation(translation, newBundle)
		return true
	})
	currentBundle = newBundle
}
