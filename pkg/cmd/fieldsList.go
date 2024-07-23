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

type FieldListProps struct {
	OutputFormatJSON *bool
}

type FieldListOption struct {
	ProjectOwner  string
	ProjectNumber int
}

func NewFieldListCmd(props *FieldListProps) *cobra.Command {
	opts := new(FieldListOption)

	// fieldListCmd represents the field-list command.
	fieldListCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "field-list",
		Short: "List the iteration fields in a project",
		Long:  `List the iteration fields in a project`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			validator := flags.NewValidator(
				flags.And(
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
			fieldListRun(props, opts)
		},
	}

	fieldListCmd.Flags().SortFlags = false
	fieldListCmd.Flags().IntVar(&opts.ProjectNumber, "project", 0, "Project number")
	fieldListCmd.Flags().StringVar(&opts.ProjectOwner, "owner", "", "User/Organization login name")
	_ = fieldListCmd.MarkFlagRequired("project")
	_ = fieldListCmd.MarkFlagRequired("owner")

	return fieldListCmd
}

func fieldListRun(props *FieldListProps, opts *FieldListOption) {
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
	fields, err := github.FetchIterationFields(project.ID)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve an iteration by field name and project: %w", err))
		os.Exit(1)
	}

	if *props.OutputFormatJSON {
		bytes, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			log.Error(fmt.Errorf("failed to marshal iterations: %w", err))
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, string(bytes))
	} else {
		s := formatIterationFieldsPlain(fields)
		_, _ = fmt.Fprint(os.Stdout, s)
	}
}

func formatIterationFieldsPlain(fields *[]github.ProjectV2IterationFieldWithoutConfiguration) string {
	maxFieldNameLen := 0
	for _, field := range *fields {
		if l := len(field.Name); l > maxFieldNameLen {
			maxFieldNameLen = l
		}
	}

	maxFieldIDLen := 0
	for _, field := range *fields {
		if l := len(field.ID); l > maxFieldIDLen {
			maxFieldIDLen = l
		}
	}

	str := fmt.Sprintf("%-"+strconv.Itoa(maxFieldNameLen)+"s  %-"+strconv.Itoa(maxFieldIDLen)+"s\n", "Name", "ID")
	format := "%-" + strconv.Itoa(maxFieldNameLen) + "s  %-" + strconv.Itoa(maxFieldIDLen) + "s\n"
	for _, field := range *fields {
		iter := fmt.Sprintf(format, field.Name, field.ID)
		str += iter
	}
	return str
}
