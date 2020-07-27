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
