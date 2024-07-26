package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/github"
	"github.com/tasshi-me/gh-iteration/pkg/log"
)

type ItemViewProps struct {
	OutputFormatJSON *bool
}

type ItemViewOption struct {
	ID string
}

func NewItemViewCmd(props *ItemViewProps) *cobra.Command {
	opts := new(ItemViewOption)

	// fieldViewCmd represents the field-view command.
	fieldViewCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "item-view",
		Short: "View a project item",
		Long:  `View a project item`,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			itemViewRun(props, opts)
		},
	}

	fieldViewCmd.Flags().SortFlags = false
	fieldViewCmd.Flags().StringVar(&opts.ID, "id", "", "ID of the project item to edit")
	_ = fieldViewCmd.MarkFlagRequired("id")

	return fieldViewCmd
}

func itemViewRun(props *ItemViewProps, opts *ItemViewOption) {
	log.Debug("Retrieve project item by ID")
	githubItem, err := github.FetchProjectItem(opts.ID)
	if err != nil {
		log.Error(fmt.Errorf("failed to retrieve a project item by item id: %w", err))
		os.Exit(1)
	}

	item := ConvertGitHubProjectItem(githubItem)
	log.Debug("Item name: " + item.Title)

	if *props.OutputFormatJSON {
		bytes, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			log.Error(fmt.Errorf("failed to marshal iterations: %w", err))
			os.Exit(1)
		}
		_, _ = fmt.Fprint(os.Stdout, string(bytes))
	} else {
		s := formatItemPlain(&item)
		_, _ = fmt.Fprint(os.Stdout, s)
	}
}

func formatItemPlain(item *ProjectItem) string {
	repoLen := max(len(item.Repository), len("Repo"))
	numLen := max(len(strconv.Itoa(item.Number)), len("Number"))
	IDLen := max(len(item.ID), len("ID"))
	titleLen := max(len(item.Title), len("Title"))

	//nolint:lll
	format := "%-" + strconv.Itoa(repoLen) + "s  %-" + strconv.Itoa(numLen) + "s  %-" + strconv.Itoa(IDLen) + "s  %-" + strconv.Itoa(titleLen) + "s\n"
	str := fmt.Sprintf(format, "Repo", "Number", "ID", "Title")
	str += fmt.Sprintf(format, item.Repository, strconv.Itoa(item.Number), item.ID, item.Title)
	return str
}

type ProjectItem struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Repository string                 `json:"repository"`
	Number     int                    `json:"number"`
	Fields     map[string]interface{} `json:"fields"`
	IsArchived bool                   `json:"isArchived"`
	Type       string                 `json:"type"` // DRAFT_ISSUE, ISSUE, PULL_REQUEST, REDACTED
}

type FieldCommon struct {
	FieldType string `json:"fieldType"`
}

type FieldIteration struct {
	FieldType   string      `json:"fieldType"`
	IterationID string      `json:"iterationId"`
	StartDate   string      `json:"startDate"`
	Duration    graphql.Int `json:"duration"`
	Title       string      `json:"title"`
	TitleHTML   string      `json:"titleHtml"`
}

type FieldNumber struct {
	FieldType string        `json:"fieldType"`
	Number    graphql.Float `json:"number"`
}

type FieldSingleSelect struct {
	FieldType string `json:"fieldType"`
	Name      string `json:"name"`
	OptionID  string `json:"optionId"`
}

func ConvertGitHubProjectItem(item *github.ProjectItem) ProjectItem {
	itemID := item.ID
	title := item.Content.DraftIssue.Title
	repository := item.Content.Issue.Repository.NameWithOwner
	number := item.Content.Issue.Number
	fields := map[string]interface{}{}
	for _, fieldValue := range item.FieldValues.Nodes {
		fieldName := fieldValue.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.Name
		fields[fieldName] = ConvertGitHubProjectItemField(fieldValue)
	}

	isArchived := item.IsArchived
	ttype := item.Type

	return ProjectItem{
		ID:         itemID,
		Title:      title,
		Repository: repository,
		Number:     number,
		Fields:     fields,
		IsArchived: isArchived,
		Type:       ttype,
	}
}

//nolint:cyclop
func ConvertGitHubProjectItemField(fieldValue github.FieldValue) interface{} {
	fieldType := fieldValue.ProjectV2ItemFieldValueCommon.Field.ProjectV2FieldCommon.DataType

	switch fieldType {
	case "ASSIGNEES":
		return FieldCommon{FieldType: fieldType}
	case "DATE":
		return FieldCommon{FieldType: fieldType}
	case "ITERATION":
		return FieldIteration{
			FieldType:   fieldType,
			IterationID: fieldValue.ProjectV2ItemFieldIterationValue.IterationID,
			StartDate:   fieldValue.ProjectV2ItemFieldIterationValue.StartDate,
			Duration:    fieldValue.ProjectV2ItemFieldIterationValue.Duration,
			Title:       fieldValue.ProjectV2ItemFieldIterationValue.Title,
			TitleHTML:   fieldValue.ProjectV2ItemFieldIterationValue.TitleHTML,
		}
	case "LABELS":
		return FieldCommon{FieldType: fieldType}
	case "LINKED_PULL_REQUESTS":
		return FieldCommon{FieldType: fieldType}
	case "MILESTONE":
		return FieldCommon{FieldType: fieldType}
	case "NUMBER":
		return FieldNumber{
			FieldType: fieldType,
			Number:    fieldValue.ProjectV2ItemFieldNumberValue.Number,
		}
	case "REPOSITORY":
		return FieldCommon{FieldType: fieldType}
	case "REVIEWERS":
		return FieldCommon{FieldType: fieldType}
	case "SINGLE_SELECT":
		return FieldSingleSelect{
			FieldType: fieldType,
			Name:      fieldValue.ProjectV2ItemFieldSingleSelectValue.Name,
			OptionID:  fieldValue.ProjectV2ItemFieldSingleSelectValue.OptionID,
		}
	case "TEXT":
		return FieldCommon{FieldType: fieldType}
	case "TITLE":
		return FieldCommon{FieldType: fieldType}
	case "TRACKED_BY":
		return FieldCommon{FieldType: fieldType}
	case "TRACKS":
		return FieldCommon{FieldType: fieldType}
	default:
		return FieldCommon{FieldType: fieldType}
	}
}
