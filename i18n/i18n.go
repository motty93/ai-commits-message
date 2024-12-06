package i18n

import (
	"embed"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const (
	JPN         string = "JPN"
	ENG         string = "ENG"
	DefaultLang string = JPN
)

func SetLanguage() {
	if lang == "" {
		lang = DefaultLang
	}
}

func SetLanguageTag() {
	if lang == ENG {
		tag = language.English
	} else {
		tag = language.Japanese
	}
}

var (
	//go:embed *.yaml
	files embed.FS

	lang string
	tag  language.Tag
	loc  *i18n.Localizer
)

// WANT: new service
func Init() {
	SetLanguage()
	SetLanguageTag()

	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if _, err := bundle.LoadMessageFileFS(files, string(lang)+".yaml"); err != nil {
		log.Fatalf("Cannot load message file: %v", err)
	}

	loc = i18n.NewLocalizer(bundle)
}

func GetText(key string) string {
	if loc == nil {
		return ""
	}

	return loc.MustLocalize(&i18n.LocalizeConfig{MessageID: key})
}
