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
			if got := FindFiles(tt.args.dir, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
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
