package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rsoi3/src/car-service/internal/models"
	"rsoi3/src/car-service/internal/repository"
)

func GetAvailableCarsHandler(w http.ResponseWriter, r *http.Request) {
	carRepo := repository.CarRepository{}

	all := false
	query := r.URL.Query()
	showAll := query.Get("showAll")
	if showAll == "true" {
		all = true
	}

	cars, err := carRepo.GetAvailableCarsRepo(all)
	if err != nil {
		log.Printf("failed to get cars: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		log.Printf("failed to encode cars: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCarByUidHandler(w http.ResponseWriter, r *http.Request) {
	carRepo := repository.CarRepository{}

	params := mux.Vars(r)
	carUID := params["CarUid"]

	car, err := carRepo.GetCarByUidRepo(carUID)
	if err != nil {
		log.Printf("failed to get car: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(car)

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(car)
	if err != nil {
		log.Printf("failed to encode car: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ReserveCarHandler(w http.ResponseWriter, r *http.Request) {
	carRepo := repository.CarRepository{}

	var car models.Car

	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	carPrice, err := carRepo.ReserveCarRepo(car.CarUID)
	if err != nil {
		log.Printf("failed to reserve car: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(carPrice)
	if err != nil {
		log.Printf("failed to encode car: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func EndCarReserveHandler(w http.ResponseWriter, r *http.Request) {
	carRepo := repository.CarRepository{}

	params := mux.Vars(r)
	carUID := params["CarUid"]

	err := carRepo.EndCarReserveRepo(carUID)
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
