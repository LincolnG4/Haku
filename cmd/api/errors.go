package main

import (
	"net/http"

	"github.com/LincolnG4/Haku/internal/utils"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	utils.WriteJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

	utils.WriteJsonError(w, http.StatusForbidden, "forbidden")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	utils.WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	utils.WriteJsonError(w, http.StatusNotFound, "resource not found")
}

func (app *application) unauthorizedBasicAuthErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	utils.WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	utils.WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}
