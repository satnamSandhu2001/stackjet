package nodejs

import (
	"fmt"

	"github.com/satnamSandhu2001/stackjet/internal/cli/pm2"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

func DeployStack() error {

	// check for package.json file in the root folder of the project
	fmt.Println("\n⚓ Checking for package.json file...")
	if err := commands.FileExists("./package.json"); err != nil {
		return err
	}

	fmt.Println("\n🧩 Installing dependencies...")
	if _, err := commands.RunCommand("npm", "ci"); err != nil {
		return err
	}

	// TODO: give process name as foldername
	if err := pm2.StartProcess("dummy-repo"); err != nil {
		return err
	}

	return nil
}
