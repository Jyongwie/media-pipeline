package presentation

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/Jyongwie/media-pipeline/backend/internal/domain"
	"github.com/Jyongwie/media-pipeline/backend/internal/infrastructure"
)

// RenderHandler struct holds our dependencies (like the database repository)
type RenderHandler struct {
	repo *infrastructure.Repository
	hub *Hub
}

// NewRenderHandler is the constructor
func NewRenderHandler(repo *infrastructure.Repository, hub *Hub) *RenderHandler {
	return &RenderHandler{repo: repo, hub: hub}
}

// CreateJob handles the POST /api/jobs request
func (h *RenderHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	// 1. We define a temporary struct just for parsing the incoming JSON.
	// This keeps our web logic separate from our core domain logic.
	var req struct {
		AssetType        string `json:"asset_type"`
		Resolution       string `json:"resolution"`
		LightingProfile  string `json:"lighting_profile"`
		CameraEffect     string `json:"camera_effect"`
		AudioSensitivity string `json:"audio_sensitivity"`
	}

	// 2. Decode the incoming JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 3. Map the web request to our Core Domain Entity
	job := &domain.RenderJob{
		AssetType:        req.AssetType,
		Resolution:       req.Resolution,
		LightingProfile:  req.LightingProfile,
		CameraEffect:     req.CameraEffect,
		AudioSensitivity: req.AudioSensitivity,
	}

	// 4. In a full app, we'd pass this to an Application Use Case, 
	// but for this iteration, we'll save it directly via the repository.
	// (Note: You will need to add a SaveJob method to your postgres_repository.go!)
	
	err := h.repo.SaveJob(r.Context(), job)
	if err != nil {
		log.Printf("Database error while saving job: %v\n", err)
		http.Error(w, "Failed to queue job", http.StatusInternalServerError)
		return
	}
	h.hub.Broadcast()

	// 5. Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "job queued successfully"})
}

// GetJobs handles the GET /api/jobs request
func (h *RenderHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetAllJobs(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	// Because we didn't add JSON tags to our Domain struct, Go will 
	// automatically capitalize the keys in the JSON response (e.g., "AssetType").
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jobs)
}