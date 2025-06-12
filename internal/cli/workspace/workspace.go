package workspace

import (
	"fmt"
	"io"
	"os"

	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

// Enter workspace of the project
func EnterWorkspace(logger io.Writer, stack *models.Stack) error {
	fmt.Fprintln(logger, "📁 Entering workspace...")
	if err := os.Chdir(stack.Directory); err != nil {
		return err
	}
	checkDir, _ := commands.RunCommand(logger, "pwd")
	fmt.Fprintf(logger, "📁 Working dir: %v \n", checkDir)

	return nil
}
