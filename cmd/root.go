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
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackjet",
	Short: "Forget shell scripts - deploy with StackJet",
	Long: `StackJet is a modern deployment automation tool built for developers managing multi-framework applications.

With StackJet, you can:
  - Pull and sync code from Git repositories (with optional rollback)
  - Restart required services for your tech stack (PM2, systemd, supervisor, etc.)
  - Handle builds, migrations, and runtime prep (Spring Boot, Django, Laravel, Go, Node.js)
  - Update and reload NGINX configurations automatically
  - Manage SSL certificates and HTTPS setups
  - Sync DNS records and proxy rules via Cloudflare
  - Trigger deployments via CLI, webhook, or web panel

Whether it's a Spring Boot JAR, a Go binary, Django with Gunicorn, npm scripts or Laravel app with Artisan — StackJet ensures consistent, reliable deployment every time.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}

func init() {

}

func printBanner() {
	fmt.Printf(`%s%s
	╭─────────────────────────╮
	│     🚀 StackJet 🚀      │
	╰─────────────────────────╯
	   Deploy Fast, Fly High!
%s%s`, "\033[1m", "\033[32m", "\033[0m", "\n\n")
}
