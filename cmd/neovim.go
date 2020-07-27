/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/josa42/project/pkg/project"
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
	"github.com/spf13/cobra"
)

func isAny(str string, options []string) bool {
	for _, o := range options {
		if o == str {
			return true
		}
	}
	return false
}

func currentFilePath(p *plugin.Plugin) string {
	b, _ := p.Nvim.CurrentBuffer()

	return bufferFilePath(p, b)
}

func bufferFilePath(p *plugin.Plugin, b nvim.Buffer) string {
	filePath := ""
	p.Nvim.Call("expand", &filePath, fmt.Sprintf("#%d:p", b))

	return filePath
}

func completeRelatedKey(p *plugin.Plugin) func(args []interface{}) ([]string, error) {
	return func(args []interface{}) ([]string, error) {
		proj := project.MustLoad(".")
		filePath := currentFilePath(p)
		keys := []string{}

		for _, rk := range proj.RelatedKeys(filePath) {
			keys = append(keys, rk.String())
		}

		return keys, nil
	}
}

func completeKey(p *plugin.Plugin) func(args []interface{}) ([]string, error) {
	return func(args []interface{}) ([]string, error) {
		proj := project.MustLoad(".")
		keys := []string{}

		for key := range proj.Files {
			keys = append(keys, key)
		}

		return keys, nil
	}
}

func completeOpen(p *plugin.Plugin) func(args []interface{}) ([]string, error) {
	return func(args []interface{}) ([]string, error) {
		pre := strings.Split(stringArg(args, 0), " ")

		proj := project.MustLoad(".")

		if len(pre) <= 1 {
			keys := []string{}

			for key := range proj.Files {
				keys = append(keys, key)
			}

			return fuzzyMatch(keys, pre[0]), nil
		}

		key := pre[0]
		pathPre := strings.Join(pre[1:], " ")

		files := fuzzyMatch(proj.FindFiles(key), pathPre)

		return prefix(files, fmt.Sprintf("%s ", key)), nil
	}
}

func findTabWindow(p *plugin.Plugin, filePath string) (nvim.Tabpage, nvim.Window) {

	tabs, _ := p.Nvim.Tabpages()

	for _, tab := range tabs {
		wins, _ := p.Nvim.TabpageWindows(tab)
		for _, win := range wins {
			buf, _ := p.Nvim.WindowBuffer(win)
			filePathAbs, _ := filepath.Abs(filePath)
			if bufferFilePath(p, buf) == filePathAbs {
				return tab, win
			}
		}
	}

	return nvim.Tabpage(-1), nvim.Window(-1)
}

func stringArg(args []interface{}, idx int) string {
	if len(args) < idx+1 {
		return ""
	}

	if str, ok := args[idx].(string); ok {
		return str
	}

	return ""
}

func commandArg(args []string) (string, bool, error) {
	target := "tab"
	if len(args) >= 1 {
		target = args[0]
	}

	command := ""
	force := strings.HasSuffix(target, "!")

	switch strings.Replace(target, "!", "", 1) {
	case "window":
		command = "edit"
	case "tab":
		command = "tabedit"
	case "split":
		command = target
	case "vsplit":
		command = target
	default:
		return "", force, errors.New("command missing")
	}

	return command, force, nil
}

func alternate(p *plugin.Plugin) func(args []string) error {
	return func(args []string) error {

		command, force, err := commandArg(args)
		if err != nil {
			return err
		}

		forceKey := ""
		if len(args) >= 2 {
			forceKey = args[1]
		}

		proj := project.MustLoad(".")
		filePath := currentFilePath(p)

		if forceKey != "" {
			files := proj.RelatedFiles(forceKey, filePath)

			for _, f := range files {
				return openFile(p, command, f, force)
			}
			return nil

		}

		key := ""
		keys, related := proj.AllRelatedFiles(filePath)

		if len(keys) == 1 {
			key = keys[0]
		} else if len(keys) > 1 {
			key = prompt(p, "key", "", "CompleteRelatedKey")
		}

		if f, ok := related[key]; ok {
			return openFile(p, command, f, force)
		}

		return nil
	}
}

func openFile(p *plugin.Plugin, command, filePath string, force bool) error {
	if !force {
		if _, win := findTabWindow(p, filePath); win > 0 {
			p.Nvim.SetCurrentWindow(win)
			return nil
		}
	}

	return p.Nvim.Command(fmt.Sprintf(`%s %s`, command, filePath))

}

func open(p *plugin.Plugin) func(args []string) error {
	var runOpen func(args []string) error

	runOpen = func(args []string) error {

		command, force, err := commandArg(args)
		if err != nil {
			return nil
		}

		text := ""
		if len(args) >= 2 {
			if len(args[1]) == 0 {
				return nil
			}

			text = fmt.Sprintf("%s ", args[1])
		}

		inpt := strings.Split(prompt(p, "open", text, "CompleteOpen"), " ")

		if len(inpt) == 1 {
			return runOpen([]string{args[0], inpt[0]})
		}

		if len(inpt) < 2 {
			return nil
		}

		filePath := strings.Join(inpt[1:], " ")

		return openFile(p, command, filePath, force)
	}

	return runOpen
}

func prompt(p *plugin.Plugin, label string, text string, complete string) string {
	out := ""
	p.Nvim.Call("input", &out, fmt.Sprintf("%s: ", label), text, fmt.Sprintf("customlist,%s", complete))
	return out
}

func prefix(options []string, prefix string) []string {

	m := []string{}
	for _, str := range options {
		m = append(m, fmt.Sprintf("%s%s", prefix, str))
	}

	return m
}

func fuzzyMatch(options []string, input string) []string {

	r := fuzzyRegexp(input)
	m := []string{}
	for _, str := range options {
		if r.MatchString(str) {
			m = append(m, str)
		}
	}

	return m
}

func fuzzyRegexp(input string) *regexp.Regexp {
	exp := []string{".*"}
	for _, c := range strings.Split(input, "") {
		exp = append(exp, regexp.QuoteMeta(c), ".*")
	}

	return regexp.MustCompile(strings.Join(exp, ""))
}

// neovimCmd represents the neovim command
var neovimCmd = &cobra.Command{
	Use:    "neovim",
	Short:  "",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Main(func(p *plugin.Plugin) error {
			p.HandleFunction(&plugin.FunctionOptions{Name: "Alternate"}, alternate(p))
			p.HandleFunction(&plugin.FunctionOptions{Name: "ProjectOpen"}, open(p))
			p.HandleFunction(&plugin.FunctionOptions{Name: "CompleteRelatedKey"}, completeRelatedKey(p))
			p.HandleFunction(&plugin.FunctionOptions{Name: "CompleteKey"}, completeKey(p))
			p.HandleFunction(&plugin.FunctionOptions{Name: "CompleteOpen"}, completeOpen(p))

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(neovimCmd)
}
