package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LincolnG4/Haku/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PipelineService struct {
	*pgxpool.Pool
}

// NewPipelineService creates a new PipelineService with the given database connection.
func NewPipelineService() *PipelineService {
	// Connect to Postgres
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PWD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if conn.Ping(context.Background()) != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &PipelineService{conn}
}

func (p *PipelineService) InsertPipeline(ctx context.Context, pipeline *models.Pipeline) error {
	query := `INSERT INTO pipelines (
	id, user_id,
	name, description,
	schedule,
	created_at, updated_at,
	metadata, status
	) VALUES (@id, @user_id, @name, @description, @schedule, @created_at, @updated_at, @metadata, @status)`

	args := pgx.NamedArgs{
		"id":          uuid.New(),
		"user_id":     pipeline.UserID,
		"name":        pipeline.Name,
		"description": pipeline.Description,
		"schedule":    pipeline.Schedule,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
		"metadata":    pipeline.Metadata,
		"status":      "created",
	}

	_, err := p.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
