package main

import (
	"log"
	"net/http"

	"github.com/LincolnG4/Haku/internal/controllers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /pipelines", controllers.CreatePipeline)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
