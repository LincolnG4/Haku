package main

import (
	"net/http"

	"github.com/LincolnG4/Haku/internal/utils"
)

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// IMPLEMENT HEALTH CHECKS
	data := map[string]string{
		"status": "ok",
		"env":    a.config.env,
	}

	if err := utils.WriteJson(w, http.StatusOK, data); err != nil {
		utils.WriteJsonError(w, http.StatusServiceUnavailable, err.Error())
	}

}
