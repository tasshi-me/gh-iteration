package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/log"
)

type RootOptions struct {
	Verbose          bool
	Trace            bool
	LogFormatJSON    bool
	OutputFormatJSON bool
}

func NewRootCmd() *cobra.Command {
	opts := new(RootOptions)

	// rootCmd represents the base command when called without any subcommands.
	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "gh iteration",
		},
		Use:   "gh-iteration",
		Short: "Work with iteration fields of GitHub Projects",
		Long: `Work with iteration fields of GitHub Projects.
This command enables you to retrieve iteration field information or update an iteration of a project item.

To run commands, your token should have 'project' scope.
To verify your token scope, run 'gh auth status'.
To add the 'project' scope, run 'gh auth refresh -s project'.
`,
		Args: cobra.NoArgs,
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
	}

	rootCmd.DisableAutoGenTag = true
	rootCmd.Flags().SortFlags = false
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) //nolint:exhaustruct
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Output verbose logs")
	rootCmd.PersistentFlags().BoolVar(&opts.Trace, "trace", false, "[INTERNAL] Output trace logs")
	rootCmd.Flag("trace").Hidden = true
	rootCmd.PersistentFlags().BoolVar(&opts.LogFormatJSON, "log-json", false, "Output log in JSON")
	rootCmd.PersistentFlags().BoolVar(&opts.OutputFormatJSON, "json", false, "Output result in JSON")

	rootCmd.AddCommand(NewListCmd(&ListProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))
	rootCmd.AddCommand(NewFieldListCmd(&FieldListProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))
	rootCmd.AddCommand(NewFieldViewCmd(&FieldViewProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))
	rootCmd.AddCommand(NewItemViewCmd(&ItemViewProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))
	rootCmd.AddCommand(NewItemEditCmd(&ItemEditProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))
	rootCmd.AddCommand(NewItemsEditCmd(&ItemsEditProps{
		OutputFormatJSON: &opts.OutputFormatJSON,
	}))

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewRootCmd()
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
