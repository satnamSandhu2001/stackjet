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

	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/internal/cli/stack"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/spf13/cobra"
)

// flags
var (
	dir string

	gitBranch string
	gitRemote string
	gitReset  bool
	gitHash   string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your app with Git sync, NGINX, SSL, service restarts and more",
	Long: `Deploy your application end-to-end with automated Git synchronization and process management.

This command handles the complete deployment workflow:
  - Pulls the latest code from your Git repository
  - Executes build commands if specified
  - Manages application processes (PM2 for Node.js applications)
  - Executes post-deployment commands
  - Provides rollback capabilities to specific commits

The deployment works with applications previously added via 'stackjet add' command.
StackJet automatically detects the application configuration from the target directory.

Examples:
  # Deploy from current directory
  stackjet deploy

  # Deploy from specific directory
  stackjet deploy --dir "/var/www/sites/my-app"

  # Deploy specific branch
  stackjet deploy --branch "production"

  # Deploy with custom git remote
  stackjet deploy --git-remote "upstream" --branch "main"

  # Rollback to specific commit
  stackjet deploy --git-hash "abc123def456"

  # Deploy without git reset (preserve local changes)
  stackjet deploy --git-reset=false

Note: The directory must contain a StackJet-managed application (added via 'stackjet add').`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// set default values
		if !cmd.Flags().Changed("git-reset") {
			gitReset = pkg.Config().GIT_RESET
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		dbConn := database.Connect()
		defer dbConn.Close()
		stackService := services.NewStackService(dbConn)

		// deploy stack logic
		if err := stack.DeployStack(os.Stdout, context.Background(), *stackService, &dto.Stack_DeployRequest{
			Directory: dir,
			Remote:    gitRemote,
			Branch:    gitBranch,
		}); err != nil {
			fmt.Printf("⭕ Failed to deploy stack: %s\n", err)
			return
		}

	}}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&dir, "dir", "d", "./", "Root directory of the project to deploy")
	deployCmd.Flags().StringVar(&gitBranch, "branch", "", "Git branch name to deploy")
	deployCmd.Flags().StringVar(&gitRemote, "git-remote", "", "Git remote name (e.g., 'origin', 'upstream')")
	deployCmd.Flags().StringVar(&gitHash, "git-hash", "", "Rollback to specific commit hash")
	deployCmd.Flags().BoolVar(&gitReset, "git-reset", true, "Force reset Git state before deployment")

}
