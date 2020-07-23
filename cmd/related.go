/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"fmt"

	"github.com/josa42/project/pkg/project"
	"github.com/spf13/cobra"
)

// relatedCmd represents the related command
var relatedCmd = &cobra.Command{
	Use:     "related",
	Aliases: []string{"alt", "rel"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		proj := project.MustLoad(".")

		for _, f := range proj.RelatedFiles(args[0], args[1]) {
			fmt.Println(f)
		}
	},
}

func init() {
	rootCmd.AddCommand(relatedCmd)
}
