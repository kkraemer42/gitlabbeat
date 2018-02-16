package main

import (
	"os"

	"github.com/kkraemer42/gitlabbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
