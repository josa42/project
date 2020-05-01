package license

//go:generate go run ../../generate/license.go "MIT"

import (
	"github.com/josa42/go-stringutils"
	"github.com/josa42/project/pkg/template"
)

type Placeholders interface {
	AuthorFull() string
	Year() string
}

func Get(key string, p Placeholders) string {
	if text, ok := licenses[key]; ok {
		text = stringutils.Wrap(text, 80)
		return template.Apply(text, p)
	}

	return ""
}

var licenses = map[string]string{}

func license(key, content string) {
	licenses[key] = content
}
