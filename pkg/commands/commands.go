package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand runs a command with the given name and arguments
func RunCommand(name string, args ...string) (string, error) {
	// if cmd.Verbose {
	fmt.Printf("Executing: %s %s\n", name, strings.Join(args, " "))
	// }

	// Create command
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("⭕ Command Failed: %v, %v, %s\n", name, err, stdout.String())
		fmt.Println(stderr.String())
		return stdout.String(), fmt.Errorf("%v: %v\n%s", name, err, stderr.String())
	}
	fmt.Print(stdout.String())
	return stdout.String(), nil
}
