/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/josa42/project/pkg/project"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("create called")

		key := cmd.Flags().Lookup("template").Value.String()
		// license := cmd.Flags().Lookup("license").Value.String()

		pwd, _ := os.Getwd()

		var p *project.Template

		if key != "" {
			usr, _ := user.Current()
			tplPath := filepath.Join(usr.HomeDir, ".config", "project", "templates", key)
			p = project.LoadTemplate(tplPath)
		} else {
			p = project.DefaultTemplate()
		}

		p.Create(pwd)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("template", "t", "", "template")
	// createCmd.Flags().StringP("license", "l", "mit", "template")
	// createCmd.MarkFlagRequired("template")
}
