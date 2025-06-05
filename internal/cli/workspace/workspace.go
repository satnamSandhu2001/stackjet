package workspace

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/satnamSandhu2001/stackjet/pkg/commands"
)

// Returns base folder's (originalName, formattedName, error)
func BaseFolderName() (string, string, error) {
	path, err := commands.RunCommand("pwd")
	if err != nil {
		return "", "", err
	}
	trimmed := strings.TrimRight(path, "/")
	originalName := strings.TrimSpace(trimmed[strings.LastIndex(trimmed, "/")+1:])
	formattedName := strings.ToUpper(regexp.MustCompile(`[-_]`).ReplaceAllString(originalName, " "))
	return originalName, formattedName, nil
}

// Enter workspace of the project
func EnterWorkspace(path string) error {
	_, folderNameFormatted, err := BaseFolderName()
	if err != nil {
		return err
	}

	fmt.Printf("\n~~~~~~ REDEPLOYING (%s) ~~~~~~\n", folderNameFormatted)

	checkDir, _ := commands.RunCommand("pwd")
	if err := os.Chdir(path); err != nil {
		return err
	}
	fmt.Printf("📁 Working dir: %v \n", checkDir)

	return nil
}
