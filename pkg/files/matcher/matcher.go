package matcher

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-zglob"
)

var (
	wildcards = regexp.MustCompile(`{(\*\*?)(:([a-z]+))?}`)
)

func FindFiles(dir, pattern string) []string {

	globPattern := toGlobPattern(pattern)

	defer cd(dir)()
	matches, _ := zglob.Glob(globPattern)

	sort.Strings(matches)

	return matches
}

func toGlobPattern(pattern string) string {

	matches := wildcards.FindAllStringSubmatch(pattern, -1)

	for _, m := range matches {
		pattern = strings.Replace(pattern, m[0], starReplace(m[1]), 1)
	}

	return pattern
}

func starReplace(star string) string {
	switch star {
	case "*":
		return "*"
	case "**":
		return "**/*"
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
