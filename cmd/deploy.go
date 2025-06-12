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
	Long: `Deploy your application end-to-end with a single command.

This includes:
  - Pulling the latest code from Git (with optional rollback)
  - Restarting required services to bring your app live (e.g., pm2, systemd, etc.)
  - Logging and error handling with rollback to a previous version
  - Updating NGINX configuration
  - Managing SSL certificates
  - Syncing DNS and proxy settings via Cloudflare

StackJet works via CLI, webhook triggers, or a web panel — making deployments simple, repeatable, and reliable.`,
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

	deployCmd.Flags().StringVarP(&dir, "dir", "d", "./", "Root directory of project")
	deployCmd.Flags().StringVar(&gitBranch, "branch", "", "Git branch name")
	deployCmd.Flags().StringVar(&gitRemote, "git-remote", "", "Git remote name")
	deployCmd.Flags().StringVar(&gitHash, "git-hash", "", "Rollback to specific commit hash")
	deployCmd.Flags().BoolVar(&gitReset, "git-reset", true, "Force reset git state")

}
