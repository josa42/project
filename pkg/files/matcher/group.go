package matcher

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
	s := g.Pattern()
	return s != "*" && s != "**"

}
