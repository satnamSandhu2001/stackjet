package commands

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/satnamSandhu2001/stackjet/pkg/logger"
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
	logger.EmitLog(args.Logger, fmt.Sprintf("> Executing: %s %s %s\n", strings.Join(envVarsList, " "), args.Name, strings.Join(args.Args, " ")))
	// }

	// Create command
	cmd := exec.Command(args.Name, args.Args...)
	cmd.Env = os.Environ()
	for k, v := range args.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("could not get stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("could not get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("could not start command: %w", err)
	}

	var output strings.Builder

	done := make(chan error, 2)
	go func() {
		done <- streamOutput(stdoutPipe, args.Logger, &output)
	}()
	go func() {
		done <- streamOutput(stderrPipe, args.Logger, &output)
	}()

	// wait for both stdout and stderr
	for range 2 {
		<-done
	}

	if err := cmd.Wait(); err != nil {
		return output.String(), fmt.Errorf("command failed: %w", err)
	}
	return output.String(), nil
}

// streams the output of a command to a writer
func streamOutput(reader io.Reader, w io.Writer, output *strings.Builder) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		logger.EmitLog(w, line)
		output.WriteString(line + "\n")

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	return scanner.Err()
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
