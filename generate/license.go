package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var tpl = `package license

func init() {
	license("%[1]s", ` + "`%[2]s`" + `)
}

`

func main() {
	keys := os.Args[1:]

	for _, key := range keys {

		fmt.Println(key)

		res, _ := http.Get(fmt.Sprintf(`https://spdx.org/licenses/%s.txt`, key))
		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		license := string(body)
		license = addPlaceholders(license)

		content := fmt.Sprintf(tpl, key, license)

		ioutil.WriteFile(fmt.Sprintf(`license-%s.go`, strings.ToLower(key)), []byte(content), 0644)
	}
}

func addPlaceholders(text string) string {

	text = strings.ReplaceAll(text, `<year>`, `{{.Year}}`)
	text = strings.ReplaceAll(text, `<copyright holders>`, `{{.AuthorFull}}`)

	return text
}

