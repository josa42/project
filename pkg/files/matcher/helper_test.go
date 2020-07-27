package matcher

import (
	"reflect"
	"testing"
)

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
