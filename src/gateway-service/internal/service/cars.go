package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"rsoi3/src/gateway-service/internal/models"
)

func GetAvailableCarsRequest(serviceAddress string, showAll bool) (*[]models.Car, error) {
	requestURL := fmt.Sprintf(serviceAddress + "/api/v1/cars")

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return nil, err
	}

	if showAll {
		q := req.URL.Query()
		q.Add("showAll", "true")
		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to car service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	result := &[]models.Car{}
	err = json.NewDecoder(res.Body).Decode(result)
	log.Println("cars result: ", result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func GetCar(carServiceAddress string, carUID string) (*models.Car, error) {
	requestURL := fmt.Sprintf(carServiceAddress+"/api/v1/cars/%s", carUID)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to car service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	car := &models.Car{}
	err = json.NewDecoder(res.Body).Decode(car)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return car, nil
}

func ReserveCar(carServiceAddress, CarUid string) (int, error) {
	requestURL := fmt.Sprintf(carServiceAddress + "/api/v1/cars")

	car := &models.Car{CarUID: CarUid}

	data, err := json.Marshal(car)
	if err != nil {
		return 0, fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed request to car service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	var responsePrice int
	err = json.NewDecoder(res.Body).Decode(&responsePrice)
	if err != nil {
		return 0, fmt.Errorf("failed to encode response: %w", err)
	}

	return responsePrice, nil
}

func EndCarReserve(carServiceAddress, carUid string) error {
	requestURL := fmt.Sprintf(carServiceAddress+"/api/v1/cars/%s", carUid)

	req, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return nil
}
