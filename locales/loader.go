package locales

import (
	"fmt"
	"github.com/eduardolat/goeasyi18n"
	"github.com/labstack/gommon/log"
	"os"
)

func Init() *goeasyi18n.I18n {
	i18n := goeasyi18n.NewI18n()

	entries, err := os.ReadDir("locales/")
	if err != nil {
		panic("no found directories on locales dir")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			langData, err := goeasyi18n.LoadFromJsonFiles(fmt.Sprintf("locales/%s/*.json", entry.Name()))
			if err != nil {
				log.Error(err)
			}
			i18n.AddLanguage(entry.Name(), langData)
		}
	}

	return i18n
}
