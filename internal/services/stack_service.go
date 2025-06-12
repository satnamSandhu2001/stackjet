package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/models"
	"github.com/satnamSandhu2001/stackjet/pkg/helpers"
)

type StackService struct {
	db *sqlx.DB
}

func NewStackService(db *sqlx.DB) *StackService {
	return &StackService{
		db: db,
	}
}

func (s *StackService) CreateStack(ctx context.Context, data *dto.Stack_CreateRequest) (int64, error) {
	directory := helpers.GenerateStackDirPath(data.RepoUrl)
	uuid := uuid.New().String()
	if data.Name == "" {
		data.Name = strings.Split(directory, "/")[len(strings.Split(directory, "/"))-1]
	}
	columns := []string{"name", "uuid", "type", "directory", "port", "commands"}
	values := []any{data.Name, uuid, data.Type, directory, data.Port, data.Commands}

	if data.RepoUrl != "" {
		columns = append(columns, "repo_url")
		values = append(values, data.RepoUrl)
	}
	if data.Branch != "" {
		columns = append(columns, "branch")
		values = append(values, data.Branch)
	}
	if data.Remote != "" {
		columns = append(columns, "remote")
		values = append(values, data.Remote)
	}

	query_stack, args_stack, err := sq.Insert("stacks").Columns(columns...).Values(values...).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return 0, err
	}

	result_stack, err := s.db.ExecContext(ctx, query_stack, args_stack...)
	if err != nil {
		return 0, err
	}
	newStackID, err := result_stack.LastInsertId()
	if err != nil {
		return 0, err
	}

	return newStackID, nil
}

func (s *StackService) GetStackByID(ctx context.Context, id int64) (*models.Stack, error) {
	var stack models.Stack

	query, args, err := sq.Select("*").From("stacks").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return nil, err
	}
	if err := s.db.GetContext(ctx, &stack, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &stack, nil
}

func (s *StackService) GetStackByDirectory(ctx context.Context, directory string) (*models.Stack, error) {
	var stack models.Stack

	query, args, err := sq.Select("*").From("stacks").Where(sq.Eq{"directory": directory}).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return nil, err
	}

	if err := s.db.GetContext(ctx, &stack, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &stack, nil
}

func (s *StackService) UpdateStack(ctx context.Context, data *dto.Stack_UpdateRequest) error {
	// check if stack exists
	existingStack, err := s.GetStackByID(ctx, data.ID)
	if err != nil {
		return err
	}
	if existingStack.ID == 0 {
		return errors.New("stack not found")
	}

	builder := sq.Update("stacks").Where(sq.Eq{"id": existingStack.ID})
	if data.Name != "" {
		builder = builder.Set("name", data.Name)
	}
	if data.RepoUrl != "" {
		builder = builder.Set("repo_url", data.RepoUrl)
	}
	if data.Branch != "" {
		builder = builder.Set("branch", data.Branch)
	}
	if data.Remote != "" {
		builder = builder.Set("remote", data.Remote)
	}
	if data.CreatedSuccessfully != nil {
		builder = builder.Set("created_successfully", data.CreatedSuccessfully)
	}

	query, args, err := builder.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (s *StackService) CreateDeployment(ctx context.Context, data *dto.Deployment_CreateRequest) (int64, error) {
	cols := []string{"stack_id", "status"}
	values := []any{data.StackID, data.Status}

	if data.CommitHash != nil {
		cols = append(cols, "commit_hash")
		values = append(values, *data.CommitHash)
	}
	if data.RolledBackFromID != nil {
		cols = append(cols, "rolled_back_from_id")
		values = append(values, *data.RolledBackFromID)
	}

	query, args, err := sq.Insert("deployments").Columns(cols...).Values(values...).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return 0, err
	}
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	newDeploymentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newDeploymentID, nil
}

func (s *StackService) UpdateDeployment(ctx context.Context, data *dto.Deployment_UpdateRequest) (*models.Deployment, error) {
	if data.ID == 0 {
		return nil, errors.New("deployment id is required")
	}

	builder := sq.Update("deployments").Where(sq.Eq{"id": data.ID})
	if data.Status != "" {
		builder = builder.Set("status", data.Status)
	}
	if data.CommitHash != "" {
		builder = builder.Set("commit_hash", data.CommitHash)
	}
	if data.RolledBackFromID != nil {
		builder = builder.Set("rolled_back_from_id", *data.RolledBackFromID)
	}

	query, args, err := builder.PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return nil, err
	}

	var deployment *models.Deployment
	if err := s.db.GetContext(ctx, &deployment, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return deployment, nil
}

func (s *StackService) CreatePM2(ctx context.Context, data *dto.PM2_CreateRequest) (int64, error) {
	cols := []string{"stack_id", "name", "script"}
	values := []any{data.StackID, data.Name, data.Script}

	if data.Watch != nil {
		cols = append(cols, "watch")
		values = append(values, *data.Watch)
	}
	if data.Instances != nil {
		cols = append(cols, "insances")
		values = append(values, *data.Instances)
	}

	query, args, err := sq.Insert("pm2_configs").Columns(cols...).Values(values...).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return 0, err
	}
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (s *StackService) GetPM2byStackID(ctx context.Context, id int64) (*models.PM2, error) {
	var pm2 models.PM2

	query, args, err := sq.Select("*").From("pm2_configs").Where(sq.Eq{"stack_id": id}).PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return nil, err
	}
	if err := s.db.GetContext(ctx, &pm2, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &pm2, nil

}
