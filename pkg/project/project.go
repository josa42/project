package project

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/josa42/project/pkg/files/matcher"
	"gopkg.in/yaml.v2"
)

func FindProjeFile(dirPath string) string {
	if !filepath.IsAbs(dirPath) {
		dirPath, _ = filepath.Abs(dirPath)
	}

	filePath := filepath.Join(dirPath, "project.yml")
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	}

	dirPath, last := filepath.Split(dirPath)
	if last == "" {
		return ""
	}

	return FindProjeFile(dirPath)
}

func LoadProjeFile(filePath string) *Project {
	p := Project{}

	content, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(content, &p)

	for key, task := range p.Tasks {
		task.Key = key
		p.Tasks[key] = task
	}

	for key, file := range p.Files {
		file.Key = key
		p.Files[key] = file
	}

	return &p
}

type Project struct {
	Files map[string]FileType `yaml:"files"`
	Tasks map[string]Task     `yaml:"tasks"`
}

func (p Project) FindFiles(key string) []string {
	files := []string{}

	ft := p.Files[key]

	for _, fp := range ft.PathPatterns {
		for _, f := range fp.Find(".") {
			if !ft.isExcluded(f) {
				files = append(files, f)
			}
		}
	}

	return files
}

type FileType struct {
	Key             string                `yaml:"-"`
	PathPatterns    []matcher.FilePattern `yaml:"path"`
	ExcludePatterns []matcher.FilePattern `yaml:"exclude"`
	RelatedKeys     []string              `yaml:"related"`
}

func (ft *FileType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := map[string]interface{}{}
	err := unmarshal(&d)
	if err != nil {
		return err
	}

	ft.PathPatterns = patternOrSlice(d["path"])
	ft.ExcludePatterns = patternOrSlice(d["exclude"])
	ft.RelatedKeys = stringOrSlice(d["related"])

	return nil
}

func (ft FileType) isExcluded(filePath string) bool {
	for _, ex := range ft.ExcludePatterns {
		if ex.Match(filePath) != nil {
			return true
		}
	}
	return false
}

func patternOrSlice(in interface{}) []matcher.FilePattern {
	patterns := []matcher.FilePattern{}
	for _, p := range stringOrSlice(in) {
		patterns = append(patterns, matcher.FilePattern(p))
	}
	return patterns
}

func stringOrSlice(in interface{}) []string {
	rel := []string{}

	if v, ok := in.(string); ok {
		rel = append(rel, v)
	} else if v, ok := in.([]interface{}); ok {
		for _, r := range v {
			if rv, ok := r.(string); ok {
				rel = append(rel, rv)
			}
		}
	}

	return rel

}

type Task struct {
	Key     string
	Command string
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cmd := ""
	if err := unmarshal(&cmd); err != nil {
		return err
	}

	t.Command = cmd

	return nil
}
