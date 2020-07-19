package create

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/josa42/project/pkg/template"
)

func createFiles(baseDir string, nodes []FileNode, p *Placeholders) {
	for _, node := range nodes {
		if !node.IsValid() {
			continue
		}

		if node.IsDir() {
			dirname := filepath.Join(baseDir, node.Name(p))
			os.MkdirAll(dirname, 0755)

			createFiles(dirname, node.Files, p)
		}

		if node.IsFile() {
			filename := filepath.Join(baseDir, node.Name(p))
			ioutil.WriteFile(filename, []byte(node.FileContent(p)), 0644)
		}
	}

}

type FileNode struct {
	Dir          string     `yaml:"dir"`
	Files        []FileNode `yaml:"files"`
	File         string     `yaml:"file"`
	Template     string     `yaml:"template"`
	TemplateFile string     `yaml:"template_file"`
	TemplateURL  string     `yaml:"template_url"`
	Content      string     `yaml:"content"`
	ContentFile  string     `yaml:"content_file"`
	ContentURL   string     `yaml:"content_url"`
}

func (n FileNode) IsDir() bool {
	return n.Dir != "" && n.File == ""
}

func (n FileNode) IsFile() bool {
	return n.File != "" && n.Dir == ""
}

func (n FileNode) IsValid() bool {
	return n.IsDir() || n.IsFile()
}

func (n FileNode) Name(p *Placeholders) string {
	if n.IsDir() {
		return template.Apply(n.Dir, p)
	}
	return template.Apply(n.File, p)
}

func (n FileNode) FileContent(p *Placeholders) string {

	if n.TemplateFile != "" {
		tmpl := readFile(p.template.path, n.TemplateFile)
		return template.Apply(tmpl, p)
	}

	if n.TemplateURL != "" {
		tmpl := readURL(n.TemplateURL)
		return template.Apply(tmpl, p)
	}

	if n.Template != "" {
		return template.Apply(n.Template, p)
	}

	if n.ContentURL != "" {
		return readURL(n.ContentURL)
	}

	if n.ContentFile != "" {
		return readFile(p.template.path, n.ContentFile)
	}

	return n.Content
}

func readFile(tmplDir, path string) string {
	filePath := filepath.Join(tmplDir, "files", path)
	content, _ := ioutil.ReadFile(filePath)

	return string(content)

}
func readURL(url string) string {
	resp, _ := http.Get(url)
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(content)
}

