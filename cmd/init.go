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

	"github.com/satnamSandhu2001/stackjet/pkg/initializer"
	"github.com/spf13/cobra"
)

var forceRecreateConfig bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize StackJet configuration and database",
	Long: `Initialize StackJet in your environment by setting up the necessary configuration files and database.

This command creates the required configuration structure and prepares StackJet for managing your applications.
Run this command once after installation to set up StackJet.

Examples:
  # Initialize StackJet with default settings
  stackjet init

  # Force recreate configuration (overwrites existing config)
  stackjet init --force

After initialization, you can add your first application:
  stackjet add --tech nodejs --port 3000 --repo https://github.com/username/my-app.git`,
	Run: func(cmd *cobra.Command, args []string) {
		initializer.InitializeApp(forceRecreateConfig)
		fmt.Println("✅ StackJet initialized successfully.")
		fmt.Print("\nRun \033[1;34mstackjet add --tech nodejs --port 3000 --repo <git repo url>\033[0m to add new app.\n\n")
		fmt.Print("\nOr \033[1;34mstackjet add --help\033[0m for more information.\n\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&forceRecreateConfig, "force", "f", false, "Force recreate config. Use with caution! This will overwrite existing config and and remove all apps data from StackJet (except app folders)")
}
