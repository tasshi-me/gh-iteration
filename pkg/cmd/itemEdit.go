package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/github"
	"github.com/tasshi-me/gh-iteration/pkg/log"
)

type ItemEditProps struct {
	OutputFormatJSON *bool
}

type ItemEditOption struct {
	FieldName      string
	ID             string
	Clear          bool
	Current        bool
	IterationTitle string
}

func NewItemEditCmd(props *ItemEditProps) *cobra.Command {
	opts := new(ItemEditOption)

	// fieldEditCmd represents the field-edit command.
	fieldEditCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "item-edit",
		Short: "Edit iteration of a project item",
		Long:  `Edit iteration of a project item`,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			itemEditRun(props, opts)
		},
	}

	fieldEditCmd.Flags().SortFlags = false
	fieldEditCmd.Flags().StringVar(&opts.ID, "id", "", "ID of the project item to edit")
	fieldEditCmd.Flags().StringVar(&opts.FieldName, "field", "", "Iteration field name")
	fieldEditCmd.Flags().BoolVar(&opts.Clear, "clear", false, "Clear iteration field value")
	fieldEditCmd.Flags().BoolVar(&opts.Current, "current", false, "Set current iteration as the iteration field value")
	fieldEditCmd.Flags().StringVar(&opts.IterationTitle, "iteration", "", "Iteration title to set")
	fieldEditCmd.MarkFlagsOneRequired("clear", "current", "iteration")
	_ = fieldEditCmd.MarkFlagRequired("id")
	_ = fieldEditCmd.MarkFlagRequired("field")

	return fieldEditCmd
}

//nolint:funlen,gocognit,cyclop
func itemEditRun(props *ItemEditProps, opts *ItemEditOption) {
	log.Debug("Retrieve project item by ID")
	githubItem, err := github.FetchProjectItem(opts.ID)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve a project item by item id: %w", err))
		os.Exit(1)
	}

	project := &githubItem.Project
	log.Debug("Project ID: " + project.ID)

	item := ConvertGitHubProjectItem(githubItem)
	log.Debug("Item name: " + item.Title)

	log.Debug("Retrieve an iteration field by field name and project")
	iterationField, err := github.FetchIterationFieldByName(project.ID, opts.FieldName)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve an iteration by field name and project: %w", err))
		os.Exit(1)
	}

	var updatedID string
	skipped := false

	iterationIDFromCurrentItem := ""
	if item.Fields[opts.FieldName] != nil {
		iterationFromCurrentItem, ok := item.Fields[opts.FieldName].(FieldIteration)
		if ok {
			iterationIDFromCurrentItem = iterationFromCurrentItem.IterationID
		}
	}

	log.Debug("Update an iteration field")
	switch {
	case opts.Current:
		{
			log.Debug("Update iteration field to current sprint")
			currentIteration := iterationField.Configuration.Iterations[0]
			log.Debug("current sprint: " + currentIteration.Title)

			if iterationIDFromCurrentItem == currentIteration.ID {
				log.Debug("No need to update. Skip.")
				skipped = true
			} else {
				updatedID, err = github.UpdateIterationField(project.ID, iterationField.ID, item.ID, currentIteration.ID)
			}
		}
	case len(opts.IterationTitle) > 0:
		{
			log.Debug("Update iteration field to specified sprint:" + opts.IterationTitle)
			var iteration github.ProjectV2IterationFieldIteration
			for _, iter := range iterationField.Configuration.Iterations {
				if iter.Title == opts.IterationTitle {
					iteration = iter
				}
			}
			for _, iter := range iterationField.Configuration.CompletedIterations {
				if iter.Title == opts.IterationTitle {
					iteration = iter
				}
			}
			if len(iteration.ID) > 0 {
				if iterationIDFromCurrentItem == iteration.ID {
					log.Debug("No need to update. Skip.")
					skipped = true
				} else {
					updatedID, err = github.UpdateIterationField(project.ID, iterationField.ID, item.ID, iteration.ID)
				}
			} else {
				err = fmt.Errorf("cannot find specified iteration: %s", opts.IterationTitle)
			}
		}
	case opts.Clear:
		{
			log.Debug("Clear iteration field")
			if len(iterationIDFromCurrentItem) == 0 {
				log.Debug("No need to update. Skip.")
				skipped = true
			} else {
				updatedID, err = github.ClearIterationField(project.ID, iterationField.ID, item.ID)
			}
		}
	}

	if err != nil {
		log.Error(fmt.Errorf("failed to update an iteration field: %w", err))
		os.Exit(1)
	}

	result := struct {
		ID      string `json:"id"`
		Skipped bool   `json:"skipped"`
	}{ID: updatedID, Skipped: skipped}

	if *props.OutputFormatJSON {
		bytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Error(fmt.Errorf("failed to marshal result: %w", err))
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, string(bytes))
	} else {
		if skipped {
			_, _ = fmt.Fprint(os.Stdout, "No need to update. Skipped.")
		} else {
			_, _ = fmt.Fprint(os.Stdout, result.ID)
		}
	}
}
