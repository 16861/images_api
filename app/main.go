package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rest-and-go/store"

	"github.com/gorilla/handlers"
)

func Run() {
	port := os.Getenv("API_PORT")

	if port == "" {
		log.Fatal("$API_PORT must be set")
	}

	router := store.NewRouter()

	allowedOrigins := handlers.AllowedOrigins([]string{""})
	allowedMethods := handlers.AllowedMethods([]string{"POST", "GET"})

	fmt.Println("Starting server...")
	log.Println("Starting server at port " + port)

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
