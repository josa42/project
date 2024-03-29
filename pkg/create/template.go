package create

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/josa42/project/pkg/license"
	"github.com/josa42/project/pkg/out"
	"github.com/josa42/project/pkg/readme"
	"github.com/josa42/project/pkg/template"
	"gopkg.in/yaml.v2"
)

type Config struct {
	License string `yaml:"license"`
	Author  string `yaml:"author"`
	Email   string `yaml:"email"`
}

type Template struct {
	path  string     `yaml:"-"`
	Files []FileNode `yaml:"files"`
	Init  []string   `yaml:"init"`
}

func DefaultTemplate() *Template {
	t := &Template{}
	t.Init = []string{}

	return t
}

func GetConfig() *Config {
	return &Config{
		License: "MIT",
		Author:  "Josa Gesell",
		Email:   "josa@gesell.me",
	}
}

func LoadTemplate(path string) *Template {
	t := &Template{path: path}

	data, err := ioutil.ReadFile(filepath.Join(path, "project-create.yml"))
	if err != nil {
		out.Logf(`error: %v`, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &t)
	if err != nil {
		out.Logf(`error: %v`, err)
		os.Exit(1)
	}

	return t
}

func (t *Template) Create(baseDir string) error {
	c := GetConfig()

	t.CreateReadme(baseDir, c)
	t.CreateLicense(baseDir, c)
	t.CreateFileTree(baseDir, c)
	t.RunGitInit(baseDir, c)
	t.RunInit(baseDir, c)
	t.RunCommit(baseDir, c)

	return nil
}

func (t *Template) CreateLicense(baseDir string, config *Config) error {
	out.Log("Create: LICENSE")

	l := license.Get(config.License, t.placeholders(baseDir, config))

	ioutil.WriteFile(filepath.Join(baseDir, "LICENSE"), []byte(l), 0644)

	return nil
}

func (t *Template) CreateReadme(baseDir string, config *Config) error {
	out.Log("Create: README.md")
	r := readme.Get(t.placeholders(baseDir, config))

	ioutil.WriteFile(filepath.Join(baseDir, "README.md"), []byte(r), 0644)

	return nil
}

func (t *Template) CreateFileTree(baseDir string, config *Config) error {
	out.Log("Create: Filetree")
	createFiles(baseDir, t.Files, t.placeholders(baseDir, config))

	return nil
}

func (t *Template) RunInit(baseDir string, config *Config) error {

	for _, init := range t.Init {
		init = template.Apply(init, t.placeholders(baseDir, config))
		out.Logf("Run: %s", init)
		if err := run(baseDir, "bash", "-c", init); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) RunGitInit(baseDir string, config *Config) error {
	out.Logf("Run: Commit")

	if err := run(baseDir, "git", "init"); err != nil {
		return err
	}
	return nil
}

func (t *Template) RunCommit(baseDir string, config *Config) error {
	out.Logf("Run: Commit")

	if err := run(baseDir, "git", "add", "-A"); err != nil {
		return err
	}
	if err := run(baseDir, "git", "commit", "-m", "🎉 Initial Commit"); err != nil {
		return err
	}
	return nil
}

func (t *Template) placeholders(baseDir string, config *Config) *Placeholders {
	return &Placeholders{config, t, baseDir}
}

func run(baseDir, bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = baseDir
	return cmd.Run()
}
