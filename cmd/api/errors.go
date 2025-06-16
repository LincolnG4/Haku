package main

import (
	"log"
	"net/http"

	"github.com/LincolnG4/Haku/internal/utils"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	log.Printf("internal error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	utils.WriteJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	log.Printf("bad request: %s path: %s error: %s", r.Method, r.URL.Path, err)

	utils.WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	log.Printf("not found: %s path: %s error: %s", r.Method, r.URL.Path, err)
	utils.WriteJsonError(w, http.StatusNotFound, "resource not found")
}

func (app *application) unauthorizedBasicAuthErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// to do: create a struct
	log.Printf("unauthorized error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	utils.WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	utils.WriteJsonError(w, http.StatusUnauthorized, "unauthorized")
}
