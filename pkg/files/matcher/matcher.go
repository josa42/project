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

type Group struct {
	str  string
	name string
}

func (g Group) parse() []string {
	return wildcards.FindStringSubmatch(g.str)
}

func (g Group) String() string {
	return g.str
}

func (g Group) Pattern() string {
	return g.parse()[wildcardIdxPattern]
}

func (g Group) Name() string {
	if n := g.name; n != "" {
		return n
	}
	return g.parse()[wildcardIdxName]
}

func (g Group) Transform() string {
	return g.parse()[wildcardIdxTransform]
}

func (g Group) IsConstant() bool {
	return isConstantGroup(g.Pattern())
}

type FilePattern struct {
	Path           string
	ConstantGroups map[string]string
}

func (fp FilePattern) Expr() *regexp.Regexp {
	return regexp.MustCompile(toExpr(fp.Path))
}

func (fp FilePattern) GroupNames() map[int]string {
	matches := wildcards.FindAllStringSubmatch(fp.Path, -1)

	groups := map[int]string{}

	for idx, m := range matches {
		groups[idx] = m[wildcardIdxName]
	}

	if len(groups) == 1 && groups[0] == "" {
		groups[0] = "path"
	}

	return groups
}

func (fp FilePattern) GroupMatches() []Group {
	matches := wildcards.FindAllString(fp.Path, -1)

	groups := []Group{}

	for _, m := range matches {
		g := Group{str: m}
		if g.Name() == "" && len(matches) == 1 {
			g.name = "path"
		}

		groups = append(groups, g)
	}

	return groups

}

func (fp FilePattern) String() string {
	return fp.Path
}

func (fp FilePattern) Find(dir string) []string {
	pattern := toGlobPattern(fp.Path)
	return findFiles(dir, pattern)
}

func (fp FilePattern) Match(path string) map[string]string {

	gm := fp.GroupMatches()
	matches := fp.Expr().FindAllStringSubmatch(path, -1)

	if len(matches) > 0 {
		groups := map[string]string{}

		for idx, m := range matches {
			name := gm[idx].Name()
			groups[name] = m[wildcardIdxPattern]
		}

		return groups
	}

	return nil
}

func (fp FilePattern) Fill(groupValues map[string]string) (string, error) {

	gm := fp.GroupMatches()

	path := fp.Path

	for _, m := range gm {
		name := m.Name()

		if m.IsConstant() && m.Pattern() != groupValues[name] {
			return "", errors.New("const group not matching")
		}

		path = strings.Replace(path, m.String(), groupValues[name], 1)
	}

	return path, nil
}

func (fp FilePattern) Groups(filePath string) map[string]string {
	gn := fp.GroupNames()
	matches := fp.Expr().FindStringSubmatch(filePath)

	groups := map[string]string{}
	for idx, m := range matches[1:] {
		groups[gn[idx]] = m
	}

	return groups
}

func isConstantGroup(s string) bool {
	return s != "*" && s != "**"
}
