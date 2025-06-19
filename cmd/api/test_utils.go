package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LincolnG4/Haku/internal/auth"
	"github.com/LincolnG4/Haku/internal/store"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	testAuth := &auth.TestAuthenticator{}
	return &application{
		logger:        logger,
		store:         mockStore,
		authenticator: testAuth,
	}
}

func executeRequest(mux http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkCode(t *testing.T, expected, response int) {
	if expected != response {
		t.Errorf("should allow authenticated requests: expected %d and we got %d", expected, response)
	}
}
