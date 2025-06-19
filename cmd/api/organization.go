package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
	"github.com/go-chi/chi"
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

	// Create organization
	if err := app.store.Organizations.Create(ctx, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Get user from context
	user := getUserFromContext(r)
	member := store.OrganizationMember{
		UserID:         user.ID,
		OrganizationID: organization.ID,
		RoleID:         store.AdminRole,
	}

	// Add user
	if err := app.store.Organizations.AddMember(ctx, &member); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type AddMemberPayload struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

func (app *application) addMemberHandler(w http.ResponseWriter, r *http.Request) {
	organization := app.getOrganizationFromContext(r)

	var payload AddMemberPayload
	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	member := &store.OrganizationMember{
		UserID:         payload.UserID,
		OrganizationID: organization.ID,
		RoleID:         payload.RoleID,
	}

	if err := app.store.Organizations.AddMember(ctx, member); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, member); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getMembersHandler(w http.ResponseWriter, r *http.Request) {
	organization := app.getOrganizationFromContext(r)

	ctx := r.Context()
	members, err := app.store.Organizations.GetMembers(ctx, organization.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JsonResponse(w, http.StatusOK, members); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	// Temporarily handle organizationID extraction here since middleware is disabled
	orgIDStr := chi.URLParam(r, "organizationID")
	if orgIDStr == "" {
		app.logger.Error("organizationID parameter is empty")
		app.internalServerError(w, r, fmt.Errorf("organizationID parameter is empty"))
		return
	}

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		app.logger.Error("failed to parse organizationID", "error", err, "orgIDStr", orgIDStr)
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	organization, err := app.store.Organizations.GetByID(ctx, orgID)
	if err != nil {
		app.logger.Error("failed to get organization from database", "error", err, "orgID", orgID)
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JsonResponse(w, http.StatusOK, organization); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getOrganizationFromContext(r *http.Request) *store.Organization {
	organization, ok := r.Context().Value(organizationCtx).(*store.Organization)
	if !ok {
		panic("organization not found in context")
	}
	return organization
}

func (app *application) organizationContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add debugging
		app.logger.Info("organizationContextMiddleware called", "url", r.URL.Path)

		// Try direct chi parameter extraction first
		orgIDStr := chi.URLParam(r, "organizationID")
		app.logger.Info("raw organizationID from chi", "orgIDStr", orgIDStr, "isEmpty", orgIDStr == "")

		if orgIDStr == "" {
			app.logger.Error("organizationID parameter is empty")
			app.internalServerError(w, r, fmt.Errorf("organizationID parameter is empty"))
			return
		}

		orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
		if err != nil {
			app.logger.Error("failed to parse organizationID", "error", err, "orgIDStr", orgIDStr)
			app.internalServerError(w, r, err)
			return
		}

		app.logger.Info("organizationID extracted", "orgID", orgID)

		ctx := r.Context()
		organization, err := app.store.Organizations.GetByID(ctx, orgID)
		if err != nil {
			app.logger.Error("failed to get organization from database", "error", err, "orgID", orgID)
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		app.logger.Info("organization found", "orgID", orgID, "orgName", organization.Name)

		ctx = context.WithValue(ctx, organizationCtx, &organization)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) isUserAdmin(ctx context.Context, orgID, userID int64) bool {
	member, err := app.store.Organizations.GetMember(ctx, orgID, userID)
	if err != nil {
		return false
	}

	return member.RoleID == store.AdminRole
}

func (app *application) isUserMemberOfOrganization(ctx context.Context, orgID, userID int64) bool {
	_, err := app.store.Organizations.GetMember(ctx, orgID, userID)
	return err == nil
}
