package matcher

import (
	"errors"
	"regexp"
	"strings"
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

type FilePattern struct {
	Path           string
	ConstantGroups map[string]string
}

func (fp FilePattern) Expr() *regexp.Regexp {
	return regexp.MustCompile(toExpr(fp.Path))
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

func (fp FilePattern) GroupValues(filePath string) map[string]string {
	gn := fp.GroupMatches()
	matches := fp.Expr().FindStringSubmatch(filePath)

	groups := map[string]string{}
	for idx, m := range matches[1:] {
		groups[gn[idx].Name()] = m
	}

	for key, value := range fp.ConstantGroups {
		groups[key] = value
	}

	return groups
}
