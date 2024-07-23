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

type ListProps struct {
	OutputFormatJSON *bool
}

type ListOption struct {
	ProjectOwner  string
	ProjectNumber int
	FieldName     string
	Completed     bool
}

func NewListCmd(props *ListProps) *cobra.Command {
	opts := new(ListOption)

	// listCmd represents the list command.
	listCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "list",
		Short: "List the iterations for an iteration field",
		Long:  `List the iterations for an iteration field`,
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
			listRun(props, opts)
		},
	}

	listCmd.Flags().SortFlags = false
	listCmd.Flags().StringVar(&opts.FieldName, "field", "", "Iteration field name")
	listCmd.Flags().IntVar(&opts.ProjectNumber, "project", 0, "Project number")
	listCmd.Flags().StringVar(&opts.ProjectOwner, "owner", "", "User/Organization login name")
	listCmd.Flags().BoolVar(&opts.Completed, "completed", false, "List completed iterations")
	_ = listCmd.MarkFlagRequired("field")
	_ = listCmd.MarkFlagRequired("project")
	_ = listCmd.MarkFlagRequired("owner")

	return listCmd
}

func listRun(props *ListProps, opts *ListOption) {
	iterationField, err := retrieveIterationField(opts)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Debug("Iteration field ID: " + iterationField.ID)

	if *props.OutputFormatJSON {
		s, err := formatIterationFieldJSON(iterationField, opts.Completed)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, s)
	} else {
		s := formatIterationFieldPlain(iterationField, opts.Completed)
		_, _ = fmt.Fprint(os.Stdout, s)
	}
}

func retrieveIterationField(opts *ListOption) (*github.ProjectV2IterationField, error) {
	log.Debug("Retrieve owner by login name")
	projectOwner, err := github.FetchOwnerByLogin(opts.ProjectOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve owner by owner login: %w", err)
	}
	log.Debug("Owner: " + projectOwner.Login)

	log.Debug("Retrieve project by owner and project number")
	project, err := github.FetchProjectByNumber(opts.ProjectNumber, projectOwner.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a project by project number: %w", err)
	}
	log.Debug("Project ID: " + project.ID)

	log.Debug("Retrieve an iteration field by field name and project")
	i, err := github.FetchIterationFieldByName(project.ID, opts.FieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve an iteration by field name and project: %w", err)
	}
	return i, nil
}

func formatIterationFieldPlain(iterationField *github.ProjectV2IterationField, completed bool) string {
	var iterations []github.ProjectV2IterationFieldIteration
	if completed {
		iterations = iterationField.Configuration.CompletedIterations
	} else {
		iterations = iterationField.Configuration.Iterations
	}

	maxTitleLen := 0
	for _, iteration := range iterations {
		if l := len(iteration.Title); l > maxTitleLen {
			maxTitleLen = l
		}
	}

	str := fmt.Sprintf("%-"+strconv.Itoa(maxTitleLen)+"s  %-10s  %-8s  %-8s\n", "Title", "StartDate", "Duration", "ID")
	format := "%-" + strconv.Itoa(maxTitleLen) + "s  %-10s  %8d  %-8s\n"
	for _, iteration := range iterations {
		iter := fmt.Sprintf(format, iteration.Title, iteration.StartDate, iteration.Duration, iteration.ID)
		str += iter
	}
	return str
}

type JSONFormattedIterations struct {
	Iterations []JSONFormattedIteration `json:"iterations"`
}

type JSONFormattedIteration struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartDate string `json:"startDate"`
	Duration  int    `json:"duration"`
}

func formatIterationFieldJSON(iterationField *github.ProjectV2IterationField, completed bool) (string, error) {
	var iterations []github.ProjectV2IterationFieldIteration
	if completed {
		iterations = iterationField.Configuration.CompletedIterations
	} else {
		iterations = iterationField.Configuration.Iterations
	}

	iters := make([]JSONFormattedIteration, 0, len(iterations))
	for _, iteration := range iterations {
		iter := JSONFormattedIteration{
			ID:        iteration.ID,
			Title:     iteration.Title,
			StartDate: iteration.StartDate,
			Duration:  iteration.Duration,
		}
		iters = append(iters, iter)
	}
	obj := JSONFormattedIterations{Iterations: iters}

	bytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal iterations: %w", err)
	}
	return string(bytes), nil
}
