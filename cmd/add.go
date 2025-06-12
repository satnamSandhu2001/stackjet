/*
Copyright © 2025 Satnam Sandhu <satnamsandhu70002@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/internal/cli/stack"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/helpers"
	"github.com/spf13/cobra"
)

var (
	stackType    string
	repoUrl      string
	branch       string
	remote       string
	port         int
	buildCommand string
	startCommand string
	postCommand  string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new application to StackJet for deployment management",
	Long: `Add a new application to StackJet's deployment pipeline. This command registers your application
with StackJet and prepares it for automated deployments.

Currently supported technology stacks:
  - nodejs: Node.js applications with PM2 process management

Required information:
  - Technology stack type (--tech)
  - Git repository URL (--repo)
  - Application port (--port)

Optional customizations:
  - Custom build commands (--build)
  - Custom start commands (--start, defaults to "npm start" for Node.js)
  - Post-deployment commands (--post)
  - Git branch and remote settings

Examples:
  # Add a basic Node.js application
  stackjet add --tech nodejs --port 3000 --repo https://github.com/username/my-app.git

  # Add with custom commands
  stackjet add --tech nodejs --port 8080 --repo https://github.com/username/api.git \
    --build "npm install && npm run build" \
    --start "npm run prod" \
    --post "npm run migrate"

  # Add with specific branch
  stackjet add --tech nodejs --port 3000 --repo https://github.com/username/app.git \
    --branch production

After adding an application, deploy it with:
  stackjet deploy --dir /path/to/deployed/app`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		// validate stack flag
		if ok := stack.IsValidStackType(strings.TrimSpace(stackType)); !ok {
			return fmt.Errorf(`⭕ Invalid stack: "%s".
   Valid options are: "%s"`, stackType, strings.Join(pkg.Config().VALID_STACKS, `", "`))
		}
		if strings.TrimSpace(repoUrl) == "" {
			return fmt.Errorf("⭕ Git repository URL is required. Use -r or --git-repo to specify the repository URL")
		}
		if port == 0 {
			return fmt.Errorf("⭕ Port is required. Use -p or --port to specify the port")
		}
		if err := commands.ValidatePort(port); err != nil {
			return err
		}
		// validate start commands
		startCommand = strings.TrimSpace(startCommand)
		if startCommand != "" {
			switch stackType {
			case "nodejs":
				if err := helpers.ValidateNodeStartCommand(startCommand); err != nil {
					return err
				}
			}
		}
		if startCommand == "" {
			switch stackType {
			case "nodejs":
				startCommand = "npm start"
			}
		}

		// set default values
		if branch == "" {
			branch = pkg.Config().GIT_BRANCH
		}
		if remote == "" {
			remote = pkg.Config().GIT_REMOTE
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		dbConn := database.Connect()
		defer dbConn.Close()
		stackService := services.NewStackService(dbConn)

		// deploy stack logic
		appCommands := models.StackCommands{
			Start: startCommand,
		}
		if buildCommand != "" {
			appCommands.Build = buildCommand
		}
		if postCommand != "" {
			appCommands.Post = postCommand
		}

		if err := stack.CreateNewStack(os.Stdout, context.Background(), *stackService, &dto.Stack_CreateRequest{
			Type:     stackType,
			RepoUrl:  repoUrl,
			Branch:   branch,
			Remote:   remote,
			Port:     port,
			Commands: appCommands,
		}); err != nil {
			fmt.Printf("⭕ Failed to deploy stack: %s\n", err)
			return
		}

	},
}

func init() {

	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&stackType, "tech", "t", "", "App's Technology Stack Type (currently supports: nodejs)")
	addCmd.Flags().StringVarP(&repoUrl, "repo", "r", "", "Git repository URL")
	addCmd.Flags().IntVarP(&port, "port", "p", 0, "Port number for the application")
	addCmd.Flags().StringVar(&branch, "branch", "", "Git branch name (default master)")
	addCmd.Flags().StringVar(&remote, "git-remote", "", "Git remote name (default origin)")
	addCmd.Flags().StringVar(&buildCommand, "build", "", "Build commands (e.g. 'npm i && npm run build', 'mvn clean package', 'gradle build', etc...)")
	addCmd.Flags().StringVar(&startCommand, "start", "", "App start commands (e.g. 'npm start', 'mvn spring-boot:run', 'gradle bootRun', etc...)")
	addCmd.Flags().StringVar(&postCommand, "post", "", "Post deployment commands (e.g. 'npm run post-deploy', 'mvn post-deploy', 'gradle post-deploy', etc...)")

	// register auto completion for stack flag
	addCmd.RegisterFlagCompletionFunc("stack", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return pkg.Config().VALID_STACKS, cobra.ShellCompDirectiveNoFileComp
	})
}
