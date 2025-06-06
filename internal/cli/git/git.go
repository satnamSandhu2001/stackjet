package git

import (
	"fmt"
	"strings"

	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

// Updates the local git-repo to specific version from remote-repo and returns error if failed
func UpdateRepo(gitBranch string, gitRemote string, gitReset bool, gitSkip bool, gitHash string) error {
	// skip git repo update if gitSkip is true
	if gitSkip {
		fmt.Println("Skipping git repo update...")
		return nil
	}

	// trim whitespace from input strings
	gitBranch = strings.TrimSpace(gitBranch)
	gitRemote = strings.TrimSpace(gitRemote)
	gitHash = strings.TrimSpace(gitHash)

	// get current active branch
	activeBranch, err := commands.RunCommand("git", "branch", "--show-current")
	if err != nil {
		return err
	}

	// switch to specified branch if it's different from the current branch
	if strings.TrimSpace(activeBranch) != gitBranch {
		fmt.Printf("\n⛓ Changing git branch to %v \n", gitBranch)
		if _, err := commands.RunCommand("git", "checkout", gitBranch); err != nil {
			return err
		}
	}

	// fetch git status
	fmt.Println("\n🖇 Checking Git Status")
	if _, err := commands.RunCommand("git", "fetch", "--all", "--tags"); err != nil {
		return err
	}

	// reset to specific commit if gitHash is provided
	if gitHash != "" {
		fmt.Println("\n🎯 Resetting git to specific commit...")
		if _, err := commands.RunCommand("git", "reset", "--hard", gitHash); err != nil {
			return err
		}
		return nil // no need to pull latest

	} else if gitReset { // force reset git state if gitReset is true
		fmt.Println("\n🧹 Forcing clean state with git reset...")
		if _, err := commands.RunCommand("git", "reset", "--hard", "origin/"+gitBranch); err != nil {
			return err
		}
	}

	// check if there are any commits behind the remote branch
	gitStatus, err := commands.RunCommand("git", "rev-list", "--count", fmt.Sprintf("HEAD...%s/%s", gitRemote, gitBranch))
	if err != nil {
		return err
	}
	fmt.Printf("Commits behind remote: %s\n", strings.TrimSpace(gitStatus))
	if strings.TrimSpace(gitStatus) == "0" {
		fmt.Println("\n✅ Repo Already up to date.")
		return nil
	}

	// pull latest changes from remote branch
	fmt.Println("\n🔄 Pulling latest changes...")
	gitOutput, err := commands.RunCommand("git", "pull", gitRemote, gitBranch)
	if err != nil {
		return err
	}
	if strings.Contains(gitOutput, "Already up to date") {
		fmt.Println("\n✅ Repo Already up to date.")
		return nil
	}

	return nil
}
