package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Stack struct {
	ID                       int64         `db:"id" json:"id"`
	Name                     string        `db:"name" json:"name"`
	Uuid                     string        `db:"uuid" json:"uuid"`
	Directory                string        `db:"directory" json:"directory"`
	Type                     string        `db:"type" json:"type"`
	RepoUrl                  string        `db:"repo_url" json:"repo_url"`
	Branch                   string        `db:"branch" json:"branch"`
	Remote                   string        `db:"remote" json:"remote"`
	Port                     int           `db:"port" json:"port"`
	Commands                 StackCommands `db:"commands" json:"commands"`
	CreatedSuccessfully      bool          `db:"created_successfully" json:"created_successfully"`
	InitialDeploymentSuccess bool          `db:"initial_deployment_success" json:"initial_deployment_success"`
	CreatedAt                string        `db:"created_at" json:"created_at"`
}

type StackCommands struct {
	Build string `json:"build"`
	Start string `json:"start"`
	Post  string `json:"post"`
}

// For saving to DB
func (s StackCommands) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// For reading from DB
func (s *StackCommands) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, s)
}

const (
	DEPLOYMENT_STATUS_IN_PROGRESS = "in_progress"
	DEPLOYMENT_STATUS_SUCCESS     = "success"
	DEPLOYMENT_STATUS_FAILED      = "failed"
)

type Deployment struct {
	ID               int64  `db:"id" json:"id"`
	StackID          int64  `db:"stack_id" json:"stack_id"`
	Status           string `db:"status" json:"status"`
	CommitHash       string `db:"commit_hash" json:"commit_hash"`
	RolledBackFromID int64  `db:"rolled_back_from_id" json:"rolled_back_from_id"`
	DeployedAt       string `db:"deployed_at" json:"deployed_at"`
}

type DeploymentLog struct {
	ID           int64  `db:"id" json:"id"`
	DeploymentID int64  `db:"deployment_id" json:"deployment_id"`
	Log          string `db:"log" json:"log"`
}

type PM2 struct {
	ID        int64  `json:"id" db:"id"`
	StackID   int64  `json:"stack_id" db:"stack_id"`
	Script    string `json:"script" db:"script"`
	Name      string `json:"name" db:"name"`
	Watch     bool   `json:"watch" db:"watch"`
	Instances int    `json:"instances" db:"instances"`
}
