package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LincolnG4/Haku/internal/models"
	"github.com/LincolnG4/Haku/internal/services"
)

type PipelineController struct {
	pipelineService *services.PipelineService
}

func NewPipelineController(pipelineService *services.PipelineService) *PipelineController {
	return &PipelineController{pipelineService: pipelineService}
}

func (p *PipelineController) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	newPipeline := models.Pipeline{}
	if err := json.NewDecoder(r.Body).Decode(&newPipeline); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validations
	if newPipeline.ID == "" {
		http.Error(w, "ID can't be empty", http.StatusBadRequest)
		return
	}

	if newPipeline.Name == "" {
		http.Error(w, "Name can't be empty", http.StatusBadRequest)
		return
	}

	// Validate DAG pipeline

	// Insert into database
	err := p.pipelineService.InsertPipeline(r.Context(), &newPipeline)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting pipeline: %v", err), http.StatusInternalServerError)
		return
	}

}
