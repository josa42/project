package matcher

func fp(p string, cgs ...map[string]string) FilePattern {
	var cg map[string]string
	if len(cgs) == 1 {
		cg = cgs[0]
	}

	return FilePattern{
		Path:           p,
		ConstantGroups: cg,
	}
}
