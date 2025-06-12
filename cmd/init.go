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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initializer.InitializeApp(forceRecreateConfig)
		fmt.Println("✅ Stackjet initialized successfully.")
		fmt.Print("\nRun \033[1;34mstackjet add --tech nodejs -p 3000 -repo <git repo url>\033[0m to add new app.\n\n")
		fmt.Print("\nOr \033[1;34mstackjet add --help\033[0m for more information.\n\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&forceRecreateConfig, "force", "f", false, "Force recreate config. Use with caution! This will overwrite existing config.")
}
