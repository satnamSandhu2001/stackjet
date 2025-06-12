package nodejs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/satnamSandhu2001/stackjet/internal/cli/pm2"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

func DeployStack(logger io.Writer, ctx context.Context, service services.StackService, stack *models.Stack) error {
	// verify installation
	if err := verifyInstallation(logger); err != nil {
		return err
	}

	// check for package.json file in the root folder of the project
	fmt.Fprintln(logger, "\n⚓ Checking for package manager file...")

	pkgManager, err := detectPackageManager(stack.Directory)
	if err != nil {
		return err
	}
	if pkgManager == "" {
		return errors.New("no supported package manager (npm, yarn or pnpm) found in project root folder")
	}

	// install and build commands
	if stack.Commands.Build != "" {
		fmt.Fprintln(logger, "\n🛠️ Building application...")
		if _, err := commands.RunCommand(logger, "bash", "-c", stack.Commands.Build); err != nil {
			return err
		}
	}

	//  start pm2
	fmt.Fprintf(logger, "\n🚀 Starting %v application...\n", stack.Type)
	if err := pm2.StartProcess(logger, ctx, service, stack); err != nil {
		return err
	}

	fmt.Fprintln(logger, "\n🚀 Application started successfully")
	return nil
}

func verifyInstallation(logger io.Writer) error {
	if _, err := commands.RunCommand(logger, "node", "--version"); err != nil {
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
