package main

import (
	"net/http"

	"github.com/LincolnG4/Haku/internal/utils"
)

func (a *application) createPipelineHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status": "ok",
		"env":    a.config.env,
	}

	if err := utils.WriteJson(w, http.StatusOK, data); err != nil {
		utils.WriteJsonError(w, http.StatusServiceUnavailable, err.Error())
	}

}
