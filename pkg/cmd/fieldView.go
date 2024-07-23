package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/flags"
	"github.com/tasshi-me/gh-iteration/pkg/github"
	"github.com/tasshi-me/gh-iteration/pkg/log"
)

type FieldViewProps struct {
	OutputFormatJSON *bool
}

type FieldViewOption struct {
	ProjectOwner  string
	ProjectNumber int
	FieldName     string
}

func NewFieldViewCmd(props *FieldViewProps) *cobra.Command {
	opts := new(FieldViewOption)

	// fieldListCmd represents the field-list command.
	fieldListCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "field-view",
		Short: "View an iteration field",
		Long:  `View an iteration field`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			validator := flags.NewValidator(
				flags.And(
					flags.Flag("field"),
					flags.Flag("project"),
					flags.Flag("owner"),
				),
			)
			err := validator.Validate(cmd)
			if err != nil {
				return fmt.Errorf("flags: %w", err)
			}
			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			fieldViewRun(props, opts)
		},
	}

	fieldListCmd.Flags().SortFlags = false
	fieldListCmd.Flags().StringVar(&opts.FieldName, "field", "", "Iteration field name")
	fieldListCmd.Flags().IntVar(&opts.ProjectNumber, "project", 0, "Project number")
	fieldListCmd.Flags().StringVar(&opts.ProjectOwner, "owner", "", "User/Organization login name")
	_ = fieldListCmd.MarkFlagRequired("field")
	_ = fieldListCmd.MarkFlagRequired("project")
	_ = fieldListCmd.MarkFlagRequired("owner")

	return fieldListCmd
}

func fieldViewRun(props *FieldViewProps, opts *FieldViewOption) {
	log.Debug("Retrieve owner by login name")
	projectOwner, err := github.FetchOwnerByLogin(opts.ProjectOwner)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve owner by owner login: %w", err))
		os.Exit(1)
	}
	log.Debug("Owner: " + projectOwner.Login)

	log.Debug("Retrieve project by owner and project number")
	project, err := github.FetchProjectByNumber(opts.ProjectNumber, projectOwner.ID)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve a project by project number: %w", err))
		os.Exit(1)
	}
	log.Debug("Project ID: " + project.ID)

	log.Debug("Retrieve an iteration field by field name and project")
	field, err := github.FetchIterationFieldByName(project.ID, opts.FieldName)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve an iteration by field name and project: %w", err))
		os.Exit(1)
	}

	if *props.OutputFormatJSON {
		bytes, err := json.MarshalIndent(field, "", "  ")
		if err != nil {
			log.Error(fmt.Errorf("failed to marshal iterations: %w", err))
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, string(bytes))
	} else {
		s := formatIterationFieldPlain(field)
		_, _ = fmt.Fprint(os.Stdout, s)
	}
}

func formatIterationFieldPlain(field *github.ProjectV2IterationField) string {
	currentIteration := field.Configuration.Iterations[0]

	maxFieldNameLen := len(field.Name)
	maxFieldIDLen := len(field.ID)
	maxTitleLen := len(currentIteration.Title)

	format := "%-" + strconv.Itoa(maxFieldNameLen) + "s  %-" + strconv.Itoa(maxFieldIDLen) + "s  %-" + strconv.Itoa(maxTitleLen) + "s  %-10s\n"
	str := fmt.Sprintf(format, "Name", "ID", "Current", "StartDate")
	str += fmt.Sprintf(format, field.Name, field.ID, currentIteration.Title, currentIteration.StartDate)
	return str
}
