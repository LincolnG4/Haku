package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	ctx := r.Context()

	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Users.Create(ctx, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	ctx := r.Context()

	if err := utils.ReadJson(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.unauthorizedBasicAuthErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if !user.Password.Matches(payload.Password) {
		app.unauthorizedBasicAuthErrorResponse(w, r, errors.New("invalid credentials"))
		return
	}

	// generate the token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.expiration).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
