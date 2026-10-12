// Generate deterministic CLI reference pages directly from the command tree.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saurlax/dlut-pan-cli/internal/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	if err := os.MkdirAll("docs/commands", 0755); err != nil {
		panic(err)
	}
	root := cli.New("dev")
	root.DisableAutoGenTag = true
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	if err := doc.GenMarkdownTree(root, "docs/commands"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	files, err := filepath.Glob("docs/commands/*.md")
	if err != nil {
		panic(err)
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		// Cobra emits an extra blank line at EOF; keep generated files diff-clean.
		if err := os.WriteFile(path, []byte(strings.TrimRight(string(data), "\r\n")+"\n"), 0644); err != nil {
			panic(err)
		}
	}
}
