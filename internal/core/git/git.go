package git

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/logger"
)

// Verifies access to git repo
func VerifyAccess(w io.Writer, repoUrl string) error {
	repoUrl = strings.TrimSpace(repoUrl)
	logger.EmitLog(w, "")
	logger.EmitLog(w, "📡 Verifying Git Repo Access ...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"ls-remote", repoUrl}, Env: map[string]string{"GIT_TERMINAL_PROMPT": "0"}}); err != nil {
		return err
	}
	return nil
}

func CloneRepo(w io.Writer, gitRepo string, gitBranch string, gitRemote string) error {
	// trim whitespace from input strings
	gitRepo = strings.TrimSpace(gitRepo)
	gitBranch = strings.TrimSpace(gitBranch)
	gitRemote = strings.TrimSpace(gitRemote)

	// clone git repo
	logger.EmitLog(w, "")
	logger.EmitLog(w, "📡 Cloning Repository ...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"clone", "-b", gitBranch, "-o", gitRemote, gitRepo, "."}}); err != nil {
		return err
	}

	// switch to specified branch

	logger.EmitLog(w, fmt.Sprintf("⛓ Changing git branch to %v \n", gitBranch))
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"checkout", gitBranch}}); err != nil {
		return err
	}

	return nil
}

// Updates the local git-repo to specific version from remote-repo and returns error if failed
func UpdateRepo(w io.Writer, ctx context.Context, service services.StackService, deploymentID int64, gitBranch string, gitRemote string, gitReset bool, gitHash string) error {
	// trim whitespace from input strings
	gitBranch = strings.TrimSpace(gitBranch)
	gitRemote = strings.TrimSpace(gitRemote)
	gitHash = strings.TrimSpace(gitHash)

	// get current active branch
	activeBranch, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"branch", "--show-current"}})
	if err != nil {
		return err
	}

	// switch to specified branch if it's different from the current branch
	if strings.TrimSpace(activeBranch) != gitBranch {
		logger.EmitLog(w, "")
		logger.EmitLog(w, fmt.Sprintf("⛓ Changing git branch to %v \n", gitBranch))
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"checkout", gitBranch}}); err != nil {
			return err
		}
	}

	// fetch git status
	logger.EmitLog(w, "")
	logger.EmitLog(w, "🖇 Checking Git Status")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"fetch", "--all", "--tags"}}); err != nil {
		return err
	}

	// reset to specific commit if gitHash is provided
	if gitHash != "" {
		logger.EmitLog(w, "🎯 Resetting git to specific commit...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"reset", "--hard", gitHash}}); err != nil {
			return err
		}
		// update hash
		if err := updateHashToDB(w, ctx, service, deploymentID); err != nil {
			return err
		}
		return nil // no need to pull latest

	} else if gitReset { // force reset git state if gitReset is true
		logger.EmitLog(w, "")
		logger.EmitLog(w, "🧹 Forcing clean state with git reset...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"reset", "--hard", "origin/" + gitBranch}}); err != nil {
			return err
		}
	}

	// check if there are any commits behind the remote branch
	gitStatus, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"rev-list", "--count", fmt.Sprintf("HEAD...%s/%s", gitRemote, gitBranch)}})
	if err != nil {
		return err
	}
	logger.EmitLog(w, fmt.Sprintf("Commits behind remote: %s\n", strings.TrimSpace(gitStatus)))
	if strings.TrimSpace(gitStatus) == "0" {
		logger.EmitLog(w, "")
		logger.EmitLog(w, "✅ Repo Already up to date.")
		// update hash
		if err := updateHashToDB(w, ctx, service, deploymentID); err != nil {
			return err
		}
		return nil
	}

	// pull latest changes from remote branch
	logger.EmitLog(w, "")
	logger.EmitLog(w, "🔄 Pulling latest changes...")
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"pull", gitRemote, gitBranch}}); err != nil {
		return err
	}
	// update hash
	if err := updateHashToDB(w, ctx, service, deploymentID); err != nil {
		return err
	}

	return nil
}

// Fetch current hash from local repo and updates it to DB
func updateHashToDB(w io.Writer, ctx context.Context, service services.StackService, deploymentID int64) error {
	currentHash, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "git", Args: []string{"rev-parse", "HEAD"}})
	if err != nil {
		return err
	}
	// update hash to deployment
	updateDeploymentData := &dto.Deployment_Update_Request{
		ID:         deploymentID,
		CommitHash: strings.TrimSpace(currentHash),
	}
	if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
		return err
	}
	logger.EmitLog(w, "")
	logger.EmitLog(w, fmt.Sprintf("✅ Repo Updated to hash: %s\n", strings.TrimSpace(currentHash)))

	return nil
}
