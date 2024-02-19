package github

import (
	"errors"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type OwnerType int

const (
	OwnerTypeUser OwnerType = 1 << iota
	OwnerTypeOrganization
)

type Owner struct {
	ID    string    `json:"id"`
	Login string    `json:"login"`
	Name  string    `json:"name"`
	Type  OwnerType `json:"type"`
}

func FetchOwnerByLogin(login string) (*Owner, error) {
	if len(login) == 0 || login == "@me" {
		viewer, err := FetchUserByViewer()
		if err != nil {
			return nil, err
		}
		return &Owner{
			ID:    viewer.ID,
			Login: viewer.Login,
			Name:  viewer.Name,
			Type:  OwnerTypeUser,
		}, nil
	}

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Organization Organization `graphql:"organization(login: $login)"`
		User         User         `graphql:"user(login: $login)"`
	}
	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	err = client.Query("OrgOrUser", &query, variables)

	if len(query.User.Login) > 0 {
		return &Owner{
			ID:    query.User.ID,
			Login: query.User.Login,
			Name:  query.User.Name,
			Type:  OwnerTypeUser,
		}, nil
	}
	if len(query.Organization.Login) > 0 {
		return &Owner{
			ID:    query.Organization.ID,
			Login: query.Organization.Login,
			Name:  query.Organization.Name,
			Type:  OwnerTypeOrganization,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve owner by login: %w", err)
	}
	return nil, errors.New("failed to retrieve owner by login")
}

// Organization
// https://docs.github.com/ja/graphql/reference/objects#organization
type Organization struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func FetchOrganizationByLogin(login string) (*Organization, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Organization Organization `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	err = client.Query("Organization", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a organization: %w", err)
	}

	return &query.Organization, nil
}

// User
// https://docs.github.com/ja/graphql/reference/objects#user
type User struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

func FetchUserByLogin(login string) (*User, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		User User `graphql:"user(login: $login)"`
	}
	variables := map[string]interface{}{
		"login": graphql.String(login),
	}

	err = client.Query("User", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a user: %w", err)
	}

	return &query.User, nil
}

func FetchUserByViewer() (*User, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		User User `graphql:"viewer"`
	}

	err = client.Query("Viewer", &query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a user: %w", err)
	}

	return &query.User, nil
}
