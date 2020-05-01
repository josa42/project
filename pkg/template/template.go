package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/josa42/go-stringutils"
)

func Apply(text string, width int, p interface{}) string {
	buf := bytes.NewBuffer(nil)

	key := fmt.Sprintf("template-%d", time.Now().Second())
	t := template.Must(template.New(key).Parse(text))
	t.Execute(buf, p)

	r := buf.String()
	r = stringutils.Wrap(r, width)

	return r
}

