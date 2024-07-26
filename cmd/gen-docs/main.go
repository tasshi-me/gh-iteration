package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/cmd"
	"github.com/tasshi-me/gh-iteration/pkg/docs"
	"github.com/tasshi-me/gh-iteration/pkg/log"
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

	description := "GitHub CLI extension to list/view/edit iterations in GitHub Projects.\n"

	installation := "" +
		"#### 1. Install GitHub CLI (`gh`)\n" +
		"\n" +
		"See the [document](https://github.com/cli/cli).\n" +
		"\n" +
		"#### 2. Install gh-iteration\n" +
		"\n" +
		"```shell\n" +
		"gh extension install tasshi-me/iteration\n" +
		"````\n" +
		"#### 3. Refresh token with `project` scope (if required)\n\n" +
		"To access to GitHub Projects, your token must have `project` scope.\n\n" +
		"You can check current your token scopes by:\n" +
		"```shell\n" +
		"$ gh auth status\n" +
		"```\n" +
		"If `project` is not in token scopes, add `project` scope by:\n" +
		"```shell\n" +
		"$ gh auth refresh -s project\n" +
		"```\n" +
		"\n"

	examples := "" +
		"```shell\n" +
		"# Set current iteration to the project items that match the query condition.\n" +
		"$ gh iteration items-edit \\\n" +
		"    --owner <OWNER> \\\n" +
		"    --project <PROJECT_NUM> \\\n" +
		"    --field <FIELD_NAME> \\\n" +
		"    --query \"Items.Fields.Status.Name == \\\"In progress\\\"\" \\\n" +
		"    --iteration-current\n" +
		"\n" +
		"```\n" +
		"\n"

	links := []docs.Link{
		{Label: "Repository", URL: "https://github.com/tasshi-me/gh-iteration"},
		{Label: "Release notes", URL: "https://github.com/tasshi-me/gh-iteration/releases"},
	}

	info := docs.AdditionalInformation{
		Description:  description,
		Installation: installation,
		Examples:     examples,
		Links:        links,
	}
	err := docs.GenMarkdownTree(targetCmd, info, opts.OutDir)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
