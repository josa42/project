package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/josa42/go-changecase"
)

var funcMap = template.FuncMap{
	"camelCase":  changecase.ToCamel,
	"constant":   changecase.ToConstant,
	"dot":        changecase.ToDot,
	"get":        get,
	"lower":      changecase.ToLower,
	"lowerFirst": changecase.ToLowerFirst,
	"param":      changecase.ToParam,
	"pascal":     changecase.ToPascal,
	"path":       changecase.ToPath,
	"replace":    replace,
	"run":        run,
	"snake":      changecase.ToSnake,
	"title":      changecase.ToTitle,
	"upper":      changecase.ToUpper,
	"upperFirst": changecase.ToUpperFirst,
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

func run(command string) string {
	cmd := exec.Command("bash", "-c", command)
	out, _ := cmd.Output()
	return string(out)
}

func get(url string) string {
	resp, _ := http.Get(url)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(content)
}

