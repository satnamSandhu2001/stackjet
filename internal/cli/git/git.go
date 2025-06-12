package git

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

// Verifies access to git repo
func VerifyAccess(logger io.Writer, repoUrl string) error {
	repoUrl = strings.TrimSpace(repoUrl)
	fmt.Fprintln(logger, "\n📡 Verifying Git Repo Access ...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"ls-remote", repoUrl}, Env: map[string]string{"GIT_TERMINAL_PROMPT": "0"}}); err != nil {
		return err
	}
	return nil
}

func CloneRepo(logger io.Writer, gitRepo string, gitBranch string, gitRemote string) error {
	// trim whitespace from input strings
	gitRepo = strings.TrimSpace(gitRepo)
	gitBranch = strings.TrimSpace(gitBranch)
	gitRemote = strings.TrimSpace(gitRemote)

	// clone git repo
	fmt.Fprintln(logger, "\n📡 Cloning Repository ...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"clone", "-b", gitBranch, "-o", gitRemote, gitRepo, "."}}); err != nil {
		return err
	}

	// switch to specified branch
	fmt.Fprintf(logger, "\n⛓ Changing git branch to %v \n", gitBranch)
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"checkout", gitBranch}}); err != nil {
		return err
	}

	return nil
}

// Updates the local git-repo to specific version from remote-repo and returns error if failed
func UpdateRepo(logger io.Writer, ctx context.Context, service services.StackService, deploymentID int64, gitBranch string, gitRemote string, gitReset bool, gitHash string) error {
	// trim whitespace from input strings
	gitBranch = strings.TrimSpace(gitBranch)
	gitRemote = strings.TrimSpace(gitRemote)
	gitHash = strings.TrimSpace(gitHash)

	// get current active branch
	activeBranch, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"branch", "--show-current"}})
	if err != nil {
		return err
	}

	// switch to specified branch if it's different from the current branch
	if strings.TrimSpace(activeBranch) != gitBranch {
		fmt.Fprintf(logger, "\n⛓ Changing git branch to %v \n", gitBranch)
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"checkout", gitBranch}}); err != nil {
			return err
		}
	}

	// fetch git status
	fmt.Fprintln(logger, "\n🖇 Checking Git Status")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"fetch", "--all", "--tags"}}); err != nil {
		return err
	}

	// reset to specific commit if gitHash is provided
	if gitHash != "" {
		fmt.Fprintln(logger, "\n🎯 Resetting git to specific commit...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"reset", "--hard", gitHash}}); err != nil {
			return err
		}
		return nil // no need to pull latest

	} else if gitReset { // force reset git state if gitReset is true
		fmt.Fprintln(logger, "\n🧹 Forcing clean state with git reset...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"reset", "--hard", "origin/" + gitBranch}}); err != nil {
			return err
		}
	}

	// check if there are any commits behind the remote branch
	gitStatus, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"rev-list", "--count", fmt.Sprintf("HEAD...%s/%s", gitRemote, gitBranch)}})
	if err != nil {
		return err
	}
	fmt.Fprintf(logger, "Commits behind remote: %s\n", strings.TrimSpace(gitStatus))
	if strings.TrimSpace(gitStatus) == "0" {
		fmt.Fprintln(logger, "\n✅ Repo Already up to date.")
		return nil
	}

	// pull latest changes from remote branch
	fmt.Fprintln(logger, "\n🔄 Pulling latest changes...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"pull", gitRemote, gitBranch}}); err != nil {
		return err
	}

	currentHash, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "git", Args: []string{"rev-parse", "HEAD"}})
	if err != nil {
		return err
	}
	// update hash to deployment
	updateDeploymentData := &dto.Deployment_UpdateRequest{
		ID:         deploymentID,
		CommitHash: strings.TrimSpace(currentHash),
	}
	if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
		return err
	}

	fmt.Fprintf(logger, "\n✅ Repo Updated to hash: %s\n", strings.TrimSpace(currentHash))

	return nil
}
