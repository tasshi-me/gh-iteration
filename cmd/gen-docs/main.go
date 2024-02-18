package main

import (
	"os"

	"github.com/mshrtsr/gh-iteration/pkg/cmd"
	"github.com/mshrtsr/gh-iteration/pkg/docs"
	"github.com/mshrtsr/gh-iteration/pkg/log"
	"github.com/spf13/cobra"
)

func main() {
	command := NewCmd()
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type Options struct {
	Verbose          bool
	Trace            bool
	LogFormatJSON    bool
	OutputFormatJSON bool
	OutDir           string
}

func NewCmd() *cobra.Command {
	opts := new(Options)

	command := &cobra.Command{ //nolint:exhaustruct
		Use:   "gen-docs",
		Short: "Generate markdown document",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if opts.Verbose {
				log.SetLevel(log.ConfigLevelDebug)
			}
			if opts.Trace {
				log.SetLevel(log.ConfigLevelTrace)
			}
			if opts.LogFormatJSON {
				log.SetFormat(log.FormatJSON)
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			run(opts)
		},
	}

	command.Flags().SortFlags = false
	command.SetHelpCommand(&cobra.Command{Hidden: true}) //nolint:exhaustruct
	command.Flags().StringVar(&opts.OutDir, "out-dir", "", "Directory to generate documents")
	err := command.MarkFlagRequired("out-dir")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = command.MarkFlagDirname("out-dir")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	command.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Output verbose logs")
	command.Flags().BoolVar(&opts.Trace, "trace", false, "Output trace logs")
	command.Flags().BoolVar(&opts.LogFormatJSON, "log-json", false, "Output log in JSON")

	return command
}

func run(opts *Options) {
	targetCmd := cmd.NewRootCmd()
	info := docs.AdditionalInformation{
		Description:  "Hoge",
		Installation: "",
		Examples:     "",
		Links:        []docs.Link{{Label: "Repository on GitHub", URL: "https://github.com/mshrtsr/gh-iteration"}},
	}
	err := docs.GenMarkdownTree(targetCmd, info, opts.OutDir)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
