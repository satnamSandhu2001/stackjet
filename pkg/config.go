package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type configutation struct {
	PORT          uint
	JWT_TOKEN     string
	DB_URL        string
	GIT_BRANCH    string
	GIT_REMOTE    string
	GIT_RESET     bool
	VALID_STACKS  []string
	APPS_BASE_DIR string
}

var (
	config *configutation
	once   sync.Once
)

// Config returns the singleton configuration instance for whole app
func Config() *configutation {
	once.Do(func() {
		// create stackjet config dir
		home, err := os.UserHomeDir()
		if err != nil {
			panic("could not determine user home directory")
		}
		stackjetDir := filepath.Join(home, ".stackjet")
		if err := os.MkdirAll(stackjetDir, 0755); err != nil {
			panic(fmt.Sprintf("could not create stackjet config dir: %v", err))
		}
		dbPath := filepath.Join(stackjetDir, "stackjet.db")

		config = &configutation{
			PORT:          uint(GetEnvInt("PORT", 8080)),
			JWT_TOKEN:     GetEnv("JWT_SECRET", "some_secret_token"),
			DB_URL:        fmt.Sprintf("file:%s?_fk=1", dbPath),
			GIT_BRANCH:    "master",
			GIT_REMOTE:    "origin",
			GIT_RESET:     false,
			VALID_STACKS:  []string{"nodejs"},
			APPS_BASE_DIR: "/var/www/html",
		}
	})

	return config
}
