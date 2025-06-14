package stack

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/satnamSandhu2001/stackjet/internal/core/git"
	"github.com/satnamSandhu2001/stackjet/internal/core/nodejs"
	"github.com/satnamSandhu2001/stackjet/internal/core/workspace"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/helpers"
	"github.com/satnamSandhu2001/stackjet/pkg/logger"
)

// DeployStack deploys a stack and returns the deployment ID
func DeployStack(w io.Writer, ctx context.Context, service services.StackService, opts *dto.Stack_Deploy_Request) (int64, error) {
	logger.EmitLog(w, "🛠️ Validating and preparing stack...")

	// get stack
	var stack *models.Stack
	var err error
	if opts.Directory != "" {
		stack, err = service.GetStackByDirectory(ctx, opts.Directory) // for cli deploy
	} else {
		stack, err = service.GetStackByID(ctx, opts.ID) // for web deploy
	}
	if err != nil {
		return 0, err
	}
	if stack == nil {
		return 0, errors.New("stack not found")
	}
	if !stack.CreatedSuccessfully {
		return 0, errors.New("app was not created successfully. Please create app first")
	}
	// update git data if new data is provided
	if opts.Branch != "" || opts.Remote != "" {
		updateStackData := &dto.Stack_Update_Request{
			ID: stack.ID,
		}
		if opts.Branch != stack.Branch {
			updateStackData.Branch = opts.Branch
		}
		if opts.Remote != stack.Remote {
			updateStackData.Remote = opts.Remote
		}
		if updateStackData.Remote != "" || updateStackData.Branch != "" {
			if err := service.UpdateStack(ctx, updateStackData); err != nil {
				return 0, err
			}
		}
	}

	logger.EmitLog(w, "")
	logger.EmitLog(w, fmt.Sprintf("------ Deploying: %s ------\n", stack.Name))
	// add new deployment to deployments table
	updateDeploymentData := &dto.Deployment_Create_Request{
		StackID: stack.ID,
		Status:  models.DEPLOYMENT_STATUS_IN_PROGRESS,
	}
	deploymentID, err := service.CreateDeployment(ctx, updateDeploymentData)
	if err != nil {
		return 0, err
	}
	if deploymentID == 0 {
		return deploymentID, errors.New("failed to create deployment")
	}

	//  workspace logic
	if err := workspace.EnterWorkspace(w, stack); err != nil {
		return deploymentID, err
	}

	// git logic
	if err := git.UpdateRepo(w, ctx, service, deploymentID, stack.Branch, stack.Remote, opts.GitReset, opts.GitHash); err != nil {
		// update deployment status
		updateDeploymentData := &dto.Deployment_Update_Request{
			ID:     deploymentID,
			Status: models.DEPLOYMENT_STATUS_FAILED,
		}
		if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
			return deploymentID, err
		}

		return deploymentID, err
	}

	// nodejs + pm2 logic logic
	if stack.Type == "nodejs" {
		if err := nodejs.DeployStack(w, ctx, service, stack); err != nil {
			// update deployment status
			updateDeploymentData := &dto.Deployment_Update_Request{
				ID:     deploymentID,
				Status: models.DEPLOYMENT_STATUS_FAILED,
			}
			if _, err := service.UpdateDeployment(ctx, updateDeploymentData); err != nil {
				return deploymentID, err
			}

			return deploymentID, err
		}
	}

	// update stack success if deployed for the first time
	if !stack.InitialDeploymentSuccess {
		updateStackStatusData := &dto.Stack_Update_Request{
			ID:                       stack.ID,
			InitialDeploymentSuccess: helpers.Bool(true),
		}
		if err := service.UpdateStack(ctx, updateStackStatusData); err != nil {
			return deploymentID, err
		}
	}

	logger.EmitLog(w, "🎉 Stack deployed successfully!")

	// update deployment success
	updateDeploymentSuccess := &dto.Deployment_Update_Request{
		ID:     deploymentID,
		Status: models.DEPLOYMENT_STATUS_SUCCESS,
	}
	if _, err := service.UpdateDeployment(ctx, updateDeploymentSuccess); err != nil {
		return deploymentID, err
	}

	return deploymentID, nil
}

// createNewStack creates new stack
func CreateNewStack(w io.Writer, ctx context.Context, service services.StackService, opts *dto.Stack_Create_Request) error {
	logger.EmitLog(w, "🛠️ Validating and preparing stack...")

	// validate stack type
	if !IsValidStackType(opts.Type) {
		return errors.New("invalid stack type. Valid types: " + strings.Join(pkg.Config().VALID_STACKS, ", "))
	}

	// set default start command
	if opts.Commands.Start == "" {
		switch opts.Type {
		case "nodejs":
			opts.Commands.Start = "npm start"
		}
	}

	// validate git repo access
	if err := git.VerifyAccess(w, opts.RepoUrl); err != nil {
		return err
	}
	// validate port
	if err := commands.ValidatePort(opts.Port); err != nil {
		return err
	}
	logger.EmitLog(w, "🚧 Creating new stack...")
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
	logger.EmitLog(w, "📁 Creating stack directory...")
	// create stack folder in system
	if err := commands.CreateDir(newStack.Directory); err != nil {
		return err
	}
	//  enter workspace
	if err := workspace.EnterWorkspace(w, newStack); err != nil {
		return err
	}
	// clone repo to directory
	if err := git.CloneRepo(w, newStack.RepoUrl, newStack.Branch, newStack.Remote); err != nil {
		return err
	}
	// update stack created_successfully status in db
	if err := service.UpdateStack(ctx, &dto.Stack_Update_Request{ID: newStack.ID, CreatedSuccessfully: helpers.Bool(true)}); err != nil {
		return err
	}
	logger.EmitLog(w, "🎉 Stack created successfully!")
	if w == os.Stdout {
		logger.EmitLog(w, fmt.Sprintf("    Run \x1b[34mstackjet deploy -d %s\x1b[0m to deploy \n    or \x1b[34mcd %s\x1b[0m then \x1b[34mstackjet deploy -d ./ \x1b[0m\n\n", newStack.Directory, newStack.Directory))
	}

	return nil
}

// IsValidStackType checks if stack type is valid from config file
func IsValidStackType(stack string) bool {
	return slices.Contains(pkg.Config().VALID_STACKS, stack)
}
