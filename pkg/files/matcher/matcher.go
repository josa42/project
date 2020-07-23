package matcher

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

var (
	wildcards = regexp.MustCompile(`{(\*\*?)(\|([^:}]+))?(:([a-z]+))?}`)
)

func findFiles(dir, pattern string) []string {
	globPattern := toGlobPattern(pattern)

	defer cd(dir)()
	matches, _ := zglob.Glob(globPattern)

	sort.Strings(matches)

	return matches
}

func toGlobPattern(pattern string) string {
	matches := wildcards.FindAllStringSubmatch(pattern, -1)

	for _, m := range matches {
		pattern = strings.Replace(pattern, m[0], starToGlob(m[1]), 1)
	}

	return pattern
}

func toExpr(pattern string) string {
	matches := wildcards.FindAllStringSubmatch(pattern, -1)

	replaces := map[string]string{}

	for idx, m := range matches {
		placeholder := fmt.Sprintf(`::::####%d####:::`, idx)
		pattern = strings.Replace(pattern, m[0], placeholder, 1)

		replaces[placeholder] = starToExpr(m[1])
	}

	pattern = regexp.QuoteMeta(pattern)

	for placeholder, replace := range replaces {
		pattern = strings.Replace(pattern, placeholder, replace, 1)
	}

	return pattern
}

func groupNames(pattern string) map[int]string {
	matches := wildcards.FindAllStringSubmatch(pattern, -1)

	groups := map[int]string{}

	for idx, m := range matches {
		groups[idx] = m[5]
	}

	if len(groups) == 1 && groups[0] == "" {
		groups[0] = "path"
	}

	return groups
}

func starToGlob(star string) string {
	switch star {
	case "*":
		return "*"
	case "**":
		return "**/*"
	default:
		return star
	}
}

func starToExpr(star string) string {
	switch star {
	case "*":
		return fmt.Sprintf(`([^%c]+)`, filepath.Separator)
	case "**":
		return `(.+)`
	default:
		return star
	}
}

func cd(dir string) func() {
	pwd, _ := os.Getwd()
	os.Chdir(dir)

	return func() {
		os.Chdir(pwd)
	}
}

type FilePattern string

func (fp FilePattern) Find(dir string) []string {
	pattern := toGlobPattern(string(fp))
	return findFiles(dir, pattern)
}

func (fp FilePattern) Match(path string) map[string]string {

	gn := groupNames(string(fp))
	exp := regexp.MustCompile(toExpr(string(fp)))
	matches := exp.FindAllStringSubmatch(path, -1)

	if len(matches) > 0 {
		groups := map[string]string{}

		for idx, m := range matches {
			name := gn[idx]
			groups[name] = m[1]
		}

		return groups

	}

	return nil
}

func (fp FilePattern) Fill(groups map[string]string) (string, error) {

	matches := wildcards.FindAllStringSubmatch(string(fp), -1)
	gn := groupNames(string(fp))

	path := string(fp)
	for idx, m := range matches {
		name := gn[idx]
		path = strings.Replace(path, m[0], groups[name], 1)
	}

	return path, nil
}

func (fp FilePattern) Groups(filePath string) map[string]string {
	gn := groupNames(string(fp))
	exp := regexp.MustCompile(toExpr(string(fp)))
	matches := exp.FindAllStringSubmatch(filePath, -1)

	groups := map[string]string{}
	for idx, m := range matches {
		groups[gn[idx]] = m[1]
	}

	return groups
}
