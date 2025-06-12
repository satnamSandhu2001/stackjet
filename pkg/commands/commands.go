package commands

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type RunCommandArgs struct {
	Logger io.Writer
	Name   string
	Args   []string
	Env    map[string]string
}

// RunCommand runs a command with the given name and arguments
func RunCommand(args RunCommandArgs) (string, error) {
	// if cmd.Verbose {
	envVarsList := make([]string, 0, len(args.Env))
	for k, v := range args.Env {
		envVarsList = append(envVarsList, fmt.Sprintf("%s=%v", k, v))
	}
	fmt.Fprintf(args.Logger, "Executing: %s %s %s\n", strings.Join(envVarsList, " "), args.Name, strings.Join(args.Args, " "))
	// }

	// Create command
	cmd := exec.Command(args.Name, args.Args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()

	// Set environment variables
	for key, value := range args.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// run command
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(args.Logger, "⭕ command Failed: %v, %v, %s\n", args.Name, err, stderr.String())
		return stdout.String(), fmt.Errorf("%v: %v\n%s", args.Name, err, stderr.String())
	}
	fmt.Fprintln(args.Logger, stdout.String())
	return stdout.String(), nil
}

// FileExists checks if a file exists at the given path on system
func FileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("system cannot find the file %s", path)
	}
	return nil
}

// checks if a stack folder exists in default base directory
func StackDirExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", path)
	}
	return nil
}

func CreateDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("directory path is empty")
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to directory: %w", err)
	}
	return nil
}

// validate and checks if port is available
func ValidatePort(port int) error {
	if port < 1024 || port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535")
	}

	// Check if port is already in use
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("port %d is already in use", port)
	}
	defer listener.Close()

	return nil
}
