package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
)

type pipelineKey string

const pipelineCtx pipelineKey = "pipeline"

type CreatePipelinePayload struct {
	Name string `json:"name" validate:"required,max=255"`
}

func (app *application) createPipelineHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePipelinePayload
	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// TO DO CREATE `USER` IMPLEMENTATION
	ctx := r.Context()
	pipeline := &store.Pipelines{
		UserID: 2,
		Name:   payload.Name,
	}

	if err := app.store.Pipelines.Create(ctx, pipeline); err != nil {
		//todo multiple error handling
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, pipeline); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPipelineHandler(w http.ResponseWriter, r *http.Request) {
	pipeline := getPipelineFromContext(r)

	if err := utils.JsonResponse(w, http.StatusOK, pipeline); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePipelineHandler(w http.ResponseWriter, r *http.Request) {
	pipelineID, err := utils.GetURLParamInt64(r, "pipelineID")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Pipelines.Delete(ctx, pipelineID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePipelinePayload struct {
	Name string `json:"name" validate:"required,max=255"`
}

func (app *application) updatePipelineHandler(w http.ResponseWriter, r *http.Request) {
	pipeline := getPipelineFromContext(r)

	var payload UpdatePipelinePayload
	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	pipeline.Name = payload.Name

	ctx := r.Context()
	if err := app.store.Pipelines.Update(ctx, pipeline); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusOK, pipeline); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) pipelineContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pipelineID, err := utils.GetURLParamInt64(r, "pipelineID")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		pipeline, err := app.store.Pipelines.GetByID(ctx, pipelineID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, pipelineCtx, pipeline)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPipelineFromContext(r *http.Request) *store.Pipelines {
	pipeline, ok := r.Context().Value(pipelineCtx).(*store.Pipelines)
	if !ok {
		panic("pipeline not found in context")
	}

	return pipeline
}
