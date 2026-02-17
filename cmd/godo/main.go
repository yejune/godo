package main

import (
	"os"

	"github.com/yejune/godo/internal/cli"
)

var version string

func normalizeLegacyAliases() {
	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "cc":
		os.Args[1] = "claude"
	case "self-update":
		os.Args[1] = "selfupdate"
	}

	if len(os.Args) >= 3 && os.Args[1] == "moai" && os.Args[2] == "rank" {
		os.Args = append([]string{os.Args[0], "rank"}, os.Args[3:]...)
	}
}

func main() {
	normalizeLegacyAliases()
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
