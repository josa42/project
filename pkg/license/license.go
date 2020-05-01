package license

//go:generate go run ../../generate/license.go "MIT"

import "github.com/josa42/project/pkg/template"

type Placeholders interface {
	AuthorFull() string
	Year() string
}

func Get(key string, p Placeholders) string {
	if text, ok := licenses[key]; ok {
		return template.Apply(text, 80, p)
	}

	return ""
}

var licenses = map[string]string{}

func license(key, content string) {
	licenses[key] = content
}
