package github

import (
	"fmt"
	"strconv"

	"github.com/cli/go-gh/v2"
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

func FetchProjectItem(number int, owner string) (*Project, error) {
	//nolint:godox
	// TODO: Implementation
	projectID, stderr, err := gh.Exec(
		"project", "view", strconv.Itoa(number), "--owner", owner, "--format", "json", "--jq", ".id",
	)
	if err != nil {
		if len(stderr.String()) > 0 {
			return nil, fmt.Errorf("%s: %w", stderr.String(), err)
		}
		return nil, fmt.Errorf("failed to retrieve a project item: %w", err)
	}
	return &Project{ID: projectID.String(), Number: number}, nil //nolint:exhaustruct

	//
	// client, err := api.DefaultGraphQLClient()
	// if err != nil {{
	//	return nil, err
	//}
	//
	// var query struct {
	//	Organization struct {
	//		ProjectV2 Project `graphql:"projectV2(number: $number)"`
	//	} `graphql:"organization(login: $organization)"`
	//}
	// variables := map[string]interface{}{
	//	"organization": graphql.String(owner),
	//	"number":       graphql.Int(number),
	//}
	//
	// err = client.Query("ProjectId", &query, variables)
	// if err != nil {
	//	return nil, err
	//}
	// return &query.Organization.ProjectV2, nil
}
