package handlers

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/rental", GetUserRentalsHandler).Methods("GET")
	router.HandleFunc("/api/v1/rental/{rentalUid}", GetRentalHandler).Methods("GET")
	router.HandleFunc("/api/v1/rental", CreateRentalHandler).Methods("POST")
	router.HandleFunc("/api/v1/rental/{rentalUid}/finish", EndRentalHandler).Methods("PATCH")
	router.HandleFunc("/api/v1/rental/{rentalUid}", CancelRentalHandler).Methods("PATCH")

	return router
}
