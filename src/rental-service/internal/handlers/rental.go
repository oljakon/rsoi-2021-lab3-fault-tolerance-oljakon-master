package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rsoi3/src/rental-service/internal/models"
	"rsoi3/src/rental-service/internal/repository"
)

func GetUserRentalsHandler(w http.ResponseWriter, r *http.Request) {
	rentalRepo := repository.RentalRepository{}

	username := r.Header.Get("X-User-Name")

	rentals, err := rentalRepo.GetUserRentalsRepo(username)
	if err != nil {
		log.Fatalf("failed to get rentals: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetRentalHandler(w http.ResponseWriter, r *http.Request) {
	rentalRepo := repository.RentalRepository{}

	params := mux.Vars(r)
	rentalUid := params["rentalUid"]

	rental, err := rentalRepo.GetRentalRepo(rentalUid)
	if err != nil {
		log.Fatalf("failed to get rental: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rental)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateRentalHandler(w http.ResponseWriter, r *http.Request) {
	rentalRepo := repository.RentalRepository{}

	var rental models.Rental

	err := json.NewDecoder(r.Body).Decode(&rental)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rentalRepo.CreateRentalRepo(&rental)
	if err != nil {
		log.Printf("failed to create rental: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func EndRentalHandler(w http.ResponseWriter, r *http.Request) {
	rentalRepo := repository.RentalRepository{}

	params := mux.Vars(r)
	rentalUid := params["rentalUid"]

	carUid, err := rentalRepo.EndRentalRepo(rentalUid)
	if err != nil {
		log.Printf("failed to end rental: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(carUid)
	if err != nil {
		log.Printf("failed to encode car: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CancelRentalHandler(w http.ResponseWriter, r *http.Request) {
	rentalRepo := repository.RentalRepository{}

	params := mux.Vars(r)
	rentalUid := params["rentalUid"]

	carUid, err := rentalRepo.CancelRentalRepo(rentalUid)
	if err != nil {
		log.Printf("failed to end rental: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(carUid)
	if err != nil {
		log.Printf("failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
