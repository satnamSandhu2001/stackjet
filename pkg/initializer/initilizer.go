package initializer

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/pkg"
)

// InitApp initializes the app by creating the config directory and generating a JWT token and lock file
func InitializeApp(forceRecreate bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		log.Printf("Error: Unable to determine user home directory: %v", err)
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to access user home directory.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}

	stackjetDirPath := filepath.Join(homeDir, ".stackjet")
	lockFilePath := filepath.Join(stackjetDirPath, "init.lock")

	if _, err := os.Stat(lockFilePath); err == nil {
		if !forceRecreate {
			return
		}
		fmt.Println("Recreating config forcefully...")
		if err := os.RemoveAll(stackjetDirPath); err != nil {
			log.Printf("Error while deleting stackjet directory: %v", err)
			fmt.Println("❌ StackJet Configuration Error")
			fmt.Println("Unable to delete the existing configuration directory.")
			fmt.Println("\nPlease check your system permissions and try again.")
			os.Exit(1)
		}
	}

	_, err = createStackJetDir()
	if err != nil {
		log.Printf("Error creating stackjet directory: %v", err)
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to create the configuration directory.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}
	configPath := filepath.Join(stackjetDirPath, "config.json")
	createDefaultConfig(configPath)
	generateAndStoreJWT(stackjetDirPath)

	createLockFile(lockFilePath) // lock file must exist before below logic (it uses config package)
	if err := createSitesDirectory(); err != nil {
		log.Printf("Error creating default base directory: %v", err)
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Printf("Unable to create the default base directory: %s", pkg.Config().DEFAULT_STACKS_BASE_DIR)
		fmt.Println("\nPlease check your system permissions and or manually create this folder")
		os.Exit(1)
	}
	dbConn := database.Connect()
	defer dbConn.Close()
	database.RunInitSQL()
}

func createLockFile(lockFilePath string) {
	file, err := os.Create(lockFilePath)
	if err != nil {
		log.Printf("Error creating lock file: %v", err)
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to create the lock file.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)

	}
	defer file.Close()
}

func createStackJetDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}

	stackjetDir := filepath.Join(home, ".stackjet")

	// Check if directory already exists before creating
	if info, err := os.Stat(stackjetDir); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("config path exists but is not a directory: %s", stackjetDir)
		}
		return stackjetDir, nil // Directory already exists
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("could not check config directory: %w", err)
	}

	// Only create if it doesn't exist
	if err := os.MkdirAll(stackjetDir, 0755); err != nil {
		return "", fmt.Errorf("could not create config dir: %w", err)
	}

	return stackjetDir, nil
}

func createDefaultConfig(configPath string) {
	defaultConfig := pkg.AppConfig{
		PORT:                    8080,
		GIT_BRANCH:              "master",
		GIT_REMOTE:              "origin",
		GIT_RESET:               true,
		DEFAULT_STACKS_BASE_DIR: "/var/www/sites",
	}
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to create the configuration file.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to open the configuration file.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to write to the configuration file.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}
}

// generateAndStoreJWT creates a cryptographically secure random token
// and stores it in the specified directory with restricted file permissions.
// The token file is created at <stackjetDirPath>/jwt.token with 0600 permissions.
func generateAndStoreJWT(stackjetDirPath string) {
	token, err := generateRandomHex(32)
	if err != nil {
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to generate secure token.")
		os.Exit(1)
	}
	tokenPath := filepath.Join(stackjetDirPath, "jwt.token")

	// token with permission 0600
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		fmt.Println("❌ StackJet Configuration Error")
		fmt.Println("Unable to write to the jwt token file.")
		fmt.Println("\nPlease check your system permissions and try again.")
		os.Exit(1)
	}

}

func generateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("could not generate secure token: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func createSitesDirectory() error {
	cmd := exec.Command("sudo", "mkdir", "-p", pkg.Config().DEFAULT_STACKS_BASE_DIR)
	if err := cmd.Run(); err != nil {
		return err
	}
	// Change ownership to current user and nginx group
	user := os.Getenv("USER")
	chownCmd := exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s:www-data", user), pkg.Config().DEFAULT_STACKS_BASE_DIR)
	if err := chownCmd.Run(); err != nil {
		return err
	}

	return nil
}
