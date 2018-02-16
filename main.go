package main

import (
	"os"

	"github.com/kkraemer42/countbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
