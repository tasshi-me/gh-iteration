package github

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// Project
// https://docs.github.com/en/graphql/reference/objects#projectv2
type Project struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func FetchProjectByNumber(number int, ownerID string) (*Project, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2Owner struct {
				ProjectV2 Project `graphql:"projectV2(number: $number)"`
			} `graphql:"... on ProjectV2Owner"`
		} `graphql:"node(id: $owner_id)"`
	}
	variables := map[string]interface{}{
		"owner_id": graphql.ID(ownerID),
		"number":   graphql.Int(number),
	}

	err = client.Query("Project", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a project: %w", err)
	}
	return &query.Node.ProjectV2Owner.ProjectV2, nil
}

func FetchProjectByID(projectID string) (*Project, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2 Project `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $project_id)"`
	}
	variables := map[string]interface{}{
		"project_id": graphql.ID(projectID),
	}

	err = client.Query("ProjectId", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ProjectV2 by project id: %w", err)
	}
	return &query.Node.ProjectV2, nil
}

type ProjectV2FieldConfiguration struct {
	ProjectV2FieldCommon struct {
		Name     string
		DataType string
	} `graphql:"... on ProjectV2FieldCommon"`

	// ProjectV2IterationField    struct{} `graphql:"... on ProjectV2IterationField"`
	// ProjectV2SingleSelectField struct{} `graphql:"... on ProjectV2SingleSelectField"`
}

func FetchProjectFields(projectID string) (*[]ProjectV2FieldConfiguration, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2 struct {
				Fields struct {
					Nodes []ProjectV2FieldConfiguration `graphql:"nodes"`
				} `graphql:"fields(first: 100)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $project_id)"`
	}
	variables := map[string]interface{}{
		"project_id": graphql.ID(projectID),
	}

	err = client.Query("ProjectV2FieldConfiguration", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ProjectV2FieldConfiguration by project id: %w", err)
	}
	return &query.Node.ProjectV2.Fields.Nodes, nil
}

// ProjectV2Item
// https://docs.github.com/ja/graphql/reference/objects#projectv2item
type ProjectItem struct {
	Content     ProjectItemContent `json:"content"`
	ID          string             `json:"id"`
	FieldValues struct {
		Nodes []FieldValue `json:"nodes"`
	} `json:"fieldValues" graphql:"fieldValues(first: 100)"`
	IsArchived bool    `json:"isArchived"`
	Type       string  `json:"type"` // DRAFT_ISSUE, ISSUE, PULL_REQUEST, REDACTED
	Project    Project `json:"project"`
}

type FieldValue struct {
	ProjectV2ItemFieldValueCommon struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldValueCommon"`
	ProjectV2ItemFieldDateValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
		// Date  string
	} `graphql:"... on ProjectV2ItemFieldDateValue"`
	ProjectV2ItemFieldIterationValue struct {
		Field       ProjectV2FieldConfiguration `json:"field"`
		IterationID string                      `json:"iterationId"`
		StartDate   string                      `json:"startDate"`
		Duration    graphql.Int                 `json:"duration"`
		Title       string                      `json:"title"`
		TitleHTML   string                      `graphql:"titleHTML" json:"titleHtml"`
	} `graphql:"... on ProjectV2ItemFieldIterationValue"`
	ProjectV2ItemFieldLabelValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldLabelValue"`
	ProjectV2ItemFieldMilestoneValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldMilestoneValue"`
	ProjectV2ItemFieldNumberValue struct {
		Field  ProjectV2FieldConfiguration `json:"field"`
		Number graphql.Float               `json:"number"`
	} `graphql:"... on ProjectV2ItemFieldNumberValue"`
	ProjectV2ItemFieldPullRequestValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldPullRequestValue"`
	ProjectV2ItemFieldRepositoryValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldRepositoryValue"`
	ProjectV2ItemFieldReviewerValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldReviewerValue"`
	ProjectV2ItemFieldSingleSelectValue struct {
		Field    ProjectV2FieldConfiguration `json:"field"`
		Name     string                      `json:"name"`
		OptionID string                      `json:"optionId"`
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	ProjectV2ItemFieldTextValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldTextValue"`
	ProjectV2ItemFieldUserValue struct {
		Field ProjectV2FieldConfiguration `json:"field"`
	} `graphql:"... on ProjectV2ItemFieldUserValue"`
}

// https://docs.github.com/ja/graphql/reference/unions#projectv2itemcontent
type ProjectItemContent struct {
	DraftIssue  DraftIssue  `graphql:"...on DraftIssue"  json:"draftIssue"`
	Issue       Issue       `graphql:"...on Issue"       json:"issue"`
	PullRequest PullRequest `graphql:"...on PullRequest" json:"pullRequest"`
}

// https://docs.github.com/ja/graphql/reference/objects#draftissue
type DraftIssue struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// https://docs.github.com/ja/graphql/reference/objects#issue
type Issue struct {
	ID          string     `json:"id"`
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	Closed      bool       `json:"closed"`
	Repository  Repository `json:"repository"`
	State       string     `json:"state"`       // CLOSED, OPEN
	StateReason string     `json:"stateReason"` // COMPLETED, NOT_PLANNED, REOPENED
}

// https://docs.github.com/ja/graphql/reference/objects#pullrequest
type PullRequest struct {
	ID         string     `json:"id"`
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	Closed     bool       `json:"closed"`
	IsDraft    bool       `json:"isDraft"`
	Merged     bool       `json:"merged"`
	Repository Repository `json:"repository"`
	State      string     `json:"state"` // CLOSED, MERGED, OPEN
}

// https://docs.github.com/ja/graphql/reference/objects#repository
type Repository struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	// Owner         string `json:"owner"`
}

func FetchProjectItem(itemID string) (*ProjectItem, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2Item ProjectItem `graphql:"... on ProjectV2Item"`
		} `graphql:"node(id: $item_id)"`
	}
	variables := map[string]interface{}{
		"item_id": graphql.ID(itemID),
	}

	err = client.Query("ProjectItem", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ProjectV2Item by item id: %w", err)
	}

	return &query.Node.ProjectV2Item, nil
}

func FetchProjectItems(projectID string) (*[]ProjectItem, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2 struct {
				Items struct {
					Nodes []ProjectItem `graphql:"nodes"`
				} `graphql:"items(first: 100)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $project_id)"`
	}
	variables := map[string]interface{}{
		"project_id": graphql.ID(projectID),
	}

	fields := map[string]interface{}{
		"hoge": 1,
	}
	fields["fuga"] = 2
	err = client.Query("ProjectItem", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ProjectV2 by project id: %w", err)
	}
	return &query.Node.ProjectV2.Items.Nodes, nil
}
