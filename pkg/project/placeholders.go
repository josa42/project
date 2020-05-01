package project

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Placeholders struct {
	config   *Config
	template *Template
	slug     string
}

var titleExp = regexp.MustCompile(`[-_]`)

func (p Placeholders) Name() string {
	return strings.Title(titleExp.ReplaceAllString(p.slug, " "))
}

func (p Placeholders) Slug() string {
	return p.slug
}

func (p Placeholders) Author() string {
	return p.config.Author
}

func (p Placeholders) AuthorFull() string {
	return fmt.Sprintf("%s <%s>", p.config.Author, p.config.Email)
}

func (p Placeholders) Email() string {
	return p.config.Email
}

func (p Placeholders) License() string {
	return p.config.License
}

func (p Placeholders) Year() string {
	return fmt.Sprintf("%d", time.Now().Local().Year())
}

