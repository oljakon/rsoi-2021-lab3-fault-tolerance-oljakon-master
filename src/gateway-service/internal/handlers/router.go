package handlers

import (
	"github.com/gorilla/mux"
)

type ServicesStruct struct {
	PaymentServiceAddress string
	RentalServiceAddress  string
	CarServiceAddress     string
}

type GatewayService struct {
	Config ServicesStruct
}

func NewGatewayService(config *ServicesStruct) *GatewayService {
	return &GatewayService{Config: *config}
}

func Router() *mux.Router {
	servicesConfig := ServicesStruct{
		PaymentServiceAddress: "https://rsoi2-payment-service.herokuapp.com",
		RentalServiceAddress:  "https://rsoi2-rental-service.herokuapp.com",
		CarServiceAddress:     "https://rsoi2-car-service.herokuapp.com",
		//PaymentServiceAddress: "http://localhost:8082/",
		//RentalServiceAddress:  "http://localhost:8083/",
		//CarServiceAddress:     "http://localhost:8081/",
	}

	router := mux.NewRouter()

	gs := NewGatewayService(&servicesConfig)

	router.HandleFunc("/api/v1/cars", gs.GetAvailableCars).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental", gs.GetUserRentals).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}", gs.GetRentalInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental", gs.RentCar).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}/finish", gs.EndRental).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}", gs.CancelRental).Methods("DELETE", "OPTIONS")

	return router
}
