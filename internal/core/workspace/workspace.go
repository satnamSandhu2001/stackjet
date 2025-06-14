package workspace

import (
	"fmt"
	"io"
	"os"

	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
	"github.com/satnamSandhu2001/stackjet/pkg/logger"
)

// Enter workspace of the project
func EnterWorkspace(w io.Writer, stack *models.Stack) error {
	logger.EmitLog(w, "📁 Entering workspace...")
	if err := os.Chdir(stack.Directory); err != nil {
		return err
	}
	checkDir, err := commands.RunCommand(commands.RunCommandArgs{
		Logger: w,
		Name:   "pwd",
	})
	if err != nil {
		return err
	}
	logger.EmitLog(w, fmt.Sprintf("📁 Working dir: %v \n", checkDir))
	return nil
}
