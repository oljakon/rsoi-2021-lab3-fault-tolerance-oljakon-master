package main

import (
	"log"
	"net/http"
	"os"

	"rsoi3/src/car-service/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")

	r := handlers.Router()

	log.Println("server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
	//log.Fatal(http.ListenAndServe(":8081", r))
}
