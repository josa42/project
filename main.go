/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package main

import (
	"os"

	"github.com/josa42/project/cmd"
	"github.com/josa42/project/pkg/logger"
)

func main() {
	if lf := os.Getenv("PROJECT_LOG_FILE"); lf != "" {
		defer logger.InitLogger(lf)()
	}

	cmd.Execute()
}
