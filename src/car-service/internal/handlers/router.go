package handlers

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/cars", GetAvailableCarsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/cars/{CarUid}", GetCarByUidHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/cars", ReserveCarHandler).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/v1/cars/{CarUid}", EndCarReserveHandler).Methods("PATCH", "OPTIONS")

	return router
}
