package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/josa42/go-changecase"
)

var funcMap = template.FuncMap{
	"camelCase":  changecase.ToCamel,
	"lower":      changecase.ToLower,
	"lowerFirst": changecase.ToLowerFirst,
	"replace":    replace,
	"title":      changecase.ToTitle,
	"upper":      changecase.ToUpper,
	"upperFirst": changecase.ToUpperFirst,
	"pascal":     changecase.ToPascal,
	"snake":      changecase.ToSnake,
	"param":      changecase.ToParam,
	"constant":   changecase.ToConstant,
	"dot":        changecase.ToDot,
	"path":       changecase.ToPath,
}

func Apply(text string, p interface{}) string {
	buf := bytes.NewBuffer(nil)

	key := fmt.Sprintf("template-%d", time.Now().Second())
	t := template.Must(template.New(key).Funcs(funcMap).Parse(text))
	t.Execute(buf, p)

	r := buf.String()

	return r
}

var isRegex = regexp.MustCompile(`^/.*/$`)

func replace(str, exp, replace string) string {
	if isRegex.MatchString(exp) {
		return regexp.MustCompile(exp[1:len(exp)-1]).ReplaceAllString(str, replace)
	}
	return strings.ReplaceAll(str, exp, replace)
}
