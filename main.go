package main

import (
	"os"
	"strings"

	_ "embed"

	"github.com/jasonuc/greentext/cmd"
)

//go:embed version.txt
var version string

//go:embed templates/greentext_template.html
var defaultTemplate []byte

func main() {
	err := cmd.Execute(strings.Trim(version, "\n"), defaultTemplate)
	if err != nil {
		os.Exit(1)
	}
}
