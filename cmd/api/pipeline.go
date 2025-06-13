package main

import (
	"net/http"
	"strconv"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
	"github.com/go-chi/chi"
)

type CreatePipelinePayload struct {
	Name string `json:"name"`
}

func (app *application) createPipelineHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePipelinePayload
	if err := utils.ReadJson(r, &payload); err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err.Error())
	}

	// TO DO CREATE `USER` IMPLEMENTATION
	ctx := r.Context()
	pipeline := &store.Pipelines{
		UserID: 1,
		Name:   payload.Name,
	}

	if err := app.store.Pipelines.Create(ctx, pipeline); err != nil {
		//todo multiple error handling
		utils.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.WriteJson(w, http.StatusCreated, pipeline); err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPipelineHandler(w http.ResponseWriter, r *http.Request) {
	pipelineID := chi.URLParam(r, "pipelineID")
	pID, err := strconv.ParseInt(pipelineID, 10, 64)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()

	pipeline, err := app.store.Pipelines.GetByID(ctx, pID)
	if err != nil {
		//todo multiple error handling
		utils.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.WriteJson(w, http.StatusCreated, pipeline); err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) updatePipelineHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deletePipelineHandler(w http.ResponseWriter, r *http.Request) {

}
