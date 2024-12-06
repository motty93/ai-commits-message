package i18n

import (
	"embed"
	"log"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Language string

const (
	JPN Language = "JPN"
	ENG Language = "ENG"
)

func SetLanguage() {
	lang = Language(os.Getenv("LANGUAGE"))
	if lang == "" {
		lang = JPN
	}
}

func SetLanguageTag() {
	if lang == ENG {
		tag = language.English
	} else {
		tag = language.Japanese
	}
}

//go:embed *
var (
	files embed.FS
	lang  Language
	tag   language.Tag
	loc   *i18n.Localizer
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
