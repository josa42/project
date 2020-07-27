package matcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

var (
	wildcards = regexp.MustCompile(`{(\*\*?|[^:\|]+)(\|([^:}]+))?(:([a-z]+))?}`)
)

const (
	wildcardIdxComplete  = 0
	wildcardIdxPattern   = 1
	wildcardIdxTransform = 3
	wildcardIdxName      = 5
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
		pattern = strings.Replace(pattern, m[wildcardIdxComplete], starToGlob(m[wildcardIdxPattern]), 1)
	}

	return pattern
}

func toExpr(pattern string) string {
	matches := wildcards.FindAllStringSubmatch(pattern, -1)

	replaces := map[string]string{}

	for idx, m := range matches {
		placeholder := fmt.Sprintf(`::::####%d####:::`, idx)
		pattern = strings.Replace(pattern, m[wildcardIdxComplete], placeholder, 1)

		replaces[placeholder] = starToExpr(m[wildcardIdxPattern])
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
		groups[idx] = m[wildcardIdxName]
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
		return fmt.Sprintf("(%s)", regexp.QuoteMeta(star))
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
			groups[name] = m[wildcardIdxPattern]
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

		if isConstantGroup(m[wildcardIdxPattern]) && m[wildcardIdxPattern] != groups[name] {
			return "", errors.New("const group not matching")
		}

		path = strings.Replace(path, m[wildcardIdxComplete], groups[name], 1)
	}

	return path, nil
}

func (fp FilePattern) Groups(filePath string) map[string]string {
	gn := groupNames(string(fp))
	exp := regexp.MustCompile(toExpr(string(fp)))
	matches := exp.FindAllStringSubmatch(filePath, -1)

	groups := map[string]string{}
	for idx, m := range matches {
		groups[gn[idx]] = m[wildcardIdxPattern]
	}

	return groups
}

func isConstantGroup(s string) bool {
	return s != "*" && s != "**"
}
