package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
)

type StackService struct {
	db *sqlx.DB
}

func NewStackService(db *sqlx.DB) *StackService {
	return &StackService{
		db: db,
	}
}

func (s *StackService) CreateStack(ctx context.Context, data *dto.Stack_CreateRequest) error {
	slug := strings.TrimSpace(strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(data.Name, "-"))) + fmt.Sprintf("_u%x", rand.Intn(1000))
	uniqueId := uuid.New()
	columns := []string{"uuid", "name", "slug", "type", "port"}
	values := []any{uniqueId, data.Name, slug, data.Type, data.Port}

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
	builder := sq.Insert("stacks").Columns(columns...).Values(values...).PlaceholderFormat(sq.Question)
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	log.Println("query", query)
	log.Println("args", args)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	data.ID = id

	// TODO : trigger stack create app and update stack deployment status
	return nil
}
