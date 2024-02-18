package cmd

import (
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/mshrtsr/gh-iteration/pkg/github"
	"github.com/mshrtsr/gh-iteration/pkg/log"
	"github.com/spf13/cobra"
)

type UpdateOption struct {
	ProjectOwner  string
	ProjectNumber int
	ProjectID     string
	FieldName     string
	FieldID       string
	Completed     bool
}

func NewUpdateCmd() *cobra.Command {
	opts := new(UpdateOption)

	// updateCmd represents the update command.
	updateCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "update",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(_ *cobra.Command, _ []string) {
			log.Info("update called")
			updateRun(opts)
		},
	}

	return updateCmd
}

func updateRun(opts *UpdateOption) {
	projectNumber := opts.ProjectNumber
	projectOwner := opts.ProjectOwner

	if len(projectOwner) == 0 {
		repo, err := repository.Current()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		projectOwner = repo.Owner
	}

	log.Debug("Owner: " + projectOwner)
	project, err := github.FetchProjectByNumber(projectNumber, projectOwner)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Debug("ProjectID: " + project.ID)

	iterationField, err := github.FetchIterationFieldByName(project.ID, "Sprint")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// if OutputFormatJSON {
	//	s, err := json.MarshalIndent(iterationField, "", "  ")
	//	if err != nil {
	//		log.Error(err)
	//		os.Exit(1)
	//	}
	//	fmt.Print(string(s))
	// } else {
	//	fmt.Print(iterationField)
	//}

	if len(iterationField.Configuration.Iterations) == 0 {
		log.Error("Current sprint not found")
		os.Exit(1)
	}
	currentSprint := iterationField.Configuration.Iterations[0]
	log.Debug(currentSprint)

	itemID := "PVTI_lADOAEl4zs4AMFWjzgLhmZ8"
	updatedID, err := github.UpdateIterationField(project.ID, iterationField.ID, itemID, currentSprint.ID)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Debug(updatedID)

	// fmt.Println(strconv.Itoa(projectNumber))
	// issueList, stderr, err := gh.Exec("project", "view", strconv.Itoa(projectNumber), "--owner", "kintone")
	// if err != nil {
	//	log.Println(stderr.String())
	//	log.Fatal(err)
	//}
	// fmt.Println(issueList.String())
}
