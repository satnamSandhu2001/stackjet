package nodejs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/satnamSandhu2001/stackjet/internal/core/pm2"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/logger"
)

func DeployStack(w io.Writer, ctx context.Context, service services.StackService, stack *models.Stack) error {
	// verify installation
	if err := verifyInstallation(w); err != nil {
		return err
	}

	// check for package.json file in the root folder of the project
	logger.EmitLog(w, "")
	logger.EmitLog(w, "⚓ Checking for package manager file...")

	pkgManager, err := detectPackageManager(stack.Directory)
	if err != nil {
		return err
	}
	if pkgManager == "" {
		return errors.New("no supported package manager (npm, yarn or pnpm) found in project root folder")
	}

	// execute build command
	if stack.Commands.Build != "" {
		logger.EmitLog(w, "")
		logger.EmitLog(w, "🛠️ Building application...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "bash", Args: []string{"-c", stack.Commands.Build}}); err != nil {
			return err
		}
	}

	//  handle pm2 + start app
	logger.EmitLog(w, "")
	logger.EmitLog(w, fmt.Sprintf("🚀 Starting %v application...\n", stack.Type))
	if err := pm2.StartProcess(w, ctx, service, stack); err != nil {
		return err
	}

	logger.EmitLog(w, "")
	logger.EmitLog(w, "🚀 Application started successfully")
	return nil
}

func verifyInstallation(w io.Writer) error {
	if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: w, Name: "node", Args: []string{"--version"}}); err != nil {
		return fmt.Errorf("nodejs is not installed: %w", err)
	}

	return nil
}

func detectPackageManager(projectDir string) (string, error) {
	tools := []struct {
		tool     string
		lockfile string
	}{
		{"npm", "package-lock.json"},
		{"yarn", "yarn.lock"},
		{"pnpm", "pnpm-lock.yaml"},
	}

	for _, t := range tools {
		lockPath := filepath.Join(projectDir, t.lockfile)
		if err := commands.FileExists(lockPath); err == nil {
			return t.tool, nil
		}
	}

	return "", errors.New("no supported package manager (npm, yarn or pnpm) found in project root folder")
}
