package project

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/josa42/project/pkg/files/matcher"
	"github.com/josa42/project/pkg/out"
	"gopkg.in/yaml.v2"
)

func MustLoad(dir string) *Project {
	projPath := FindProjeFile(".")
	if projPath == "" {
		out.Log("Could not find project.yml")
		os.Exit(1)
	}

	return LoadProjeFile(projPath)
}

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

func (p Project) RelatedKeys(filePath string) []string {
	keys := []string{}

	for _, ft := range p.Files {
		if ft.isMatching(filePath) {
			keys = append(keys, ft.RelatedKeys...)
		}
	}

	return keys
}

func (p Project) AllRelatedFiles(filePath string) ([]string, map[string]string) {
	related := map[string]string{}
	keys := []string{}

	for _, key := range p.RelatedKeys(filePath) {
		for _, file := range p.RelatedFiles(key, filePath) {
			if _, err := os.Stat(file); err == nil {
				related[key] = file
				keys = append(keys, key)
				break
			}
		}
	}

	return keys, related
}

func (p Project) RelatedFiles(key, filePath string) []string {
	files := []string{}

	ftr, ok := p.Files[key]
	if !ok {
		return []string{}
	}

	for _, ft := range p.Files {
		if ft.isRelated(key) {

			for _, fp := range ft.PathPatterns {

				if len(fp.Match(filePath)) > 0 {
					g := fp.Groups(filePath)

					for _, fpr := range ftr.PathPatterns {
						if rPath, err := fpr.Fill(g); err == nil {
							files = append(files, rPath)
							break
						}
					}
				}
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

func (ft FileType) isMatching(filePath string) bool {
	if !ft.isExcluded(filePath) {
		for _, ex := range ft.PathPatterns {
			if ex.Match(filePath) != nil {
				return true
			}
		}
	}
	return false
}

func (ft FileType) isRelated(key string) bool {
	for _, akey := range ft.RelatedKeys {
		if akey == key {
			return true
		}
	}
	return false
}

func patternOrSlice(in interface{}) []matcher.FilePattern {
	rel := []matcher.FilePattern{}

	if v, ok := in.(string); ok {
		rel = append(rel, matcher.FilePattern{
			Path: v,
		})
	} else if v, ok := in.([]interface{}); ok {
		for _, r := range v {
			if rv, ok := r.(string); ok {
				rel = append(rel, matcher.FilePattern{
					Path: rv,
				})
			} else if rv, ok := r.(map[interface{}]interface{}); ok {
				for p := range rv {
					if rv, ok := p.(string); ok {
						rel = append(rel, matcher.FilePattern{
							Path: rv,
						})
					}
					break
				}
			}
		}
	}

	return rel
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
