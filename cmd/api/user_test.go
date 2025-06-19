package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/LincolnG4/Haku/internal/auth"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func TestGetUsed(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	// Create a user in the mock store
	user := &store.User{ID: 1, Username: "testuser", Email: "testuser@example.com"}
	app.store.Users.Create(nil, user)

	// Generate a valid token for the user
	claims := &auth.MyClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  fmt.Sprintf("%d", user.ID),
			Issuer:   "test-aud",
			Audience: []string{"test-aud"},
		},
	}
	testToken, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(mux, req)
		checkCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("bad request, string instead a number", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1s", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should allow authenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(mux, req)
		checkCode(t, http.StatusOK, rr.Code)
	})
}
