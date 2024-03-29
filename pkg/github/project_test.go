package github_test

import (
	"os"
	"testing"

	"github.com/tasshi-me/gh-iteration/pkg/github"
)

func TestFetchProjectByNumber(t *testing.T) {
	t.Parallel()
	if os.Getenv("CI") == "true" {
		t.Skip()
	}

	login := "tasshi-playground"
	projectNumber := 2
	owner, err := github.FetchOwnerByLogin(login)
	if err != nil {
		t.Fatal(err)
	}
	project, err := github.FetchProjectByNumber(projectNumber, owner.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(project.ID) == 0 {
		t.Errorf("failed to retrieve ID")
	}
	if len(project.Title) == 0 {
		t.Errorf("failed to retrieve Title")
	}
	if project.Number != projectNumber {
		t.Errorf("wrong owner type want: %d, got %d", projectNumber, project.Number)
	}
}
