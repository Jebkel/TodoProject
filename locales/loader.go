package locales

import (
	"github.com/eduardolat/goeasyi18n"
)

func Init() *goeasyi18n.I18n {
	i18n := goeasyi18n.NewI18n()

	ruTranslations, err := goeasyi18n.LoadFromJsonFiles("locales/ru/*.json")
	if err != nil {
		panic(err)
	}

	enTranslations, err := goeasyi18n.LoadFromJsonFiles("locales/en/*.json")
	if err != nil {
		panic(err)
	}

	i18n.AddLanguage("ru", ruTranslations)
	i18n.AddLanguage("en", enTranslations)
	return i18n
}
