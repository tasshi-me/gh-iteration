package github_test

import (
	"testing"

	"github.com/mshrtsr/gh-iteration/pkg/github"
)

func TestFetchOwnerIDByLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		login     string
		ownerType github.OwnerType
	}{
		{"mshrtsr", github.OwnerTypeUser},
		{"tasshi-playground", github.OwnerTypeOrganization},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.login, func(t *testing.T) {
			t.Parallel()

			login := test.login
			owner, err := github.FetchOwnerByLogin(login)
			if err != nil {
				t.Fatal(err)
			}
			if len(owner.ID) == 0 {
				t.Errorf("failed to retrieve ID")
			}
			if len(owner.Name) == 0 {
				t.Errorf("failed to retrieve Name")
			}
			if owner.Login != login {
				t.Errorf("wrong login want: %s, got %s", login, owner.Login)
			}
			if owner.Type != test.ownerType {
				t.Errorf("wrong owner type want: %d, got %d", test.ownerType, owner.Type)
			}
		})
	}
}

func TestFetchOrganizationByLogin(t *testing.T) {
	t.Parallel()

	login := "tasshi-playground"
	org, err := github.FetchOrganizationByLogin(login)
	if err != nil {
		t.Fatal(err)
	}
	if len(org.ID) == 0 {
		t.Errorf("failed to retrieve ID")
	}
	if len(org.Name) == 0 {
		t.Errorf("failed to retrieve Name")
	}
	if org.Login != login {
		t.Errorf("wrong login want: %s, got %s", login, org.Login)
	}
}

func TestFetchUserByLogin(t *testing.T) {
	t.Parallel()

	login := "mshrtsr"
	org, err := github.FetchUserByLogin(login)
	if err != nil {
		t.Fatal(err)
	}
	if len(org.ID) == 0 {
		t.Errorf("failed to retrieve ID")
	}
	if len(org.Name) == 0 {
		t.Errorf("failed to retrieve Name")
	}
	if org.Login != login {
		t.Errorf("wrong login want: %s, got %s", login, org.Login)
	}
}

func TestFetchUserByViewer(t *testing.T) {
	t.Parallel()

	login := "mshrtsr"
	org, err := github.FetchUserByViewer()
	if err != nil {
		t.Fatal(err)
	}
	if len(org.ID) == 0 {
		t.Errorf("failed to retrieve ID")
	}
	if len(org.Name) == 0 {
		t.Errorf("failed to retrieve Name")
	}
	if org.Login != login {
		t.Errorf("wrong login want: %s, got %s", login, org.Login)
	}
}
