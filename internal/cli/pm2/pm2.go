package pm2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

func StartProcess(logger io.Writer, ctx context.Context, service services.StackService, stack *models.Stack) error {
	// verify installation
	if err := verifyInstallation(logger); err != nil {
		return err
	}

	pm2Data, err := service.GetPM2byStackID(ctx, stack.ID)
	if err != nil {
		return err
	}

	if pm2Data == nil { // create new record
		fmt.Fprintln(logger, "\n🚀 Creating pm2 process...")
		commandParts := strings.Fields(stack.Commands.Start)
		if _, err := service.CreatePM2(ctx, &dto.PM2_CreateRequest{
			StackID: stack.ID,
			Script:  commandParts[0] + " -- " + strings.Join(commandParts[1:], " "),
			Name:    stack.Name,
		}); err != nil {
			return err
		}
		pm2Data, err = service.GetPM2byStackID(ctx, stack.ID) // fetch record after creation
		if err != nil {
			return err
		}
	}

	// start pm2
	if !stack.InitialDeploymentSuccess {
		fmt.Fprintln(logger, "\n🚀 Starting pm2 process...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "pm2", Args: []string{"start", "--name", pm2Data.Name, pm2Data.Script}}); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(logger, "\n🚀 Restarting pm2 process...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "pm2", Args: []string{"restart", pm2Data.Name}}); err != nil {
			return err
		}
	}
	// post script
	if stack.Commands.Post != "" {
		fmt.Fprintln(logger, "\n🛠️ Running post commands...")
		if _, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "bash", Args: []string{"-c", stack.Commands.Post}}); err != nil {
			return err
		}
	}
	// verify status
	if err := validatePM2Process(pm2Data.Name); err != nil {
		return fmt.Errorf("pm2 process did not start properly: %w", err)
	}
	// save pm2 app list if first deployment
	if !stack.InitialDeploymentSuccess {
		commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "pm2", Args: []string{"save"}})
	}

	fmt.Fprintln(logger, "🚀 pm2 process started successfully")
	return nil
}
func verifyInstallation(logger io.Writer) error {
	version, err := commands.RunCommand(commands.RunCommandArgs{Logger: logger, Name: "pm2", Args: []string{"--version"}})
	if err != nil || version == "" {
		fmt.Fprintln(logger, "Please install pm2 (https://github.com/Unitech/pm2?tab=readme-ov-file#installing-pm2)")
		return fmt.Errorf("pm2 is not installed: %w", err)
	}

	return nil
}

// JSON-based validation function using <pm2 jlist>
func validatePM2Process(name string) error {
	out, err := exec.Command("pm2", "jlist").Output()
	if err != nil {
		return fmt.Errorf("failed to get pm2 jlist: %w", err)
	}

	var apps []map[string]interface{}
	if err := json.Unmarshal(out, &apps); err != nil {
		return fmt.Errorf("failed to parse pm2 jlist: %w", err)
	}

	for _, app := range apps {
		if appName, ok := app["name"].(string); ok && appName == name {
			// Check status
			monit, ok := app["pm2_env"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("pm2_env not found in pm2 jlist")
			}
			if status, ok := monit["status"].(string); ok && status == "online" {
				return nil
			}
			return fmt.Errorf("pm2 process status: %v", monit["status"])
		}
	}

	return fmt.Errorf("pm2 process %s not found in jlist", name)
}
