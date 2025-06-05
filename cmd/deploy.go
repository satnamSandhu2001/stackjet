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
	"fmt"

	"github.com/satnamSandhu2001/stackjet/internal/cli/git"
	"github.com/satnamSandhu2001/stackjet/internal/cli/workspace"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/spf13/cobra"
)

// flags
var (
	dir       string
	gitBranch string
	gitRemote string
	gitReset  bool
	gitSkip   bool
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

Stackjet works via CLI, webhook triggers, or a web panel — making deployments simple, repeatable, and reliable.`,

	Run: func(cmd *cobra.Command, args []string) {
		// workspace logic
		if err := workspace.EnterWorkspace(dir); err != nil {
			return
		}

		// git logic
		if err := git.UpdateRepo(gitBranch, gitRemote, gitReset, gitSkip, gitHash); err != nil {
			return
		}

		fmt.Println("\n✅ Deployment complete!")

	}}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&dir, "dir", "d", "./", "Root directory of project")
	deployCmd.Flags().StringVarP(&gitBranch, "branch", "b", pkg.Config().GIT_BRANCH, "Git branch name")
	deployCmd.Flags().StringVar(&gitRemote, "git-remote", pkg.Config().GIT_REMOTE, "Git remote name")
	deployCmd.Flags().StringVar(&gitHash, "git-hash", "", "Rollback to specific commit hash")
	deployCmd.Flags().BoolVar(&gitReset, "git-reset", pkg.Config().GIT_RESET, "Force reset git state")
	deployCmd.Flags().BoolVar(&gitSkip, "git-skip", false, "Skip git repo update")
}
