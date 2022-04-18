package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rsoi3/src/payment-service/internal/models"
	"rsoi3/src/payment-service/internal/repository"
)

func GetPaymentByUidHandler(w http.ResponseWriter, r *http.Request) {
	paymentRepo := repository.PaymentRepository{}

	params := mux.Vars(r)
	paymentUID := params["PaymentUid"]

	payment, err := paymentRepo.GetPaymentByUidRepo(paymentUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(payment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentRepo := repository.PaymentRepository{}

	var payment models.Payment

	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = paymentRepo.CreatePaymentRepo(&payment)
	if err != nil {
		log.Printf("failed to create payment: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CancelPaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentRepo := repository.PaymentRepository{}

	params := mux.Vars(r)
	paymentUID := params["PaymentUid"]

	err := paymentRepo.CancelPaymentRepo(paymentUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("failed to end car reserve: %v\n", err)
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}
