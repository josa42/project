package readme

import (
	"strings"

	"github.com/josa42/project/pkg/template"
)

var head = `# {{.Name}}`
var license = `## License

[{{.License}} © {{.Author}}](LICENSE)`

type Placeholders interface {
	Name() string
	License() string
	Author() string
	Year() string
}

func Get(p Placeholders) string {

	sections := []string{
		template.Apply(head, 80, p),
	}

	if p.License() != "" {
		sections = append(sections, template.Apply(license, 80, p))
	}

	return strings.Join(sections, "\n\n")
}
