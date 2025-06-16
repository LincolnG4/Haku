package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
	"github.com/go-chi/chi"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return

	}

	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		panic("user not found in context")
	}
	return user
}
