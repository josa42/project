package matcher

import (
	"reflect"
	"regexp"
	"testing"
)

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

func TestFilePattern_GroupValues(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		fp   FilePattern
		args args
		want map[string]string
	}{
		{"", fp("app/{controllers:type}/{**:path}.js"), args{"app/controllers/account/billing.js"}, map[string]string{
			"type": "controllers",
			"path": "account/billing",
		}},
		{"", fp("app/controllers/{**:path}.js", map[string]string{"type": "controllers"}), args{"app/controllers/account/billing.js"}, map[string]string{
			"type": "controllers",
			"path": "account/billing",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fp.GroupValues(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilePattern.Groups() = %v, want %v", got, tt.want)
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
