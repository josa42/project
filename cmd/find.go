/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"fmt"

	"github.com/josa42/project/pkg/project"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		proj := project.MustLoad(".")
		key := args[0]

		for _, f := range proj.FindFiles(key) {
			fmt.Println(f)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(findCmd)
}
