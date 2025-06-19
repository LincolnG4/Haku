package main

import (
	"net/http"
	"testing"
)

func TestGetUsed(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
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
