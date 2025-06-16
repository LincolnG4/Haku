package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
)

type organizationKey string

const organizationCtx organizationKey = "organization"

type CreateOrganizationPayload struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

func (app *application) createOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateOrganizationPayload
	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	organization := &store.Organization{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := app.store.Organization.Create(ctx, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	organization := app.getOrganizationFromContext(r)
	if err := utils.JsonResponse(w, http.StatusOK, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getOrganizationFromContext(r *http.Request) *store.Organization {
	organization, ok := r.Context().Value(organizationCtx).(*store.Organization)
	if !ok {
		panic("pipeline not found in context")
	}
	return organization
}

func (app *application) organizationContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgID, err := utils.GetURLParamInt64(r, "organizationID")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		organization, err := app.store.Organization.GetByID(ctx, orgID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, organizationCtx, organization)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
