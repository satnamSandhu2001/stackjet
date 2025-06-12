package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type AppConfig struct {
	PORT                    uint     `json:"port"`
	JWT_TOKEN               string   `json:"-"`
	GIT_BRANCH              string   `json:"git_branch"`
	GIT_REMOTE              string   `json:"git_remote"`
	GIT_RESET               bool     `json:"git_reset"`
	DEFAULT_STACKS_BASE_DIR string   `json:"default_stacks_base_dir"`
	DB_URL                  string   `json:"-"`
	VALID_STACKS            []string `json:"-"`
}

var (
	config *AppConfig
	once   sync.Once
)

// Config initializes and returns the app configuration
func Config() *AppConfig {
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil || homeDir == "" {
			fmt.Println("❌ StackJet Configuration Error")
			fmt.Println("Unable to access user home directory.")
			fmt.Println("\nPlease check your system permissions and try again.")
			os.Exit(1)
		}

		stackjetDir := filepath.Join(homeDir, ".stackjet")
		lockFilePath := filepath.Join(stackjetDir, "init.lock")
		if _, err := os.Stat(lockFilePath); err != nil {
			fmt.Println("❌ StackJet Not Initialized")
			fmt.Println("\nStackJet has not been initialized on this system.")
			fmt.Println("\n🚀 To get started, run:")
			fmt.Println("   \033[1;34mstackjet init\033[0m")
			fmt.Println("\nThis will set up the necessary configuration files.")
			os.Exit(1)
		}
		configPath := filepath.Join(stackjetDir, "config.json")
		// Load config.json
		data, err := os.ReadFile(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("❌ StackJet Configuration Missing")
				fmt.Println("\nConfiguration file not found.")
				fmt.Println("\n🔧 To fix this issue, run:")
				fmt.Println("   \033[1;34mstackjet init\033[0m")
				fmt.Println("\nThis will recreate the necessary configuration.")
				os.Exit(1)
			}
			fmt.Println("❌ StackJet Configuration Error")
			fmt.Println("Unable to read configuration file.")
			fmt.Println("\n🔧 To fix this issue, try:")
			fmt.Println("   \033[1;34mstackjet init\033[0m")
			os.Exit(1)
		}

		var loaded AppConfig
		if err := json.Unmarshal(data, &loaded); err != nil {
			fmt.Println("❌ StackJet Configuration Error")
			fmt.Println("Configuration file is corrupted or invalid.")
			fmt.Println("\n🔧 To fix this issue, run:")
			fmt.Println("   \033[1;34mstackjet init\033[0m")
			fmt.Println("\nThis will recreate a fresh configuration.")
			os.Exit(1)
		}

		// Load token
		tokenData, err := os.ReadFile(filepath.Join(stackjetDir, "jwt.token"))
		if err != nil {
			fmt.Println("❌ StackJet Authentication Error")
			fmt.Println("Authentication token not found or corrupted.")
			fmt.Println("\n🔧 To fix this issue, run:")
			fmt.Println("   \033[1;34mstackjet init\033[0m")
			fmt.Println("\nThis will regenerate the authentication token.")
			os.Exit(1)
		}

		loaded.JWT_TOKEN = string(tokenData)
		loaded.VALID_STACKS = []string{"nodejs"}
		loaded.DB_URL = fmt.Sprintf("file:%s?_fk=1", filepath.Join(stackjetDir, "stackjet.db"))

		config = &loaded
	})

	return config
}
