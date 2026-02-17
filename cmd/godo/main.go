package main

import (
	"os"

	"github.com/yejune/godo/internal/cli"
)

var version string

func main() {
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
