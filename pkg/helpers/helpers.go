package helpers

import (
	"errors"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/satnamSandhu2001/stackjet/pkg"
)

// Bool, Int, Int16, String, Uint are helpers to create a pointer to a values
func Bool(v bool) *bool       { return &v }
func Int(v int) *int          { return &v }
func Int16(v int16) *int16    { return &v }
func String(v string) *string { return &v }
func Uint(v uint) *uint       { return &v }

// GenerateStackDirPath generates a directory path for a stack based on its repo-url in default base dir
func GenerateStackDirPath(repoUrl string) string {
	// Remove .git suffix if present
	repoUrl = strings.TrimSuffix(repoUrl, ".git")

	// Extract repo name from URL
	parts := strings.Split(repoUrl, "/")
	repoName := parts[len(parts)-1]

	// Slugify
	slug := strings.TrimSpace(strings.ToLower(
		regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(repoName, "-"),
	))

	// Trim to 30 characters max
	if len(slug) > 30 {
		slug = slug[:30]
	}
	slug = strings.Trim(slug, "-") // clean up trailing dash

	uid := uuid.New().String()[:3]
	basePath := pkg.Config().DEFAULT_STACKS_BASE_DIR
	finalPath := filepath.Join(basePath, slug+"__"+uid)
	return finalPath
}

// MultiLogger is a logger that writes to multiple writers
type MultiLogger struct {
	writers []io.Writer
}

func NewMultiLogger(writers ...io.Writer) *MultiLogger {
	return &MultiLogger{writers}
}

func (m *MultiLogger) Write(p []byte) (int, error) {
	for _, w := range m.writers {
		w.Write(p) // ignore individual write errors
	}
	return len(p), nil
}

// accepts only: [npm | yarn | pnpm] start or [npm | yarn | pnpm] run <script>
func ValidateNodeStartCommand(command string) error {
	command = strings.TrimSpace(command)

	// disallow chaining, pipes
	if strings.ContainsAny(command, "&|;") {
		return errors.New("chaining or piping is not allowed in start command")
	}

	// allow only npm|yarn|pnpm start / run script-name
	validPattern := regexp.MustCompile(`^(npm|yarn|pnpm)\s+(start|run\s+[a-zA-Z0-9:_-]+)$`)
	if !validPattern.MatchString(command) {
		return errors.New("start command must be 'npm|yarn|pnpm start' or 'run <script>'")
	}

	return nil

}
