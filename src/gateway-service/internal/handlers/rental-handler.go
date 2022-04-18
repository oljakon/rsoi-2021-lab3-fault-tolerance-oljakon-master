package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rsoi3/src/gateway-service/internal/models"
	"rsoi3/src/gateway-service/internal/service"
)

func (gs *GatewayService) GetUserRentals(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentalsInfo, err := service.UsersRentalWithPaymentController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentalsInfo)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) GetRentalInfo(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentalsInfo, err := service.UsersRentalFullInfoController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if service.NeedsRetry == 1 && rentalsInfo.Payment.Status == "PAID" {
		//service.RetryPaymentCancel(gs.Config.PaymentServiceAddress, rentalsInfo.Payment.PaymentUID)
		err := service.CancelPayment(paymentServiceAddress, rentalsInfo.Payment.PaymentUID)
		if err != nil {
			log.Printf("payment service unavailable")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	service.NeedsRetry = 0

	rentalsInfo, err = service.UsersRentalFullInfoController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rentalsInfo.Payment.PaymentUID == "" {
		rentalInfoEmptyPayment := &models.RentalInfoEmptyPayment{
			RentalUID: rentalsInfo.RentalUID,
			CarUID:    rentalsInfo.CarUID,
			DateFrom:  rentalsInfo.DateFrom,
			DateTo:    rentalsInfo.DateTo,
			Status:    rentalsInfo.Status,
			Car: &models.CarField{
				CarUID:             rentalsInfo.Car.CarUID,
				Brand:              rentalsInfo.Car.Brand,
				Model:              rentalsInfo.Car.Model,
				RegistrationNumber: rentalsInfo.Car.RegistrationNumber,
			},
			Payment: "",
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rentalInfoEmptyPayment)
		if err != nil {
			log.Printf("failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rentalsInfo)
		if err != nil {
			log.Printf("failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (gs *GatewayService) RentCar(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rentInfo models.RentCarRequest
	err := json.NewDecoder(r.Body).Decode(&rentInfo)
	if err != nil {
		fmt.Println("failed to decode post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentedCar, err := service.RentCarController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, &rentInfo)
	if err != nil {
		if errors.Is(err, service.ErrPayment503) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			resp := make(map[string]string)
			resp["message"] = "Payment Service unavailable"
			json.NewEncoder(w).Encode(resp)
			return
		}
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentedCar)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) EndRental(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	err := service.EndRentalController(rentalServiceAddress, carServiceAddress, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gs *GatewayService) CancelRental(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	err := service.CancelRentalController(rentalServiceAddress, carServiceAddress, paymentServiceAddress, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
