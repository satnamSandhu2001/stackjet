package commands

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// RunCommand runs a command with the given name and arguments
func RunCommand(name string, args ...string) (string, error) {
	// if cmd.Verbose {
	fmt.Printf("executing: %s %s\n", name, strings.Join(args, " "))
	// }

	// Create command
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("⭕ command Failed: %v, %v, %s\n", name, err, stdout.String())
		fmt.Println(stderr.String())
		return stdout.String(), fmt.Errorf("%v: %v\n%s", name, err, stderr.String())
	}
	fmt.Print(stdout.String())
	return stdout.String(), nil
}

// FileExists checks if a file exists at the given path on system
func FileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("⭕ system cannot find the file: %s", path)
		return fmt.Errorf("system cannot find the file: %s", path)
	}
	return nil
}

// GetFreePort returns an available port number
func GetFreePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to find free port: %w", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// IsPortFree checks if a port is available to use
func IsPortFree(port uint16) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false // Port is in use or blocked
	}
	_ = ln.Close()
	return true
}
