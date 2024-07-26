package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/expr-lang/expr"
	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/github"
	"github.com/tasshi-me/gh-iteration/pkg/log"
)

type ItemsEditProps struct {
	OutputFormatJSON *bool
}

type ItemsEditOption struct {
	ProjectOwner   string
	ProjectNumber  int
	FieldName      string
	Query          string
	Clear          bool
	Current        bool
	IterationTitle string
	DryRun         bool
}

func NewItemsEditCmd(props *ItemsEditProps) *cobra.Command {
	opts := new(ItemsEditOption)

	// itemsEditCmd represents the items-edit command.
	itemsEditCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "items-edit",
		Short: "Edit iteration of multiple project items",
		Long:  `Edit iteration of multiple project items`,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			itemsEditRun(props, opts)
		},
	}

	itemsEditCmd.Flags().SortFlags = false
	itemsEditCmd.Flags().IntVar(&opts.ProjectNumber, "project", 0, "Project number")
	itemsEditCmd.Flags().StringVar(&opts.ProjectOwner, "owner", "", "User/Organization login name")
	itemsEditCmd.Flags().StringVar(&opts.Query, "query", "false", "Query to filter target project items")
	itemsEditCmd.Flags().StringVar(&opts.FieldName, "field", "", "Iteration field name")
	itemsEditCmd.Flags().BoolVar(&opts.Clear, "clear", false, "Clear iteration field value")
	itemsEditCmd.Flags().BoolVar(&opts.Current, "current", false, "Set current iteration as the iteration field value")
	itemsEditCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "DryRun mode")
	itemsEditCmd.Flags().StringVar(&opts.IterationTitle, "iteration", "", "Iteration title to set")
	itemsEditCmd.MarkFlagsOneRequired("clear", "current", "iteration")
	_ = itemsEditCmd.MarkFlagRequired("project")
	_ = itemsEditCmd.MarkFlagRequired("owner")
	_ = itemsEditCmd.MarkFlagRequired("query")
	_ = itemsEditCmd.MarkFlagRequired("field")

	return itemsEditCmd
}

//nolint:funlen,gocognit,cyclop,gocyclo
func itemsEditRun(props *ItemsEditProps, opts *ItemsEditOption) {
	program, err := expr.Compile(opts.Query, expr.AsBool())
	if err != nil {
		log.Error(fmt.Errorf("failed to compile input query: %w", err))
		os.Exit(1)
	}

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
	iterationField, err := github.FetchIterationFieldByName(project.ID, opts.FieldName)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve an iteration by field name and project: %w", err))
		os.Exit(1)
	}

	log.Debug("Retrieve project items")
	githubItems, err := github.FetchProjectItems(project.ID)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve a project items by item id: %w", err))
		os.Exit(1)
	}

	items := make([]ProjectItem, 0, len(*githubItems))
	for _, githubItem := range *githubItems {
		item := ConvertGitHubProjectItem(&githubItem)
		items = append(items, item)
	}

	for _, item := range items {
		log.Debug("Item name: " + item.Title)

		output, err := expr.Run(program, map[string]any{
			"Item": item,
		})
		if err != nil {
			log.Error(fmt.Errorf("failed run query: %w", err))
			os.Exit(1)
		}
		pass, ok := output.(bool)
		if !ok {
			log.Error(errors.New("the result of query is not bool value. Please update the query"))
			os.Exit(1)
		}

		log.Debug(fmt.Sprintf("query result: %t", pass))

		if !pass {
			continue
		}

		skipped := false

		iterationIDFromCurrentItem := ""
		if item.Fields[opts.FieldName] != nil {
			iterationFromCurrentItem, ok := item.Fields[opts.FieldName].(FieldIteration)
			if ok {
				iterationIDFromCurrentItem = iterationFromCurrentItem.IterationID
			}
		}

		switch {
		case opts.Current:
			{
				log.Debug("Update iteration field to current sprint")
				currentIteration := iterationField.Configuration.Iterations[0]
				log.Debug("current sprint: " + currentIteration.Title)

				if iterationIDFromCurrentItem == currentIteration.ID {
					log.Debug("No need to update. Skip.")
					skipped = true
				} else if !opts.DryRun {
					_, err = github.UpdateIterationField(project.ID, iterationField.ID, item.ID, currentIteration.ID)
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
					} else if !opts.DryRun {
						_, err = github.UpdateIterationField(project.ID, iterationField.ID, item.ID, iteration.ID)
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
				} else if !opts.DryRun {
					_, err = github.ClearIterationField(project.ID, iterationField.ID, item.ID)
				}
			}
		}

		if err != nil {
			log.Error(fmt.Errorf("failed to update an iteration field: %w", err))
			os.Exit(1)
		}

		result := struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Skipped bool   `json:"skipped"`
			DryRun  bool   `json:"dryRun"`
		}{ID: item.ID, Title: item.Title, Skipped: skipped, DryRun: opts.DryRun}

		if *props.OutputFormatJSON {
			bytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Error(fmt.Errorf("failed to marshal result: %w", err))
				os.Exit(1)
			}
			_, _ = fmt.Fprintln(os.Stdout, string(bytes))
		} else {
			switch {
			case skipped:
				_, _ = fmt.Fprintf(os.Stdout, "%s %s => No need to update. Skipped.\n", item.ID, item.Title)
			case opts.DryRun:
				_, _ = fmt.Fprintf(os.Stdout, "%s %s => DryRun.\n", item.ID, item.Title)
			default:
				_, _ = fmt.Fprintf(os.Stdout, "%s %s => Updated.\n", item.ID, item.Title)
			}
		}
	}
}
