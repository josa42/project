/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
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
		keys := proj.RelatedKeys(filePath)

		return keys, nil
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

func alternate(p *plugin.Plugin) func(args []string) error {
	return func(args []string) error {

		target := "tab"
		if len(args) >= 1 {
			target = args[0]
		}

		forceKey := ""
		if len(args) >= 2 {
			forceKey = args[1]
		}

		proj := project.MustLoad(".")
		filePath := currentFilePath(p)

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
			return errors.New("command missing")
		}

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
			p.Nvim.Call("input", &key, "key: ", "", "customlist,CompleteRelatedKey")
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

// neovimCmd represents the neovim command
var neovimCmd = &cobra.Command{
	Use:    "neovim",
	Short:  "",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Main(func(p *plugin.Plugin) error {
			p.HandleFunction(&plugin.FunctionOptions{Name: "Alternate"}, alternate(p))
			p.HandleFunction(&plugin.FunctionOptions{Name: "CompleteRelatedKey"}, completeRelatedKey(p))

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(neovimCmd)
}
