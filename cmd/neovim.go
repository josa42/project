/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/josa42/project/pkg/project"
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

func alternate(p *plugin.Plugin) func(args []string) error {
	return func(args []string) error {
		proj := project.MustLoad(".")

		filePath, _ := p.Nvim.CommandOutput("echo expand('%')")

		if len(args) == 0 || !isAny(args[0], []string{"edit", "e", "tabedit", "tabe", "vsplit", "split"}) {
			return errors.New("command missing")
		}

		key := ""

		keys := proj.RelatedKeys(filePath)
		if len(keys) == 1 {
			key = keys[0]
		} else if len(keys) > 1 {
			p.Nvim.Call("input", &key, "key: ", "", "customlist,CompleteRelatedKey")
		}

		if !isAny(key, keys) {
			return nil
		}

		for _, f := range proj.RelatedFiles(key, filePath) {
			return p.Nvim.Command(fmt.Sprintf(`%s %s`, args[0], f))
		}

		return nil
	}
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
