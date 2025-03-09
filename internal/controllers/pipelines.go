package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LincolnG4/Haku/internal/models"
)

func CreatePipeline(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	newPipeline := models.Pipeline{}
	if err := json.NewDecoder(r.Body).Decode(&newPipeline); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(newPipeline)
}
