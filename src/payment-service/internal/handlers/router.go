package handlers

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/payment/{PaymentUid}", GetPaymentByUidHandler).Methods("GET")
	router.HandleFunc("/api/v1/payment", CreatePaymentHandler).Methods("POST")
	router.HandleFunc("/api/v1/payment/{PaymentUid}", CancelPaymentHandler).Methods("PATCH")

	return router
}
