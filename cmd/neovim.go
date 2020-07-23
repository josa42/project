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

func alternate(p *plugin.Plugin) func(args []string) error {
	return func(args []string) error {
		proj := project.MustLoad(".")

		filePath, _ := p.Nvim.CommandOutput("echo expand('%')")

		if len(args) == 0 || !isAny(args[0], []string{"edit", "tabedit", "vsplit", "split"}) {
			return errors.New("command missing")
		}

		keys := proj.RelatedKeys(filePath)

		key := ""
		if len(args) == 2 && isAny(args[1], keys) {
			key = args[1]
		} else if len(keys) > 0 {
			key = keys[0]
		} else {
			return errors.New("key missing")
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

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(neovimCmd)
}
