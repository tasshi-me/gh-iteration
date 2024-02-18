package main

import (
	"log"

	"github.com/mshrtsr/gh-iteration/pkg/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
