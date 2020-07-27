package project

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/josa42/project/pkg/files/matcher"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func fixtureDir(dir string) func() {
	pwd, _ := os.Getwd()
	os.Chdir(path.Join("testdata", dir))

	return func() {
		os.Chdir(pwd)
	}
}

func fixture(p ...string) string {
	rootPath, _ := os.Getwd()
	return filepath.Join(append([]string{rootPath, "testdata"}, p...)...)
}

func fps(ps ...string) []matcher.FilePattern {
	fps := []matcher.FilePattern{}
	for _, p := range ps {
		fps = append(fps, matcher.FilePattern{
			Path: p,
		})
	}
	return fps
}

func TestFileType(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		ft := FileType{}

		err := yaml.Unmarshal([]byte(`
path: controllers/{*}.js
exclude: controllers/{*}.test.js
related: view`), &ft)
		assert.Nil(t, err)
		assert.Equal(t, "", ft.Key)
		assert.Equal(t, fps("controllers/{*}.js"), ft.PathPatterns)
		assert.Equal(t, fps("controllers/{*}.test.js"), ft.ExcludePatterns)
		assert.Equal(t, []string{"view"}, ft.RelatedKeys)
	})

	t.Run("list", func(t *testing.T) {
		ft := FileType{}

		err := yaml.Unmarshal([]byte(`
path:
  - controllers/{*}.js
  - sub/controllers/{*}.js
exclude: [ "controllers/{*}.test.js" ]
related:
  - view
  - test`), &ft)
		assert.Nil(t, err)
		assert.Equal(t, "", ft.Key)
		assert.Equal(t, fps("controllers/{*}.js", "sub/controllers/{*}.js"), ft.PathPatterns)
		assert.Equal(t, fps("controllers/{*}.test.js"), ft.ExcludePatterns)
		assert.Equal(t, []string{"view", "test"}, ft.RelatedKeys)
	})

}

func TestFindProjeFile(t *testing.T) {
	t.Run("same dir", func(t *testing.T) {
		assert.Equal(t, fixture("a", "project.yml"), FindProjeFile(fixture("a")))
	})

	t.Run("parent dir", func(t *testing.T) {
		assert.Equal(t, fixture("a", "sub-a", "project.yml"), FindProjeFile(fixture("a", "sub-a")))
	})
}

func TestLoadProjeFile(t *testing.T) {
	proj := LoadProjeFile(fixture("complete", "project.yml"))

	assert.Equal(t, "test", proj.Tasks["test"].Key)
	assert.Equal(t, "go test ./...", proj.Tasks["test"].Command)

	assert.Equal(t, "test", proj.Files["test"].Key)
	assert.Equal(t, fps("{**}_test.go"), proj.Files["test"].PathPatterns)
	assert.Equal(t, []string{"source"}, proj.Files["test"].RelatedKeys)
}

func TestProject_RelatedFiles(t *testing.T) {
	type args struct {
		key      string
		filePath string
	}
	tests := []struct {
		name string
		dir  string
		args args
		want []string
	}{
		{"ember controller => test", "emberjs-1", args{"test", "app/controllers/account/billing.js"}, []string{
			"tests/unit/controllers/account/billing-test.js",
		}},
		// {"ember controller => route", "emberjs-1", args{"route", "app/controllers/account/billing.js"}, []string{
		// 	"app/routes/account/billing.js",
		// }},
		{"ember controller", "emberjs-1", args{"template", "app/controllers/account/billing.js"}, []string{
			"app/templates/account/billing.hbs",
		}},
		{"ember route => test", "emberjs-1", args{"test", "app/routes/account/billing.js"}, []string{
			"tests/unit/routes/account/billing-test.js",
		}},
		// {"ember route => controller", "emberjs-1", args{"controller", "app/routes/account/billing.js"}, []string{
		// 	"app/controllers/account/billing.js",
		// }},
		{"ember route => template", "emberjs-1", args{"template", "app/routes/account/billing.js"}, []string{
			"app/templates/account/billing.hbs",
		}},
		{"ember component => test", "emberjs-1", args{"test", "app/components/my-component.js"}, []string{
			"tests/unit/components/my-component-test.js",
		}},
		{"ember component => template", "emberjs-1", args{"template", "app/components/my-component.js"}, []string{
			"app/templates/components/my-component.hbs",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer fixtureDir(tt.dir)()
			p := MustLoad(".")

			if got := p.RelatedFiles(tt.args.key, tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Project.RelatedFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
