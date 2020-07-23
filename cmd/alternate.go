/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"fmt"

	"github.com/josa42/project/pkg/project"
	"github.com/spf13/cobra"
)

// alternateCmd represents the alternate command
var alternateCmd = &cobra.Command{
	Use:     "alternate",
	Aliases: []string{"alt"},
	Short:   "",
	Long:    ``,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		proj := project.MustLoad(".")

		for _, f := range proj.AlternateFiles(args[0], args[1]) {
			fmt.Println(f)
		}
	},
}

func init() {
	rootCmd.AddCommand(alternateCmd)
}
