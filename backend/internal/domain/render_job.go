// backend/internal/domain/render_job.go
package domain

import (
	"errors"
	"time"
)

// RenderJob is our core domain entity. 
type RenderJob struct {
	ID               string
	AssetType        string    
	Resolution       string    
	LightingProfile  string    
	CameraEffect     string    
	AudioSensitivity string    
	Status           string    
	CreatedAt        time.Time
}

// MarkAsFailed enforces our business rules for failing a job
func (r *RenderJob) MarkAsFailed() error {
	if r.Status == "completed" {
		return errors.New("domain error: cannot fail a job that has already completed successfully")
	}
	r.Status = "failed"
	return nil
}

// StartProcessing enforces our business rules for starting a job
func (r *RenderJob) StartProcessing() error {
	if r.Status != "queued" {
		return errors.New("domain error: only queued jobs can begin processing")
	}
	r.Status = "processing"
	return nil
}