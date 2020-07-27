package matcher

import (
	"reflect"
	"regexp"
	"testing"
)

func fp(p string) FilePattern {
	return FilePattern{
		Path: p,
	}
}

func TestFindFiles(t *testing.T) {
	type args struct {
		dir     string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"basic", args{"testdata/js", "models/{*}.js"}, []string{"models/Bar.js", "models/Foo.js"}},
		{"recursive", args{"testdata/js", "models/{**}.js"}, []string{"models/Bar.js", "models/Foo.js", "models/scoped/Foo.js"}},
		{"prefix", args{"testdata/js", "models/F{*}.js"}, []string{"models/Foo.js"}},
		{"named", args{"testdata/js", "models/{*:name}.js"}, []string{"models/Bar.js", "models/Foo.js"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findFiles(tt.args.dir, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toGlobPattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"single star", args{"dir/{*}.js"}, "dir/*.js"},
		{"double star", args{"dir/{**}.js"}, "dir/**/*.js"},
		{"transform", args{"dir/{*|dashed}.js"}, "dir/*.js"},
		{"named", args{"dir/{*:path}.js"}, "dir/*.js"},
		{"transform and named", args{"dir/{*|dashed:path}.js"}, "dir/*.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toGlobPattern(tt.args.pattern); got != tt.want {
				t.Errorf("toGlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toExpr(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"single star", args{`dir/{*}.js`}, `dir/([^/]+)\.js`},
		{"double star", args{`dir/{**}.js`}, `dir/(.+)\.js`},
		{"transform", args{`dir/{*|dashed}.js`}, `dir/([^/]+)\.js`},
		{"named", args{`dir/{*:path}.js`}, `dir/([^/]+)\.js`},
		{"named 2", args{`{**:path}/file.js`}, `(.+)/file\.js`},
		{"transform and named", args{`dir/{*|dashed:path}.js`}, `dir/([^/]+)\.js`},
		{"constant", args{"src/{controllers}/file.js"}, `src/(controllers)/file\.js`},
		{"constant and named", args{"src/{controllers:type}/file.js"}, `src/(controllers)/file\.js`},
		{"constant, transformed and named", args{"src/{controllers|dashed:type}/file.js"}, `src/(controllers)/file\.js`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toExpr(tt.args.pattern); got != tt.want {
				t.Errorf("toExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_Match(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		fp   string
		args args
		want map[string]string
	}{
		{"empty", "{**}", args{""}, nil},
		{"named", "{**:path}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
		{"transformed", "{**|dashed}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
		{"named and transformed", "{**|dashed:path}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
		{"constant", "src/{controllers}/file.js", args{"src/controllers/file.js"}, map[string]string{"path": "controllers"}},
		{"constant and named", "src/{controllers:type}/file.js", args{"src/controllers/file.js"}, map[string]string{"type": "controllers"}},
		{"constant, transformed and named", "src/{controllers|dashed:type}/file.js", args{"src/controllers/file.js"}, map[string]string{"type": "controllers"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fp(tt.fp).Match(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_starToGlob(t *testing.T) {
	type args struct {
		star string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := starToGlob(tt.args.star); got != tt.want {
				t.Errorf("starToGlob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_starToExpr(t *testing.T) {
	type args struct {
		star string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := starToExpr(tt.args.star); got != tt.want {
				t.Errorf("starToExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_Fill(t *testing.T) {
	type args struct {
		groups map[string]string
	}
	tests := []struct {
		name    string
		fp      string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", "src/{*}.js", args{map[string]string{"path": "foo"}}, "src/foo.js", false},
		{"constant", "src/{controllers}/file.js", args{map[string]string{"path": "controllers"}}, "src/controllers/file.js", false},
		{"constant and named", "src/{controllers:type}/file.js", args{map[string]string{"type": "controllers"}}, "src/controllers/file.js", false},
		{"constant, transformed and named", "src/{controllers|dashed:type}/file.js", args{map[string]string{"type": "controllers"}}, "src/controllers/file.js", false},
		{"constant | error", "src/{controllers}/file.js", args{map[string]string{"path": "models"}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fp(tt.fp).Fill(tt.args.groups)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilePattern.Fill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilePattern.Fill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_Find(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		fp   FilePattern
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fp.Find(tt.args.dir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_Groups(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		fp   string
		args args
		want map[string]string
	}{
		{"", "app/{controllers:type}/{**:path}.js", args{"app/controllers/account/billing.js"}, map[string]string{
			"type": "controllers",
			"path": "account/billing",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fp(tt.fp).Groups(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.Groups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_String(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("Group.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Pattern(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.Pattern(); got != tt.want {
				t.Errorf("Group.Pattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Name(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		// TODO: Add test cases.
		// {"empty", "", map[int]string{}},
		// {"named", "{**:path}/file.js", map[int]string{0: "path"}},
		// {"transformed", "{**|dashed}/file.js", map[int]string{0: "path"}},
		// {"named and transformed", "{**|dashed:path}/file.js", map[int]string{0: "path"}},
		// {"constant", "src/{controllers}/file.js", map[int]string{0: "path"}},
		// {"constant and named", "src/{controllers:type}/file.js", map[int]string{0: "type"}},
		// {"constant, transformed and named", "src/{controllers|dashed:type}/file.js", map[int]string{0: "type"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Group{str: tt.path}).Name(); got != tt.want {
				t.Errorf("Group.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_Transform(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.Transform(); got != tt.want {
				t.Errorf("Group.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_IsConstant(t *testing.T) {
	type fields struct {
		str  string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group{
				str:  tt.fields.str,
				name: tt.fields.name,
			}
			if got := g.IsConstant(); got != tt.want {
				t.Errorf("Group.IsConstant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_Expr(t *testing.T) {
	type fields struct {
		Path           string
		ConstantGroups map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   *regexp.Regexp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := FilePattern{
				Path:           tt.fields.Path,
				ConstantGroups: tt.fields.ConstantGroups,
			}
			if got := fp.Expr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.Expr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_GroupMatches(t *testing.T) {
	type fields struct {
		Path           string
		ConstantGroups map[string]string
	}
	tests := []struct {
		name string
		path string
		want []Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fp(tt.path).GroupMatches(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.GroupMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilePattern_String(t *testing.T) {
	type fields struct {
		Path           string
		ConstantGroups map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := FilePattern{
				Path:           tt.fields.Path,
				ConstantGroups: tt.fields.ConstantGroups,
			}
			if got := fp.String(); got != tt.want {
				t.Errorf("FilePattern.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
