package dto

import "github.com/satnamSandhu2001/stackjet/internal/models"

type Stack_CreateRequest struct {
	ID       int64                `json:"id" db:"id"`
	Name     string               `json:"name" db:"name" binding:"required"`
	Type     string               `json:"type" db:"type" binding:"required"`
	RepoUrl  string               `json:"repo_url" db:"repo_url" binding:"required"`
	Branch   string               `json:"branch" db:"branch"`
	Remote   string               `json:"remote" db:"remote"`
	Port     int                  `db:"port" json:"port" binding:"required"`
	Commands models.StackCommands `db:"commands" json:"commands" binding:"required"`
}
type Stack_DeployRequest struct {
	ID        int64  `json:"id" db:"id" binding:"required"`
	Branch    string `json:"branch" db:"branch"`
	Remote    string `json:"remote" db:"remote"`
	Directory string `db:"directory"` // only used for cli created stacks
	GitHash   string `json:"git_hash"`
	GitReset  bool
}

type Stack_UpdateRequest struct {
	ID                       int64  `json:"id" db:"id" binding:"required"`
	Name                     string `json:"name" db:"name"`
	RepoUrl                  string `json:"repo_url" db:"repo_url"`
	Branch                   string `json:"branch" db:"branch"`
	Remote                   string `json:"remote" db:"remote"`
	CreatedSuccessfully      *bool  `db:"created_successfully"`
	InitialDeploymentSuccess bool   `db:"initial_deployment_success"`
}

type Deployment_CreateRequest struct {
	StackID          int64   `db:"stack_id" json:"stack_id"`
	Status           string  `db:"status" json:"status"`
	CommitHash       *string `db:"commit_hash" json:"commit_hash"`
	RolledBackFromID *int64  `db:"rolled_back_from_id" json:"rolled_back_from_id"`
}
type Deployment_UpdateRequest struct {
	ID               int64  `db:"id" json:"id"`
	Status           string `db:"status" json:"status"`
	CommitHash       string `db:"commit_hash" json:"commit_hash"`
	RolledBackFromID *int64 `db:"rolled_back_from_id" json:"rolled_back_from_id"`
}

type PM2_CreateRequest struct {
	StackID   int64  `json:"stack_id" db:"stack_id"`
	Script    string `json:"script" db:"script"`
	Name      string `json:"name" db:"name"`
	Watch     *bool  `json:"watch" db:"watch"`
	Instances *int   `json:"instances" db:"instances"`
}
