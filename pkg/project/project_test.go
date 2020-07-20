package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func fixture(p ...string) string {
	rootPath, _ := os.Getwd()
	return filepath.Join(append([]string{rootPath, "testdata"}, p...)...)
}

func TestFileType(t *testing.T) {

	t.Run("single", func(t *testing.T) {
		ft := FileType{}

		err := yaml.Unmarshal([]byte(`
path: controlers/{*}.js
related: view`), &ft)
		assert.Nil(t, err)
		assert.Equal(t, "", ft.Key)
		assert.Equal(t, []string{"controlers/{*}.js"}, ft.PathPatterns)
		assert.Equal(t, []string{"view"}, ft.RelatedKeys)
	})

	t.Run("list", func(t *testing.T) {
		ft := FileType{}

		err := yaml.Unmarshal([]byte(`
path:
  - controlers/{*}.js
  - sub/controlers/{*}.js
related:
  - view
  - test`), &ft)
		assert.Nil(t, err)
		assert.Equal(t, "", ft.Key)
		assert.Equal(t, []string{"controlers/{*}.js", "sub/controlers/{*}.js"}, ft.PathPatterns)
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
	assert.Equal(t, []string{"{**}_test.go"}, proj.Files["test"].PathPatterns)
	assert.Equal(t, []string{"source"}, proj.Files["test"].RelatedKeys)

}
