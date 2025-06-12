package stack

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/satnamSandhu2001/stackjet/internal/cli/git"
	"github.com/satnamSandhu2001/stackjet/internal/cli/nodejs"
	"github.com/satnamSandhu2001/stackjet/internal/cli/workspace"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/helpers"
)

func DeployStack(logger io.Writer, ctx context.Context, service services.StackService, opts *dto.Stack_DeployRequest) error {
	fmt.Fprintln(logger, "🛠️ Validating and preparing stack...")

	// get stack by directory
	stack, err := service.GetStackByDirectory(ctx, opts.Directory)
	if err != nil {
		return err
	}
	if stack.ID == 0 {
		return errors.New("stack not found")
	}

	// update git data if new data is provided
	if opts.Branch != "" || opts.Remote != "" {
		updateStackData := &dto.Stack_UpdateRequest{
			ID: stack.ID,
		}
		if opts.Branch != stack.Branch {
			updateStackData.Branch = opts.Branch
		}
		if opts.Remote != stack.Remote {
			updateStackData.Remote = opts.Remote
		}
		log.Println("updateStackData first", updateStackData)
		if updateStackData.Remote != "" || updateStackData.Branch != "" {
			log.Println("updateStackData", updateStackData)
			if err := service.UpdateStack(ctx, updateStackData); err != nil {
				return err
			}
		}
	}

	fmt.Fprintf(logger, "\n------ Deploying: %s ------\n", stack.Name)
	// add new deployment to deployments table
	updateDeploymentData := &dto.Deployment_CreateRequest{
		StackID: stack.ID,
		Status:  models.DEPLOYMENT_STATUS_IN_PROGRESS,
	}
	deploymentID, err := service.CreateDeployment(ctx, updateDeploymentData)
	if err != nil {
		return err
	}
	if deploymentID == 0 {
		return errors.New("failed to create deployment")
	}

	//  workspace logic
	if err := workspace.EnterWorkspace(logger, stack); err != nil {
		return err
	}

	// git logic
	if err := git.UpdateRepo(logger, ctx, service, deploymentID, stack.Branch, stack.Remote, opts.GitReset, opts.GitHash); err != nil {
		// update deployment status
		updateDeploymentData := &dto.Deployment_UpdateRequest{
			ID:     deploymentID,
			Status: models.DEPLOYMENT_STATUS_FAILED,
		}
		if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
			return err
		}

		return err
	}

	// nodejs + pm2 logic logic
	if stack.Type == "nodejs" {
		if err := nodejs.DeployStack(logger, ctx, service, stack); err != nil {
			// update deployment status
			updateDeploymentData := &dto.Deployment_UpdateRequest{
				ID:     deploymentID,
				Status: models.DEPLOYMENT_STATUS_FAILED,
			}
			if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
				return err
			}

			return err
		}
	}

	// TODO : update deploy table (Status) and logs to table

	// update deployment status
	updateStackStatusData := &dto.Stack_UpdateRequest{
		ID:                       stack.ID,
		InitialDeploymentSuccess: true,
	}
	if err := service.UpdateStack(ctx, updateStackStatusData); err != nil {
		return err
	}

	fmt.Fprintln(logger, "🎉 Stack deployed successfully!")
	return nil
}

// createNewStack creates new stack
func CreateNewStack(logger io.Writer, ctx context.Context, service services.StackService, opts *dto.Stack_CreateRequest) error {
	fmt.Fprintln(logger, "🛠️ Validating and preparing stack...")

	// validate stack type
	if !IsValidStackType(opts.Type) {
		return errors.New("invalid stack type. Valid types: " + strings.Join(pkg.Config().VALID_STACKS, ", "))
	}
	// validate git repo access
	if err := git.VerifyAccess(logger, opts.RepoUrl); err != nil {
		return err
	}
	// validate port
	if err := commands.ValidatePort(opts.Port); err != nil {
		return err
	}
	fmt.Fprintln(logger, "🚧 Creating new stack...")
	// create stack in db
	newStackID, err := service.CreateStack(ctx, opts)
	if err != nil {
		return err
	}
	// get newly created stack
	newStack, err := service.GetStackByID(ctx, newStackID)
	if err != nil {
		return err
	}
	fmt.Fprintln(logger, "📁 Creating stack directory...")
	// create stack folder in system
	if err := commands.CreateDir(newStack.Directory); err != nil {
		return err
	}
	//  enter workspace
	if err := workspace.EnterWorkspace(logger, newStack); err != nil {
		return err
	}
	// clone repo to directory
	if err := git.CloneRepo(logger, newStack.RepoUrl, newStack.Branch, newStack.Remote); err != nil {
		return err
	}
	// update stack created_successfully status in db
	if err := service.UpdateStack(ctx, &dto.Stack_UpdateRequest{ID: newStack.ID, CreatedSuccessfully: helpers.Bool(true)}); err != nil {
		return err
	}
	fmt.Fprintln(logger, "🎉 Stack created successfully!")
	if logger == os.Stdout {
		fmt.Fprintf(logger, "✅ Stack created successfully!\n    Run \x1b[34mstackjet deploy -d %s\x1b[0m to deploy \n    or \x1b[34mcd %s\x1b[0m then \x1b[34mstackjet deploy -d ./ \x1b[0m\n\n", newStack.Directory, newStack.Directory)
	}

	return nil
}

// IsValidStackType checks if stack type is valid from config file
func IsValidStackType(stack string) bool {
	return slices.Contains(pkg.Config().VALID_STACKS, stack)
}
