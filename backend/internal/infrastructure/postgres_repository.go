package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Jyongwie/media-pipeline/backend/internal/domain"
)

// Repository holds our database connection pool
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository initializes the connection to PostgreSQL
func NewRepository(ctx context.Context, connectionString string) (*Repository, error) {
	// We parse the connection string (e.g., postgres://admin:secretpassword@localhost:5432/mediadb)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %v", err)
	}

	// pgxpool automatically manages a pool of active connections, 
	// which is vital for highly concurrent Go applications.
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &Repository{db: pool}, nil
}

// Close cleanly shuts down the database connection pool.
// It should be called when the application is shutting down.
func (r *Repository) Close() {
	if r.db != nil {
		r.db.Close()
		fmt.Println("Database connection pool closed.")
	}
}

// GetNextQueuedJob safely locks and retrieves the oldest queued job.
// It instantly changes the status to 'processing' to prevent other workers from grabbing it.
func (r *Repository) GetNextQueuedJob(ctx context.Context) (*domain.RenderJob, error) {
	// The Staff Engineer SQL query: FOR UPDATE SKIP LOCKED
	query := `
		UPDATE render_jobs
		SET status = 'processing'
		WHERE id = (
			SELECT id FROM render_jobs 
			WHERE status = 'queued' 
			ORDER BY created_at ASC 
			FOR UPDATE SKIP LOCKED 
			LIMIT 1
		)
		RETURNING id, asset_type, resolution, lighting_profile, camera_effect, audio_sensitivity, status, created_at;
	`

	var j domain.RenderJob
	err := r.db.QueryRow(ctx, query).Scan(
		&j.ID, &j.AssetType, &j.Resolution, &j.LightingProfile,
		&j.CameraEffect, &j.AudioSensitivity, &j.Status, &j.CreatedAt,
	)

	// If no rows are found, it just means the queue is currently empty
	if err == pgx.ErrNoRows {
		return nil, nil 
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next job: %v", err)
	}

	return &j, nil
}

// MarkJobCompleted updates the job status to 'completed'
func (r *Repository) MarkJobCompleted(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "UPDATE render_jobs SET status = 'completed' WHERE id = $1", id)
	return err
}

// SaveJob inserts a newly created RenderJob into the database.
func (r *Repository) SaveJob(ctx context.Context, job *domain.RenderJob) error {
	// 1. Define the SQL insert statement using parameterized queries ($1, $2...) 
	// to prevent SQL injection.
	query := `
		INSERT INTO render_jobs (
			asset_type, 
			resolution, 
			lighting_profile, 
			camera_effect, 
			audio_sensitivity
		) VALUES ($1, $2, $3, $4, $5)
	`

	// 2. Execute the query using the connection pool
	_, err := r.db.Exec(ctx, query,
		job.AssetType,
		job.Resolution,
		job.LightingProfile,
		job.CameraEffect,
		job.AudioSensitivity,
	)

	// 3. Handle any database errors
	if err != nil {
		return fmt.Errorf("failed to save job to database: %w", err)
	}

	return nil
}

// GetAllJobs retrieves all rendering jobs from the database, sorted by newest.
func (r *Repository) GetAllJobs(ctx context.Context) ([]domain.RenderJob, error) {
	query := `
		SELECT id, asset_type, resolution, lighting_profile, camera_effect, audio_sensitivity, status, created_at 
		FROM render_jobs 
		ORDER BY created_at DESC
	`
	
	// Query returns multiple rows
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close() // Always defer closing rows to prevent memory leaks!

	var jobs []domain.RenderJob

	// Loop through the results and map them to our Domain struct
	for rows.Next() {
		var j domain.RenderJob
		err := rows.Scan(
			&j.ID, &j.AssetType, &j.Resolution, &j.LightingProfile, 
			&j.CameraEffect, &j.AudioSensitivity, &j.Status, &j.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}