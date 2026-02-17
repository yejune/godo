package main

import (
	"os"

	"github.com/do-focus/convert/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
