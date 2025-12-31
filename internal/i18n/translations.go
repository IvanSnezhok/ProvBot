package i18n

import "fmt"

// Translator interface for translations
type Translator interface {
	Get(key string) string
	Getf(key string, args ...interface{}) string
}

// TranslationStore holds translations for all languages
type TranslationStore struct {
	translations map[string]map[string]string
}

var globalStore *TranslationStore

func init() {
	globalStore = NewTranslationStore()
	LoadTranslations()
}

// NewTranslationStore creates a new translation store
func NewTranslationStore() *TranslationStore {
	return &TranslationStore{
		translations: make(map[string]map[string]string),
	}
}

// RegisterLanguage registers translations for a language
func (ts *TranslationStore) RegisterLanguage(lang string, translations map[string]string) {
	ts.translations[lang] = translations
}

// GetTranslation returns translation for a key in specified language
func (ts *TranslationStore) GetTranslation(lang, key string) string {
	if translations, ok := ts.translations[lang]; ok {
		if translation, ok := translations[key]; ok {
			return translation
		}
	}
	// Fallback to English if key not found
	if lang != "en" {
		if translations, ok := ts.translations["en"]; ok {
			if translation, ok := translations[key]; ok {
				return translation
			}
		}
	}
	return key // Return key if translation not found
}

// GetTranslationf returns formatted translation
func (ts *TranslationStore) GetTranslationf(lang, key string, args ...interface{}) string {
	return fmt.Sprintf(ts.GetTranslation(lang, key), args...)
}

// GetTranslator returns a translator for a specific language
func (ts *TranslationStore) GetTranslator(lang string) Translator {
	return &translator{
		store: ts,
		lang:  lang,
	}
}

// GetGlobalTranslator returns translator from global store
func GetGlobalTranslator(lang string) Translator {
	return globalStore.GetTranslator(lang)
}

// translator implements Translator interface
type translator struct {
	store *TranslationStore
	lang  string
}

func (t *translator) Get(key string) string {
	return t.store.GetTranslation(t.lang, key)
}

func (t *translator) Getf(key string, args ...interface{}) string {
	return t.store.GetTranslationf(t.lang, key, args...)
}

// Default language
const DefaultLanguage = "ua"

// Supported languages
const (
	LanguageUA = "ua"
	LanguageEN = "en"
)

// LoadTranslations loads all translations
func LoadTranslations() {
	globalStore.RegisterLanguage(LanguageUA, GetUATranslations())
	globalStore.RegisterLanguage(LanguageEN, GetENTranslations())
}

