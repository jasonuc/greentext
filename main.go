package main

import (
	"os"
	"strings"

	_ "embed"

	"github.com/jasonuc/greentext/cmd"
)

//go:embed version.txt
var version string

func main() {
	err := cmd.Execute(strings.Trim(version, "\n"))
	if err != nil {
		os.Exit(1)
	}
}
