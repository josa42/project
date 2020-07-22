package matcher

import (
	"reflect"
	"testing"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toExpr(tt.args.pattern); got != tt.want {
				t.Errorf("toExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupNames(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want map[int]string
	}{
		{"empty", args{""}, map[int]string{}},
		{"named", args{"{**:path}/file.js"}, map[int]string{0: "path"}},
		{"transformed", args{"{**|dashed}/file.js"}, map[int]string{0: "path"}},
		{"named and transformed", args{"{**|dashed:path}/file.js"}, map[int]string{0: "path"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := groupNames(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupNames() = %v, want %v", got, tt.want)
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
		fp   FilePattern
		args args
		want map[string]string
	}{
		{"empty", "{**}", args{""}, nil},
		{"named", "{**:path}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
		{"transformed", "{**|dashed}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
		{"named and transformed", "{**|dashed:path}/file.js", args{"foo/bar/file.js"}, map[string]string{"path": "foo/bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fp.Match(tt.args.path); !reflect.DeepEqual(got, tt.want) {
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
		fp      FilePattern
		args    args
		want    string
		wantErr bool
	}{
		{"basic", "src/{*}.js", args{map[string]string{"path": "foo"}}, "src/foo.js", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fp.Fill(tt.args.groups)
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
