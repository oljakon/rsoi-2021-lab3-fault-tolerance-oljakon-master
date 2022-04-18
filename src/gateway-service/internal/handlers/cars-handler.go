package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"rsoi3/src/gateway-service/internal/models"
	"rsoi3/src/gateway-service/internal/service"
)

func (gs *GatewayService) GetAvailableCars(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	log.Printf("get cars method")

	showAllCars := false
	showAllParam := params.Get("showAll")
	if showAllParam == "true" {
		showAllCars = true
	}

	serviceAddress := gs.Config.CarServiceAddress
	availableCars, err := service.GetAvailableCarsRequest(serviceAddress, showAllCars)
	if err != nil {
		log.Printf("failed to get response from car service: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pageParam := params.Get("page")
	if pageParam == "" {
		log.Println("invalid query parameter: page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		log.Printf("unable to convert the string into int:  %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		log.Println("invalid query parameter: size")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		log.Printf("unable to convert the string into int:  %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fin := page * size
	if len(*availableCars) < fin {
		fin = len(*availableCars)
	}
	carsCount := (*availableCars)[(page-1)*size : fin]
	//carsCount := *availableCars
	cars := models.CarsLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(carsCount),
		Items:         &carsCount,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
