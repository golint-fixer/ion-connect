package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ion-channel/ionic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	teamID    string
	projectID string
)

func init() {
	ProjectCmd.AddCommand(GetProjectCmd)
	RootCmd.AddCommand(ProjectCmd)

	GetProjectCmd.Flags().StringVarP(&teamID, "team-id", "t", "", "--team-id <some team id>")
	GetProjectCmd.MarkFlagRequired("team-id")
	GetProjectCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "--project-id <some project id>")
	GetProjectCmd.MarkFlagRequired("project-id")
}

// ProjectCmd - Container for holding project root and secondary commands
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project anything to the screen",
	Long: `project is for printing anything back to the screen.
For many years people have printed back to the screen.`,
}

// GetProjectCmd - Container for holding project root and secondary commands
var GetProjectCmd = &cobra.Command{
	Use:   "get-project",
	Short: "Project anything to the screen",
	Long: `project is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		clionic, err := ionic.New(viper.GetString("endpoint_url"))
		if err != nil {
			panic(err)
		}
		project, err := clionic.GetProject(projectID, teamID, viper.GetString("secret_key"))
		if err != nil {
			panic(err)
		}
		bytes, err := json.Marshal(project)
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(bytes))
	},
}
