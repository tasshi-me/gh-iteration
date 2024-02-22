package flags_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/tasshi-me/gh-iteration/pkg/flags"
)

type Flags struct {
	fieldID   bool
	field     bool
	projectID bool
	project   bool
	owner     bool
}

func setFlags(cmd *cobra.Command, flags Flags) {
	if flags.fieldID {
		cmd.Flag("field-id").Changed = true
	}
	if flags.field {
		cmd.Flag("field").Changed = true
	}
	if flags.projectID {
		cmd.Flag("project-id").Changed = true
	}
	if flags.project {
		cmd.Flag("project").Changed = true
	}
	if flags.owner {
		cmd.Flag("owner").Changed = true
	}
}

//nolint:funlen
func TestQuery(t *testing.T) {
	t.Parallel()
	// log.SetLevel(log.ConfigLevelDebug)

	newCmd := func() *cobra.Command {
		cmd := &cobra.Command{} //nolint:exhaustruct
		cmd.Flags().SortFlags = false
		cmd.Flags().String("field-id", "", "")
		cmd.Flags().String("field", "", "")
		cmd.Flags().String("project-id", "", "")
		cmd.Flags().Int("project", 0, "")
		cmd.Flags().String("owner", "", "")
		return cmd
	}

	tests := []struct {
		flags  Flags
		errMsg string
	}{
		{
			flags:  Flags{fieldID: false, field: false, projectID: false, project: false, owner: false},
			errMsg: "you must set one of [--field-id --field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: false, project: false, owner: true},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: false, project: true, owner: false},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: false, project: true, owner: true},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: true, project: false, owner: false},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: true, project: false, owner: true},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: true, project: true, owner: false},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: false, projectID: true, project: true, owner: true},
			errMsg: "when you set [--project-id --project --owner], you must set [--field]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: false, project: false, owner: false},
			errMsg: "when you set [--field], you must set [--project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: false, project: false, owner: true},
			errMsg: "when you set [--owner], you must set [--project]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: false, project: true, owner: false},
			errMsg: "when you set [--project], you must set [--owner]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: false, project: true, owner: true},
			errMsg: "",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: true, project: false, owner: false},
			errMsg: "",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: true, project: false, owner: true},
			errMsg: "when you set [--project-id], you cannot set [--project --owner]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: true, project: true, owner: false},
			errMsg: "when you set [--project-id], you cannot set [--project --owner]",
		},
		{
			flags:  Flags{fieldID: false, field: true, projectID: true, project: true, owner: true},
			errMsg: "when you set [--project-id], you cannot set [--project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: false, project: false, owner: false},
			errMsg: "",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: false, project: false, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: false, project: true, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: false, project: true, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: true, project: false, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: true, project: false, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: true, project: true, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: false, projectID: true, project: true, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: false, project: false, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: false, project: false, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: false, project: true, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: false, project: true, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: true, project: false, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: true, project: false, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: true, project: true, owner: false},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
		{
			flags:  Flags{fieldID: true, field: true, projectID: true, project: true, owner: true},
			errMsg: "when you set [--field-id], you cannot set [--field --project-id --project --owner]",
		},
	}

	query := flags.NewValidator(
		flags.Or(
			flags.Flag("field-id"),
			flags.And(
				flags.Flag("field"),
				flags.Or(
					flags.Flag("project-id"),
					flags.And(
						flags.Flag("project"), flags.Flag("owner"),
					),
				),
			),
		),
	)

	for i, tt := range tests {
		test := tt
		name := fmt.Sprintf("%d %+v", i, test.flags)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd := newCmd()
			setFlags(cmd, test.flags)

			err := query.Validate(cmd)
			if ((err == nil) && (len(test.errMsg) > 0)) ||
				((err != nil) && (len(test.errMsg) == 0)) ||
				((err != nil) && (err.Error() != test.errMsg)) {
				t.Errorf("Want %s, got %s", test.errMsg, err)
			}
		})
	}
}
